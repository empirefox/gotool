package ng

import (
	"encoding/json"
	"errors"
	"io"
	"text/template"

	"github.com/golang/glog"
)

var valueTpl = `angular.module('{{.ModuleName}}', []).{{.Type}}('{{.Name}}', {{marshal .Instance}});`

type Module struct {
	Type       string
	ModuleName string
	Name       string
	Instance   interface{}
}

var (
	value *template.Template
)

func init() {
	var err error
	value, err = template.New("value").Funcs(template.FuncMap{
		"marshal": func(v interface{}) string {
			a, _ := json.Marshal(v)
			return string(a)
		},
	}).Parse(valueTpl)
	if err != nil {
		glog.Errorln(err)
	}
}

func Write(w io.Writer, m Module) error {
	if m.Type == "" || m.ModuleName == "" || m.Name == "" || m.Instance == nil {
		return errors.New("All fields must be set")
	}
	return value.ExecuteTemplate(w, "value", m)
}
