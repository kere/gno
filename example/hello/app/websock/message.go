package websock

import (
	"fmt"
	"net/http"

	"github.com/gorilla/websocket"
	"github.com/kere/gno/websock"
)

// Message class
type Message struct {
	websock.MessageSock
}

// NewMessage f
func NewMessage() Message {
	w := Message{}
	w.Target = w
	return w
}

// Auth f
func (u Message) Auth(req *http.Request) error {
	return nil
}

// Exec f
func (u Message) Exec(conn *websocket.Conn, msg []byte) ([]byte, error) {
	fmt.Println("receive:", string(msg))
	return []byte("ni hao a..."), nil
}
