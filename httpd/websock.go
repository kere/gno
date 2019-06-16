package httpd

import (
	"fmt"
	"net/http"
	"reflect"

	"github.com/fasthttp/websocket"
	"github.com/kere/gno/libs/log"
	"github.com/kere/gno/libs/util"
	"github.com/valyala/fasthttp"
)

const (
	errorMethodNotFound = "method not found"
	wsSignField         = "sign"
)

var wsMethodMap = make(map[string]wsExec)

// IWebSock interface
type IWebSock interface {
	Auth(ctx *fasthttp.RequestCtx) error
}

type wsExec func(args util.MapData) (interface{}, error)

func buildWSExec(w IWebSock) {
	v := reflect.ValueOf(w)
	typ := v.Type()
	l := typ.NumMethod()
	clasName := typ.Name()
	for i := 0; i < l; i++ {
		m := typ.Method(i)
		if m.Name == "Auth" {
			continue
		}

		wsMethodMap[clasName+m.Name] = v.Method(i).Interface().(func(args util.MapData) (interface{}, error))
	}
}

type messageRecv struct {
	Method string       `json:"method"`
	Args   util.MapData `json:"args"`
	Result interface{}  `json:"result"`
}

// RegistWS router
func (s *SiteServer) RegistWS(rule string, w IWebSock) {
	buildWSExec(w)
	var upgrader = websocket.FastHTTPUpgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
	}

	s.Router.GET(rule, func(ctx *fasthttp.RequestCtx) {
		sign := ctx.QueryArgs().Peek(wsSignField)
		sign2 := buildWSSign(&ctx.Request)
		if string(sign) != sign2 {
			log.App.Debug("ws sign failed:", string(sign), sign2)
			return
		}

		if err := w.Auth(ctx); err != nil {
			ctx.SetStatusCode(http.StatusForbidden)
			return
		}

		err := upgrader.Upgrade(ctx, func(ws *websocket.Conn) {
			defer ws.Close()
			for {
				var recv messageRecv
				if err := ws.ReadJSON(&recv); err != nil {
					break
				}
				wsExec, isok := wsMethodMap[recv.Method]
				if !isok {
					break
				}

				dat, err := wsExec(recv.Args)
				if err != nil {
					if errW := ws.WriteJSON(map[string]interface{}{"iserror": true, "error": err.Error()}); errW != nil {
						break
					}
					continue
				}

				if dat != nil {
					recv.Result = dat
					if errW := ws.WriteJSON(recv); errW != nil {
						break
					}
				}
			}
		}) // Upgrade end

		if err != nil {
			if _, ok := err.(websocket.HandshakeError); ok {
				fmt.Println(err)
			}
			return
		}

	})
}
