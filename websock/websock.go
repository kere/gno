package websock

import (
	"net/http"

	"github.com/gorilla/websocket"
	"github.com/kere/gno/libs/log"
	"github.com/kere/gno/libs/util"
)

// IWebSock interface
type IWebSock interface {
	Auth(r *http.Request) error
	Listen(conn *websocket.Conn)
	Exec(conn *websocket.Conn, args util.MapData) (interface{}, error)
}

// WebSock class
type WebSock struct {
	Target IWebSock
}

// Auth a
func (c WebSock) Auth(r *http.Request) error {
	return nil
}

// CurrentClientID genrate id
func CurrentClientID() int {
	counter++
	return counter
}

// Listen a
func (c WebSock) Listen(conn *websocket.Conn) {
	var args util.MapData
	var dat interface{}
	for {
		err := conn.ReadJSON(&args)
		if err != nil {
			log.App.Error("ReadJSON:", err)
			break
		}

		dat, err = c.Target.Exec(conn, args)
		if err != nil {
			log.App.Error("Exec:", err)
			break
		}

		err = conn.WriteJSON(dat)
		if err != nil {
			log.App.Error("WriteJSON:", err)
			break
		}
	}
}
