package render

import (
	"html/template"
	"io"

	"github.com/kere/gno/libs/util"
	"github.com/valyala/bytebufferpool"
)

// Text class
type Text struct {
	Source string
	Data   map[string]interface{}
}

// NewText new
func NewText(txt string) Text {
	return Text{Source: txt}
}

// Render func
func (t Text) Render(w io.Writer) error {
	tmpl, _ := template.New("").Parse(t.Source)
	return tmpl.Execute(w, t.Data)
}

// String class
type String struct {
	Source string
}

// NewString new
func NewString(src string) String {
	return String{Source: src}
}

// Render func
func (t String) Render(w io.Writer) error {
	w.Write(util.Str2Bytes(t.Source))
	return nil
}

// RenderWith func
func (t String) RenderWith(w io.Writer, opt Opt) error {
	w.Write(util.Str2Bytes(t.Source))
	return nil
}

// Script return string
func Script(src string, data map[string]string) String {
	buf := bytebufferpool.Get()
	buf.WriteString("<script")

	// var s string
	if len(data) > 0 {
		// s = " "
		buf.WriteByte(' ')
		for k, v := range data {
			// s += k + "=\"" + v + "\" "
			buf.WriteString(k)
			buf.WriteString("=\"")
			buf.WriteString(v)
			buf.WriteString("\" ")
		}
	} else {
		buf.WriteString(" type=\"text/javascript\"")
	}

	// str += s + ">" + src + "</script>"
	buf.WriteString(">")
	buf.WriteString(src)
	buf.WriteString("</script>")
	r := NewString(buf.String())
	bytebufferpool.Put(buf)

	return r
}
