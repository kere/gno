package gno

import (
	"net/http"
	"reflect"

	"github.com/julienschmidt/httprouter"
	"github.com/kere/gno/libs/util"
	"github.com/kere/gno/websock"
)

type openAPIExec func(req *http.Request, ps httprouter.Params, args util.MapData) (interface{}, error)

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
		if name == "Auth" || name == "Reply" {
			continue
		}
		f := v.Method(i).Interface().(func(req *http.Request, ps httprouter.Params, args util.MapData) (interface{}, error))
		openapiMap[rule+"/"+name] = openapiItem{Exec: f, API: openapi}

		// s.Router.POST(rule+"/"+name, doOpenAPIHandle)
		s.Router.POST(rule+"/"+name, func(rw http.ResponseWriter, req *http.Request, ps httprouter.Params) {
			arg := PoolParams{Typ: 2, RW: rw, Req: req, Params: ps, Error: make(chan error, 1)}
			pool.Serve(&arg)
			err := <-arg.Error
			if err != nil {
				doAPIError(err, rw, req)
			}
		})
	}
}

// RegistGet router
func (s *SiteServer) RegistGet(rule string, factory func() IPage) {
	s.Router.GET(rule, func(rw http.ResponseWriter, req *http.Request, ps httprouter.Params) {
		p := factory()
		p.Init("GET", rw, req, ps)
		arg := PoolParams{Typ: 1, RW: rw, Req: req, Params: ps, Page: p, Error: make(chan error, 1)}
		pool.Serve(&arg)
		err := <-arg.Error
		if err != nil {
			doPageError(s.ErrorURL, err, rw, req)
		}
	})
}

// RegistPost router
func (s *SiteServer) RegistPost(rule string, factory func() IPage) {
	s.Router.POST(rule, func(rw http.ResponseWriter, req *http.Request, ps httprouter.Params) {
		p := factory()
		p.Init("POST", rw, req, ps)
		arg := PoolParams{Typ: 1, RW: rw, Req: req, Params: ps, Page: p, Error: make(chan error, 1)}
		pool.Serve(&arg)
		err := <-arg.Error
		if err != nil {
			doPageError(s.ErrorURL, err, rw, req)
		}
	})
}

// RegistPostAPI api router
func (s *SiteServer) RegistPostAPI(rule string, webapi IWebAPI) {
	s.Router.POST(rule, func(rw http.ResponseWriter, req *http.Request, ps httprouter.Params) {
		arg := PoolParams{Typ: 3, RW: rw, Req: req, Params: ps, WebAPI: webapi, Error: make(chan error, 1)}
		pool.Serve(&arg)
		err := <-arg.Error
		if err != nil {
			doAPIError(err, rw, req)
		}
	})
}

// RegistGetAPI api router
func (s *SiteServer) RegistGetAPI(rule string, webapi IWebAPI) {
	s.Router.GET(rule, func(rw http.ResponseWriter, req *http.Request, ps httprouter.Params) {
		arg := PoolParams{Typ: 3, RW: rw, Req: req, Params: ps, WebAPI: webapi, Error: make(chan error, 1)}
		pool.Serve(&arg)
		err := <-arg.Error
		if err != nil {
			doAPIError(err, rw, req)
		}
	})
}

// RegistWebSocket router
func (s *SiteServer) RegistWebSocket(rule string, ctl websock.IWebSock) {
	websock.RegistWebSocket(s.Router, rule, ctl)
}
