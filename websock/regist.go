package websock

import (
	"net/http"

	"github.com/gorilla/websocket"
	"github.com/julienschmidt/httprouter"
	"github.com/kere/gno/libs/log"
)

var upgrader websocket.Upgrader

// make(map[path] map[clientID]*websocket.Conn, 0)
var connMap = make(map[string]map[int]*websocket.Conn, 0)

// GetConn get connection
func GetConn(path string, id int) *websocket.Conn {
	if clientMap, isok := connMap[path]; isok {
		return clientMap[id]
	}
	return nil
}

// GetConnMap get connections
func GetConnMap(path string) map[int]*websocket.Conn {
	return connMap[path]
}

// RegistMessageSocket router
func RegistMessageSocket(router *httprouter.Router, path string, ctl IMessageSock) {
	router.GET(path, func(rw http.ResponseWriter, req *http.Request, ps httprouter.Params) {
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

		clientMap, isok := connMap[path]
		if !isok {
			clientMap = make(map[int]*websocket.Conn, 0)
			connMap[path] = clientMap
		}

		id := CurrentClientID()
		clientMap[id] = conn
		defer func() {
			conn.Close()
			delete(clientMap, id)
		}()

		ctl.Listen(conn)
	})
}

// RegistWebSocket router
func RegistWebSocket(router *httprouter.Router, path string, ctl IWebSock) {
	router.GET(path, func(rw http.ResponseWriter, req *http.Request, ps httprouter.Params) {
		err := ctl.Auth(req)
		if err != nil {
			log.App.Error(err)
			return
		}

		conn, err := upgrader.Upgrade(rw, req, nil)
		if err != nil {
			return
		}

		clientMap, isok := connMap[path]
		if !isok {
			clientMap = make(map[int]*websocket.Conn, 0)
			connMap[path] = clientMap
		}

		id := CurrentClientID()

		clientMap[id] = conn
		defer func() {
			conn.Close()
			delete(clientMap, id)
		}()

		ctl.Listen(conn)
	})
}
