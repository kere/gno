package openapi

import (
	"fmt"
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/kere/gno"
	"github.com/kere/gno/libs/util"
)

// App class
type App struct {
	gno.OpenAPI
}

// NewApp func
func NewApp() App {
	return App{}
}

// Auth page auth
// if require is true then do auth
func (a App) Auth(req *http.Request, ps httprouter.Params) (bool, error) {
	require := true
	return require, nil
}

// PageData func
func (a App) PageData(req *http.Request, ps httprouter.Params, args util.MapData) (interface{}, error) {
	fmt.Println(args)

	return util.MapData{"isok": true}, nil
}