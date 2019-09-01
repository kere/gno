package openapi

import (
	"time"

	"github.com/kere/gno/libs/util"
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
	time.Sleep(time.Second)
	return util.MapData{"isok": true}, nil
}
