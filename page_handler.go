package gno

import (
	"bytes"
	"net/http"
	"net/url"

	"github.com/kere/gno/libs/log"
)

func pageHandle(p IPage) error {
	uri, err := p.Auth()
	if uri != "" || err != nil {
		if uri != "" && err != nil {
			// add msg after url
			var u *url.URL
			u, err = url.Parse(uri)
			if err != nil {
				return err
			}
			if u.RawQuery == "" {
				u.RawQuery = "msg=" + url.PathEscape(err.Error())
			} else {
				u.RawQuery += "&msg=" + url.PathEscape(err.Error())
			}
			uri = u.String()
		} else if err != nil {
			uri = "/error?msg=" + url.PathEscape(err.Error())
		}

		http.Redirect(p.GetResponseWriter(), p.GetRequest(), uri, http.StatusSeeOther)
		return nil
	}

	if TryCache(p) {
		log.App.Debug("Page Cache", p.GetRequest().URL.String())
		return nil
	}

	err = p.Prepare()
	if err != nil {
		return err
	}

	buf := bytes.NewBuffer(nil)

	err = p.Render(buf)
	if err != nil {
		return err
	}

	TrySetCache(p, buf)
	_, err = p.GetResponseWriter().Write(buf.Bytes())

	p.RunAfter()
	return err
}
