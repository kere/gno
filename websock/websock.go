package websock

import (
	"net/http"

	"github.com/gorilla/websocket"
	"github.com/kere/gno/libs/util"
)

// IWebSock interface
type IWebSock interface {
	Auth(r *http.Request) error
	Exec(conn *websocket.Conn, args util.MapData) (interface{}, error)
}
