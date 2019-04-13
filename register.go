package gno

import (
	"errors"
	"net/http"
	"reflect"

	"github.com/julienschmidt/httprouter"
	"github.com/kere/gno/libs/log"
	"github.com/kere/gno/libs/util"
	"github.com/kere/gno/websock"
)

type openAPIExec func(args util.MapData, prepareDat interface{}) (interface{}, error)

type openapiItem struct {
	Exec openAPIExec
	API  IOpenAPI
}

var openapiMap = make(map[string]openapiItem)

// RegistOpenAPI init open api
func (s *SiteServer) RegistOpenAPI(rule string, openapi IOpenAPI) {
	v := reflect.ValueOf(openapi)
	typ := v.Type()
	l := typ.NumMethod()
	for i := 0; i < l; i++ {
		m := typ.Method(i)
		name := m.Name
		if name == "Prepare" || name == "Reply" {
			continue
		}

		f := v.Method(i).Interface().(func(args util.MapData, dat interface{}) (interface{}, error))
		openapiMap[rule+"/"+name] = openapiItem{Exec: f, API: openapi}

		// s.Router.POST(rule+"/"+name, openAPIHandle)

		s.Router.POST(rule+"/"+name, func(rw http.ResponseWriter, req *http.Request, ps httprouter.Params) {
			// err := openAPIHandle(rw, req, ps)
			arg := PoolParams{Typ: 2, RW: rw, Req: req, Params: ps}
			if RunMode == ModeDev {
				InvokeExec(arg)
				return
			}
			if err := pool.Invoke(arg); err != nil {
				doAPIError(errors.New("Throttle limit error"), rw, req)
			}
		})
	}
}

// RegistGet router
func (s *SiteServer) RegistGet(rule string, factory func() IPage) {
	s.Router.GET(rule, func(rw http.ResponseWriter, req *http.Request, ps httprouter.Params) {
		p := factory()
		p.Init("GET", rw, req, ps)
		// err := pageHandle(p)
		arg := PoolParams{Typ: 1, Page: p}
		if RunMode == ModeDev {
			InvokeExec(arg)
			return
		}
		if err := pool.Invoke(arg); err != nil {
			doAPIError(errors.New("Throttle limit error"), rw, req)
		}
	})
}

// RegistPost router
func (s *SiteServer) RegistPost(rule string, factory func() IPage) {
	s.Router.POST(rule, func(rw http.ResponseWriter, req *http.Request, ps httprouter.Params) {
		p := factory()
		p.Init("POST", rw, req, ps)
		// err := pageHandle(p, rw, req, ps)

		arg := PoolParams{Typ: 1, Page: p}
		if RunMode == ModeDev {
			InvokeExec(arg)
			return
		}
		if err := pool.Invoke(arg); err != nil {
			doPageError(s.ErrorURL, errors.New("Throttle limit error"), rw, req)
		}
	})
}

// RegistPostAPI api router
func (s *SiteServer) RegistPostAPI(rule string, webapi IWebAPI) {
	s.Router.POST(rule, func(rw http.ResponseWriter, req *http.Request, ps httprouter.Params) {
		arg := PoolParams{Typ: 3, RW: rw, Req: req, Params: ps, WebAPI: webapi}
		if RunMode == ModeDev {
			InvokeExec(arg)
			return
		}
		if err := pool.Invoke(arg); err != nil {
			doAPIError(errors.New("Throttle limit error"), rw, req)
		}

	})
}

// RegistGetAPI api router
func (s *SiteServer) RegistGetAPI(rule string, webapi IWebAPI) {
	s.Router.GET(rule, func(rw http.ResponseWriter, req *http.Request, ps httprouter.Params) {
		arg := PoolParams{Typ: 3, RW: rw, Req: req, Params: ps, WebAPI: webapi}
		if RunMode == ModeDev {
			InvokeExec(arg)
			return
		}
		if err := pool.Invoke(arg); err != nil {
			doAPIError(errors.New("Throttle limit error"), rw, req)
		}
	})
}

// RegistWebSocket router
func (s *SiteServer) RegistWebSocket(rule string, ctl websock.IWebSock) {
	websock.RegistWebSocket(s.Router, rule, ctl)
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

func doAPIError(err error, rw http.ResponseWriter, req *http.Request) {
	addr := req.Header.Get("X-Forwarded-For")
	if addr == "" {
		addr = req.Header.Get("X-Real-IP")
	}
	log.App.Error(err, addr)
	rw.WriteHeader(http.StatusInternalServerError)
	rw.Write([]byte(err.Error()))
}
