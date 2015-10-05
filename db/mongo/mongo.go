package mongo

import (
	"errors"
	"reflect"
	"strings"
	"time"

	"github.com/codeblanche/golibs/logr"

	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

var (
	// Session default db connection
	Session *mgo.Session
)

// OID type alias of mgo's bson.ObjectId
type OID interface {
	Counter() int32
	Hex() string
	Machine() []byte
	MarshalJSON() ([]byte, error)
	Pid() uint16
	String() string
	Time() time.Time
	// UnmarshalJSON([]byte) error //takes a pointer receiver
	Valid() bool
}

// M type alias of mgo's bson.M
type M bson.M

// Model defines the interface for models to be stored in mongo
type Model interface {
	ObjectID() OID
	DBName() string
	CName() string
	SetObjectID(interface{})
}

// SessionError for mongo session failures
type SessionError struct{}

// ErrorList collection of multiple errors
type ErrorList []error

// Error implements error interface
func (e SessionError) Error() string {
	return "Mongo session is not valid or nil"
}

// Error implements error interface
func (e ErrorList) Error() string {
	err := ""
	for _, one := range e {
		err += one.Error() + "\n"
	}
	return err
}

// Index ensures an index on the given keys
func Index(m Model, key ...string) error {
	s := session()
	defer s.Close()

	return c(s, m).EnsureIndexKey(key...)
}

// Load a db model using it's ObjectID
func Load(m Model) error {
	return One(M{"_id": m.ObjectID()}, m)
}

// LoadList loads a slice of models using their ObjectIDs
func LoadList(l interface{}) error {
	list, ok := l.([]Model)
	if !ok {
		return errors.New("Unable to use %T as []mongo.Model")
	}
	var ids []OID
	for _, e := range list {
		e := e.(Model)
		ids = append(ids, e.ObjectID())
	}
	return Find(M{"$in": ids}, "", l)
}

// Remove removes a Document by Object Id
func Remove(m ...Model) error {
	var errl ErrorList

	s := session()
	defer s.Close()

	// TODO optimize this into a single query per collection
	for _, model := range m {
		id := model.ObjectID()
		if id == nil || !id.Valid() {
			continue
		}
		err := c(s, model).Remove(M{"_id": id})
		if err != nil {
			errl = append(errl, err)
		}
	}
	if len(errl) > 0 {
		return errl
	}
	return nil

}

// Update Load a db model matching the given conditions
func Update(m Model, query M) error {
	m, err := resolveModel(m)
	if err != nil {
		return err
	}

	s := session()
	defer s.Close()

	err = c(s, m).Update(M{"_id": m.ObjectID()}, query)
	if err != nil {
		return err
	}
	return nil
}

// Count the number of document matching the given query
// TODO
func Count(m Model, query M) int {
	return 0
}

// One finds a single document
func One(query M, into interface{}) error {
	m, err := resolveModel(into)
	if err != nil {
		return err
	}

	s := session()
	defer s.Close()

	return c(s, m).Find(query).One(into)
}

// Find documents
func Find(query M, sort string, into interface{}) error {
	m, err := resolveModel(into)
	if err != nil {
		return err
	}

	s := session()
	defer s.Close()

	q := c(s, m).Find(query)

	if len(sort) > 0 {
		so := strings.Split(sort, ",")
		q.Sort(so...)
	}

	return q.All(into)
}

// Page is much live Find with the addition of paging
func Page(query M, sort string, page, size int, into interface{}) error {
	m, err := resolveModel(into)
	if err != nil {
		return err
	}

	s := session()
	defer s.Close()

	q := c(s, m).Find(query)

	if len(sort) > 0 {
		so := strings.Split(sort, ",")
		q.Sort(so...)
	}

	q.Skip(size * (page - 1))
	q.Limit(size)

	return q.All(into)
}

// resolveModel resolves the Model implementing object from a slice or pointer
func resolveModel(v interface{}) (Model, error) {
	t := reflect.TypeOf(v)
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}
	if t.Kind() == reflect.Slice {
		m, ok := reflect.New(t.Elem()).Interface().(Model)
		if ok {
			return m, nil
		}

	}
	if t.Kind() == reflect.Struct {
		m, ok := v.(Model)
		if ok {
			return m, nil
		}
	}
	return nil, errors.New("Parameter 'into' must be a slice of mongo.Model elements or a mongo.Model")
}

// Save the given db models
func Save(m ...Model) error {
	var errl ErrorList

	s := session()
	defer s.Close()

	for _, model := range m {
		id := model.ObjectID()
		if id == nil || !id.Valid() {
			id = bson.NewObjectId()
			model.SetObjectID(id)
		}
		_, err := c(s, model).Upsert(M{"_id": id}, model)
		if err != nil {
			errl = append(errl, err)
		}
	}
	if len(errl) > 0 {
		return errl
	}
	return nil
}

// session returns a *mgo.Session
func session() *mgo.Session {
	return Session.Copy()
}

// c returns a *mgo.Collection
func c(s *mgo.Session, m Model) *mgo.Collection {
	return s.DB(m.DBName()).C(m.CName())
}

// ToOID converts any given value into the internal OID type
func ToOID(id interface{}) OID {
	switch t := id.(type) {
	case string:
		id := id.(string)
		if bson.IsObjectIdHex(id) {
			return OID(bson.ObjectIdHex(id))
		}
	case bson.ObjectId:
		id := id.(bson.ObjectId)
		return OID(id)
	case OID:
		id := id.(OID)
		return id
	default:
		logr.Debugf("Unable to convert type: %T to mongo.OID", t)
	}
	return nil
}
