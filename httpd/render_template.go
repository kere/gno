package httpd

import (
	"html/template"
	"io"
	"path/filepath"

	"github.com/kere/gno/libs/i18n"
)

// Template render class
type Template struct {
	FileName  string
	Data      interface{}
	TransData i18n.TrunsData
	Locale    string

	tmpl *template.Template
}

// NewTemplate new
func NewTemplate(fileName string) *Template {
	return &Template{FileName: fileName}
}

// Render template
func (t *Template) Render(w io.Writer) error {
	return t.RenderWithData(w, t.Data)
}

// RenderWithData template
func (t *Template) RenderWithData(w io.Writer, data interface{}) error {
	if t.tmpl == nil {
		filename := filepath.Join("app/view/", t.FileName)
		var err error
		t.tmpl, err = template.ParseFiles(filename)
		if err != nil {
			return err
		}

		if TemplateLeftDelim != "" {
			t.tmpl.Delims(TemplateLeftDelim, TemplateRightDelim)
		}

		if t.Locale != "" && t.loadi18n() == nil {
			t.tmpl.Funcs(template.FuncMap{"T": t.TransData.T})
		} else {
			t.tmpl.Funcs(template.FuncMap{"T": i18n.EmptyTransFunc})
		}
	}

	return t.tmpl.Execute(w, data)
}

func (t *Template) loadi18n() error {
	dir, name := filepath.Split(t.FileName)

	d, err := i18n.Load(t.Locale, filepath.Join(dir, "lang", t.Locale, name+".json"))
	if err != nil {
		return err
	}
	t.TransData = d
	return nil
}
