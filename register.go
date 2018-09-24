package gno

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/kere/gno/websock"
)

// RegistGet router
func (s *SiteServer) RegistGet(rule string, factory func() IPage) {
	s.Router.GET(rule, func(rw http.ResponseWriter, req *http.Request, ps httprouter.Params) {
		p := factory()
		p.Init("GET", rw, req, ps)
		err := doPageHandle(p, rw, req, ps)
		if err == nil {
			return
		}
		// do error
		doPageError(s.ErrorURL, err, rw, req)
	})
}

// RegistPost router
func (s *SiteServer) RegistPost(rule string, factory func() IPage) {
	s.Router.POST(rule, func(rw http.ResponseWriter, req *http.Request, ps httprouter.Params) {
		p := factory()
		p.Init("POST", rw, req, ps)
		err := doPageHandle(p, rw, req, ps)
		if err == nil {
			return
		}
		// do error
		doPageError(s.ErrorURL, err, rw, req)
	})
}

// RegistAPI api router
func (s *SiteServer) RegistAPI(rule string, webapi IWebAPI) {
	f := func(rw http.ResponseWriter, req *http.Request, ps httprouter.Params) {
		err := doAPIHandle(webapi, rw, req, ps)
		if err == nil {
			return
		}
		// do error
		doAPIError(err, rw)
	}

	s.Router.POST(rule, f)
	s.Router.GET(rule, f)
}

// RegistPostAPI api router
func (s *SiteServer) RegistPostAPI(rule string, webapi IWebAPI) {
	s.Router.POST(rule, func(rw http.ResponseWriter, req *http.Request, ps httprouter.Params) {
		// webapi := factory()
		// var err error
		// webapi, err = webapi.Init(rw, req, ps)
		// if err != nil {
		// 	Site.Log.Error(err)
		// 	return
		// }
		err := doAPIHandle(webapi, rw, req, ps)
		if err == nil {
			return
		}
		// do error
		doAPIError(err, rw)
	})
}

// RegistGetAPI api router
func (s *SiteServer) RegistGetAPI(rule string, webapi IWebAPI) {
	s.Router.GET(rule, func(rw http.ResponseWriter, req *http.Request, ps httprouter.Params) {
		err := doAPIHandle(webapi, rw, req, ps)
		if err == nil {
			return
		}
		// do error
		doAPIError(err, rw)
	})
}

// RegistWebSocket router
func (s *SiteServer) RegistWebSocket(rule string, ctl websock.IWebSock) {
	websock.RegistWebSocket(s.Router, rule, ctl)
}
