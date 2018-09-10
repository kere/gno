package gno

import (
	"net/http"

	"github.com/gorilla/websocket"
	"github.com/julienschmidt/httprouter"
	"github.com/kere/gno/libs/log"
	"github.com/kere/gno/websock"
)

var upgrader websocket.Upgrader

// RegistMessageSocket router
func (s *SiteServer) RegistMessageSocket(path string, ctl websock.IMessageSock) {
	s.Router.GET(path, func(rw http.ResponseWriter, req *http.Request, ps httprouter.Params) {
		err := ctl.Auth(req)
		if err != nil {
			log.App.Error(err)
			return
		}

		conn, err := upgrader.Upgrade(rw, req, nil)
		if err != nil {
			return
		}
		defer conn.Close()

		ctl.Listen(conn)
	})

	// Handler f
}

// RegistWebSocket router
func (s *SiteServer) RegistWebSocket(path string, ctl websock.IWebSock) {
	s.Router.GET(path, func(rw http.ResponseWriter, req *http.Request, ps httprouter.Params) {
		err := ctl.Auth(req)
		if err != nil {
			log.App.Error(err)
			return
		}

		conn, err := upgrader.Upgrade(rw, req, nil)
		if err != nil {
			return
		}
		defer conn.Close()

		ctl.Listen(conn)
	})

	// Handler f
}
