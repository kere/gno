package websock

import (
	"net/http"

	"github.com/gorilla/websocket"
	"github.com/kere/gno/libs/log"
)

// IMessageSock interface
type IMessageSock interface {
	Auth(r *http.Request) error
	ClientID() int
	// Init(conn *websocket.Conn) IMessageSock
	Listen(conn *websocket.Conn)
	Exec(conn *websocket.Conn, msg []byte) ([]byte, error)
}

// MessageSock class
type MessageSock struct {
	ID     int
	Conn   *websocket.Conn
	Target IMessageSock
}

// // Init f
// func (c MessageSock) Init(target IMessageSock, conn *websocket.Conn) IMessageSock {
// 	counter++
// 	c.Target = target
// 	c.ID = counter
// 	c.Conn = conn
// 	return c
// }

// Auth a
func (c MessageSock) Auth(r *http.Request) error {
	return nil
}

// Exec msg
func (c MessageSock) Exec(conn *websocket.Conn, msg []byte) ([]byte, error) {
	return nil, nil
}

// ClientID a
func (c MessageSock) ClientID() int {
	return c.ID
}

var counter int

// Listen a
func (c MessageSock) Listen(conn *websocket.Conn) {
	counter++
	c.ID = counter
	for {
		mt, message, err := conn.ReadMessage()
		if err != nil {
			log.App.Error(err)
			break
		}
		message, err = c.Target.Exec(conn, message)
		if err != nil {
			log.App.Error(err)
			break
		}
		err = conn.WriteMessage(mt, message)
		if err != nil {
			log.App.Error(err)
			break
		}
	}
}
