package httpd

import (
	"html/template"
	"io"
	"io/ioutil"
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

	rawSrc []byte // static file template
	tmpl   *template.Template
}

// NewSiteTemplate add site path
func NewSiteTemplate(folder, name string) *Template {
	fileName := filepath.Join(HomeDir, "app/view/", folder, name+DefaultTemplateSubfix)
	return NewTemplate(fileName)
}

// NewTemplate new
func NewTemplate(fileName string) *Template {
	return &Template{FileName: fileName}
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
	return t.RenderD(w, t.Data)
}

// RenderD template
func (t *Template) RenderD(w io.Writer, data interface{}) error {
	if RunMode == ModeDev && len(t.rawSrc) > 0 {
		t.rawSrc = nil
	}

	if len(t.rawSrc) > 0 {
		w.Write(t.rawSrc)
		return nil
	}

	// static htm template
	if t.tmpl == nil && data == nil {
		var err error
		t.rawSrc, err = ioutil.ReadFile(t.FileName)
		if err != nil {
			return err
		}
		w.Write(t.rawSrc)
		return nil
	}

	// dynamic htm
	if t.tmpl == nil {
		t.tmpl = template.New(filepath.Base(t.FileName))

		if TemplateLeftDelim != "" {
			t.tmpl.Delims(TemplateLeftDelim, TemplateRightDelim)
		}

		var err error
		t.tmpl, err = t.tmpl.ParseFiles(t.FileName)
		if err != nil {
			panic(t.FileName + "\n" + err.Error())
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
