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
	Src       string
	Data      interface{}
	TransData i18n.TrunsData
	Locale    string

	tmpl *template.Template
}

// NewSiteTemplate add site path
func NewSiteTemplate(folder, name string) *Template {
	fileName := filepath.Join(HomeDir, "app/view/", folder, name+DefaultTemplateSubfix)
	return NewTemplate(fileName)
}

// NewTemplate new
func NewTemplate(fileName string) *Template {
	tmpl := template.New(filepath.Base(fileName))

	if TemplateLeftDelim != "" {
		tmpl.Delims(TemplateLeftDelim, TemplateRightDelim)
	}

	// filename := filepath.Join("app/view/", fileName)
	var err error
	tmpl, err = tmpl.ParseFiles(fileName)
	if err != nil {
		panic(fileName + "\n" + err.Error())
	}
	return &Template{FileName: fileName, tmpl: tmpl}
}

// NewTemplateS new
func NewTemplateS(src string) *Template {
	tmpl := template.New("")
	if TemplateLeftDelim != "" {
		tmpl.Delims(TemplateLeftDelim, TemplateRightDelim)
	}
	tmpl.Parse(src)
	return &Template{Src: src, tmpl: tmpl}
}

// Render template
func (t *Template) Render(w io.Writer) error {
	return t.RenderWithData(w, t.Data)
}

// RenderWithData template
func (t *Template) RenderWithData(w io.Writer, data interface{}) error {
	// if t.tmpl == nil {
	//
	// 	if t.Locale != "" && t.loadi18n() == nil {
	// 		t.tmpl.Funcs(template.FuncMap{"T": t.TransData.T})
	// 	} else {
	// 		t.tmpl.Funcs(template.FuncMap{"T": i18n.EmptyTransFunc})
	// 	}
	// }

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
