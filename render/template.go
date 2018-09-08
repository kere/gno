package render

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
}

// NewTemplate new
func NewTemplate(fileName string) Template {
	return Template{FileName: fileName}
}

// Render func
func (t Template) Render(w io.Writer) error {
	tmpl := template.New(filepath.Base(t.FileName))
	if TemplateLeftDelim != "" {
		tmpl.Delims(TemplateLeftDelim, TemplateRightDelim)
	}

	filename := filepath.Join("app/view/", t.FileName)
	tmpl, err := tmpl.ParseFiles(filename)
	if err != nil {
		return err
	}

	if t.Locale != "" && t.loadi18n() == nil {
		tmpl.Funcs(template.FuncMap{"T": t.TransData.T})
	} else {
		tmpl.Funcs(template.FuncMap{"T": i18n.EmptyTransFunc})
	}

	return tmpl.Execute(w, t.Data)
}

func (t Template) loadi18n() error {
	dir, name := filepath.Split(t.FileName)

	d, err := i18n.Load(t.Locale, filepath.Join(dir, "lang", t.Locale, name+".json"))
	if err != nil {
		return err
	}
	t.TransData = d
	return nil
}
