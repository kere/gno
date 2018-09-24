package websock

import (
	"fmt"
	"net/http"

	"github.com/gorilla/websocket"
	"github.com/kere/gno/libs/util"
)

// User class
type User struct {
}

// NewUser f
func NewUser() User {
	return User{}
}

// Auth f
func (u User) Auth(req *http.Request) error {
	return nil
}

// Exec f
func (u User) Exec(conn *websocket.Conn, args util.MapData) (interface{}, error) {
	fmt.Println("receive:", args)
	return util.MapData{"name": "luhan", "status": 99}, nil
}
