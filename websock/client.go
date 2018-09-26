package websock

import (
	"net/http"
	"net/url"

	"github.com/gorilla/websocket"
	"github.com/kere/gno/libs/log"
	"github.com/kere/gno/libs/util"
)

// Client class
type Client struct {
	ID      int
	Conn    *websocket.Conn
	Cookies []*http.Cookie
	Form    url.Values
	Header  http.Header
	IP      string
}

// Cookie get connect cookie
func (c Client) Cookie(name string) (*http.Cookie, error) {
	for _, cookie := range c.Cookies {
		if cookie.Name == name {
			return cookie, nil
		}
	}
	return nil, http.ErrNoCookie
}

// Close a
func (c *Client) Close() {
	c.Conn.Close()
	c.Cookies = nil
	c.Form = nil
}

// Listen a
func (c *Client) Listen(ctl IWebSock) {
	var args util.MapData
	var dat interface{}
	conn := c.Conn

	for {
		err := conn.ReadJSON(&args)
		if websocket.IsCloseError(err, 1001) {
			// Listen 方法结束后，会自动清理当前client
			break
		}
		if err != nil {
			log.App.Error("ReadJSON:", err)
			break
		}

		dat, err = ctl.Exec(conn, args)
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
