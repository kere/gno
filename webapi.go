package gno

import (
	"encoding/json"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

// IWebAPI interface
type IWebAPI interface {
	Init(rw http.ResponseWriter, req *http.Request, params httprouter.Params)
	Auth() (require bool, err error)
	Exec() (interface{}, error)
	Reply(data interface{})
}

// WebAPI class
type WebAPI struct {
	Request        *http.Request
	Params         httprouter.Params
	ResponseWriter http.ResponseWriter
}

// Init api
func (w *WebAPI) Init(rw http.ResponseWriter, req *http.Request, params httprouter.Params) {
	w.ResponseWriter = rw
	w.Request = req
	w.Params = params
}

// Auth page auth
// if require is true then do auth
func (w *WebAPI) Auth() (require bool, err error) {
	return false, nil
}

// Exec api
func (w *WebAPI) Exec() (interface{}, error) {
	return nil, nil
}

// Reply response
func (w *WebAPI) Reply(data interface{}) {
	if data == nil {
		return
	}

	src, err := json.Marshal(data)
	if err != nil {
		Site.Log.Warn(err)
		return
	}

	w.ResponseWriter.Write(src)

	w.Request = nil
	w.Params = nil
	w.ResponseWriter = nil
}
