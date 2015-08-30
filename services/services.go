package services

import "sync"

// Service interface
type Service interface {
	Init()
}

type singleton struct {
	*sync.Once
	instance Service
}

func (s *singleton) init() Service {
	s.Do(s.instance.Init)
	return s.instance
}

var (
	mutex    = sync.RWMutex{}
	services = map[string]*singleton{}
)

// Register a service with a name
func Register(name string, service Service) {
	mutex.Lock()
	services[name] = &singleton{
		instance: service,
	}
	mutex.Unlock()
}

// Get a service by name
func Get(name string) (Service, bool) {
	if s, ok := services[name]; ok {
		return s.init(), true
	}
	return nil, false
}

// GetInto loads a service into the given pointer by name
func GetInto(name string, into Service) bool {
	s, ok := Get(name)
	into = s
	return ok
}
