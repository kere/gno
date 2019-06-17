package openapi

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/kere/gno"
	"github.com/kere/gno/libs/util"
	"github.com/kere/gno/websock"
)

// App class
type App struct {
	gno.OpenAPI
}

// NewApp func
func NewApp() App {
	return App{}
}

// Prepare page auth
// if require is true then do auth
func (a App) Prepare(req *http.Request, ps httprouter.Params) (interface{}, error) {
	return nil, nil
}

// PageData func
func (a App) PageData(args util.MapData, dat interface{}) (interface{}, error) {
	// fmt.Println(args)

	return util.MapData{"isok": true}, nil
}

// ServerSend func
func (a App) ServerSend(args util.MapData, dat interface{}) (interface{}, error) {
	m := websock.GetManager("/ws")
	m.AllClients(func(c websock.Client) {
		c.Conn.WriteJSON(util.MapData{"isserver": true, "clientid": c.ID})
	})
	return util.MapData{"isok": true}, nil
}
