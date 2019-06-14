package openapi

import (
	"github.com/kere/gno/libs/util"
	"github.com/kere/gno/websock"
	"github.com/valyala/fasthttp"
)

// App class
type App struct {
}

// NewApp func
func NewApp() *App {
	return &App{}
}

// Auth page auth
// if require is true then do auth
func (a App) Auth(ctx *fasthttp.RequestCtx) error {
	return nil
}

// PageData func
func (a App) PageData(ctx *fasthttp.RequestCtx, args util.MapData) (interface{}, error) {
	// fmt.Println(args)

	return util.MapData{"isok": true}, nil
}

// ServerSend func
func (a App) ServerSend(ctx *fasthttp.RequestCtx, args util.MapData) (interface{}, error) {
	m := websock.GetManager("/ws")
	m.AllClients(func(c websock.Client) {
		c.Conn.WriteJSON(util.MapData{"isserver": true, "clientid": c.ID})
	})
	return util.MapData{"isok": true}, nil
}
