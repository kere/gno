package websock

import (
	"net/http"

	"github.com/gorilla/websocket"
	"github.com/julienschmidt/httprouter"
	"github.com/kere/gno/libs/log"
)

var upgrader websocket.Upgrader

func realip(req *http.Request) string {
	addr := req.Header.Get("X-Forwarded-For")
	if addr == "" {
		addr = req.Header.Get("X-Real-IP")
	}
	return addr
}

// RegistWebSocket router
func RegistWebSocket(router *httprouter.Router, path string, ctl IWebSock) {
	connMap[path] = NewManager()
	router.GET(path, func(rw http.ResponseWriter, req *http.Request, ps httprouter.Params) {
		err := ctl.Auth(req)
		if err != nil {
			log.App.Error(err, realip(req))
			return
		}

		conn, err := upgrader.Upgrade(rw, req, nil)
		if err != nil {
			log.App.Error(err, realip(req))
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
		defer DoRecover(req)

		client.Listen(ctl)
	})
}

// DoRecover dillwith panic
func DoRecover(req *http.Request) {
	err := recover()
	if err != nil {
		log.App.Error(err, realip(req)).Stack()
	}
}
