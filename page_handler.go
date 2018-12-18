package gno

import (
	"bytes"
	"net/http"
	"net/url"

	"github.com/kere/gno/libs/log"
)

func pageHandle(p IPage) error {
	isReq, isOK, urlstr, err := p.Auth()
	if isReq && !isOK {
		if urlstr != "" {
			u, _ := url.Parse(urlstr)
			if u.RawQuery == "" {
				u.RawQuery = "msg=" + url.PathEscape(err.Error())
			} else {
				u.RawQuery += "&msg=" + url.PathEscape(err.Error())
			}

			http.Redirect(p.GetResponseWriter(), p.GetRequest(), u.String(), http.StatusSeeOther)
		}
		return nil
	}

	if TryCache(p) {
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

	p.Finish()
	return err
}

func doPageError(errorURL string, err error, rw http.ResponseWriter, req *http.Request) {
	log.App.Error(err)
	if errorURL == "" {
		rw.WriteHeader(http.StatusInternalServerError)
		rw.Write([]byte(err.Error()))
		return
	}
	// ErrorURL redirect to
	http.Redirect(rw, req, errorURL+"?msg="+err.Error(), http.StatusSeeOther)
}
