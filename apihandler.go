package gno

import (
	"crypto/md5"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"

	"github.com/julienschmidt/httprouter"
	"github.com/kere/gno/libs/log"
	"github.com/kere/gno/libs/util"
)

const (
	// APIFieldSrc post field
	APIFieldSrc = "_src"
	// APIFieldTS post field
	APIFieldTS = "now"
	// APIFieldMethod post field
	APIFieldMethod = "method"
	// APIFieldToken post field
	APIFieldToken = "token"
)

func doPageHandle(p IPage, rw http.ResponseWriter, req *http.Request, ps httprouter.Params) error {
	isReq, isOK, urlstr, err := p.Auth()
	if isReq && !isOK {
		if urlstr != "" {
			u, _ := url.Parse(urlstr)
			if u.RawQuery == "" {
				u.RawQuery = "msg=" + url.PathEscape(err.Error())
			} else {
				u.RawQuery += "&msg=" + url.PathEscape(err.Error())
			}

			http.Redirect(rw, req, u.String(), http.StatusSeeOther)
		}
		return nil
	}

	err = p.Prepare()
	if err != nil {
		return err
	}

	return p.Render()
}

func doAPIHandle(webapi IWebAPI, rw http.ResponseWriter, req *http.Request, ps httprouter.Params) error {
	if isReq, err := webapi.Auth(); isReq && err != nil {
		return err
	}

	var args util.MapData
	src := req.PostFormValue(APIFieldSrc)
	if src != "" {
		err := json.Unmarshal([]byte(src), &args)
		if err != nil {
			return err
		}
	}

	if !webapi.IsSkipToken(req.Method) {
		// method+now+jsonStr+now
		ts := req.PostFormValue(APIFieldTS)
		token := req.PostFormValue(APIFieldToken)
		method := req.PostFormValue(APIFieldMethod)

		u32 := fmt.Sprintf("%x", md5.Sum([]byte(ts+method+ts+src+ts)))
		if u32 != token {
			return errors.New("api token failed")
		}
	}

	data, err := webapi.Exec(args)
	if err != nil {
		return err
	}

	return webapi.Reply(data)
}

func doOpenAPIHandle(rw http.ResponseWriter, req *http.Request, ps httprouter.Params) {
	uri := req.URL.Path
	// method := ps.ByName(methodFieldName)

	var item openapiItem
	var isok bool
	if item, isok = openapiMap[uri]; !isok {
		doAPIError(errors.New(uri+" openapi not found"), rw)
		return
	}

	if isReq, err := item.API.Auth(req, ps); isReq && err != nil {
		doAPIError(err, rw)
		return
	}

	var args util.MapData
	src := req.PostFormValue(APIFieldSrc)
	if src != "" {
		err := json.Unmarshal([]byte(src), &args)
		if err != nil {
			doAPIError(err, rw)
			return
		}
	}

	ts := req.PostFormValue(APIFieldTS)
	token := req.PostFormValue(APIFieldToken)
	method := req.PostFormValue(APIFieldMethod)

	// method+now+jsonStr+now
	u32 := fmt.Sprintf("%x", md5.Sum([]byte(ts+method+ts+src+ts)))
	if u32 != token {
		doAPIError(errors.New("open api token failed"), rw)
		return
	}

	data, err := item.Exec(req, ps, args)
	if err != nil {
		doAPIError(err, rw)
		return
	}

	err = item.API.Reply(rw, data)
	if err != nil {
		doAPIError(err, rw)
		return
	}
}

func doPageError(errorURL string, err error, rw http.ResponseWriter, req *http.Request) {
	fmt.Println("doPageErr", err)
	log.App.Warn(err)
	if errorURL == "" {
		rw.WriteHeader(http.StatusInternalServerError)
		rw.Write([]byte(err.Error()))
		return
	}
	// ErrorURL redirect to
	http.Redirect(rw, req, errorURL+"?msg="+err.Error(), http.StatusSeeOther)
}

func doAPIError(err error, rw http.ResponseWriter) {
	log.App.Warn(err)
	rw.WriteHeader(http.StatusInternalServerError)
	rw.Write([]byte(err.Error()))
}
