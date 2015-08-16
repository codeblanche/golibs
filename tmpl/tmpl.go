package tmpl

import (
	"fmt"
	"html/template"
	"io"
	"io/ioutil"

	"github.com/codeblanche/golibs/logr"
)

type (
	// Tmpl wraps the native template.Template
	Tmpl struct {
		template *template.Template
	}

	// Layout provides the API for a selected layout
	Layout struct {
		tmpl    Tmpl
		context interface{}
		name    string
	}
)

var (
	fmap = template.FuncMap{
		"yield": func() string {
			// placeholder to be filled at render time
			return ""
		},
		"raw": func(in interface{}) template.HTML {
			return template.HTML(fmt.Sprintf("%s", in))
		},
	}
)

// New creates a new Tmpl
func New(name string) *Tmpl {
	tmpl := &Tmpl{
		template: template.New(name),
	}
	return tmpl
}

// Load a new template
func (t *Tmpl) Load(name string, file string) error {
	f, err := ioutil.ReadFile(file)
	if err != nil {
		logr.Error(err.Error())
		return nil
	}
	t.template.New(name).Funcs(fmap).Parse(string(f))
	return nil
}

// Templates wraps template.Templates
func (t *Tmpl) Templates() []*template.Template {
	return t.template.Templates()
}

// ExecuteTemplate wraps template.ExecuteTemplate
func (t *Tmpl) ExecuteTemplate(w io.Writer, lname string, pname string, data interface{}) error {
	fm := fmap
	fm["yield"] = func() string {
		t.template.ExecuteTemplate(w, pname, data)
		return ""
	}
	return t.template.Funcs(fm).ExecuteTemplate(w, lname, data)
}

// With layout by given name
func (t Tmpl) With(name string, context interface{}) *Layout {
	logr.Debugf("Layout: %s", name)
	l := t.Templates()
	for _, lt := range l {
		logr.Debug(lt.Name())
	}
	return &Layout{
		tmpl:    t,
		context: context,
		name:    name,
	}
}

// Render a template into the layout with the given context
func (l *Layout) Render(name string, w io.Writer) error {
	logr.Debugf("Render: %s", name)
	err := l.tmpl.ExecuteTemplate(w, l.name, name, l.context)
	return err
}
