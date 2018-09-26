package websock

import (
	"net/http"

	"github.com/gorilla/websocket"
	"github.com/julienschmidt/httprouter"
	"github.com/kere/gno/libs/log"
)

var upgrader websocket.Upgrader

// RegistWebSocket router
func RegistWebSocket(router *httprouter.Router, path string, ctl IWebSock) {
	connMap[path] = NewManager()
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

		m := connMap[path]
		id := m.GenrateClientID()
		client := m.SetClient(id, conn, req)

		// conn.SetCloseHandler(func(code int, text string) error {
		// 	return nil
		// })
		// fmt.Println(m.ClientCount())

		defer m.Close(id)

		client.Listen(ctl)
	})
}
