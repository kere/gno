package websock

import (
	"fmt"
	"net/http"

	"github.com/gorilla/websocket"
	"github.com/kere/gno/libs/util"
	"github.com/kere/gno/websock"
)

// User class
type User struct {
	websock.WebSock
}

// NewUser f
func NewUser() User {
	w := User{}
	w.Target = w
	return w
}

// Auth f
func (u User) Auth(req *http.Request) error {
	return nil
}

// Exec f
func (u User) Exec(conn *websocket.Conn, args util.MapData) (interface{}, error) {
	fmt.Println("receive: obj:", args)
	return util.MapData{"name": "luhan", "status": 99}, nil
}
