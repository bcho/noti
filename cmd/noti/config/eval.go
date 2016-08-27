package config

import (
	"bytes"
	"text/template"

	"github.com/variadico/noti/cmd/noti/run"
)

type tmplFielder interface {
	TmplFields() []*string
}

func EvalTmplFields(s tmplFielder, st run.Stats) error {
	var err error

	for _, field := range s.TmplFields() {
		*field, err = eval(*field, st)
	}

	return err
}

func eval(s string, st run.Stats) (string, error) {
	tmpl, err := template.New("").Parse(s)
	if err != nil {
		return "", err
	}

	buf := new(bytes.Buffer)
	if err := tmpl.Execute(buf, st); err != nil {
		return "", err
	}

	return buf.String(), nil
}
