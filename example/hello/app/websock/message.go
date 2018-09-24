package websock

import (
	"fmt"
	"net/http"

	"github.com/gorilla/websocket"
	"github.com/kere/gno/libs/util"
)

// Message class
type Message struct {
}

// NewMessage f
func NewMessage() Message {
	return Message{}
}

// Auth f
func (u Message) Auth(req *http.Request) error {
	return nil
}

// Exec f
func (u Message) Exec(conn *websocket.Conn, args util.MapData) ([]byte, error) {
	fmt.Println("receive:", args.String("content"))
	return []byte("ni hao a..."), nil
}
