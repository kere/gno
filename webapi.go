package gno

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

const (
	// ReplyTypeJSON reply json
	ReplyTypeJSON = iota
	// ReplyTypeXML reply xml
	ReplyTypeXML
	// ReplyTypeText reply text
	ReplyTypeText
)

// IWebAPI interface
type IWebAPI interface {
	Init(rw http.ResponseWriter, req *http.Request, params httprouter.Params)
	Auth() (require bool, err error)
	Exec() (interface{}, error)
	Reply(data interface{}) error
}

// WebAPI class
type WebAPI struct {
	Request        *http.Request
	Params         httprouter.Params
	ResponseWriter http.ResponseWriter
	ReplyType      int //json, xml, text
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
	return require, nil
}

// Exec api
func (w *WebAPI) Exec() (interface{}, error) {
	return nil, nil
}

// Reply response
func (w *WebAPI) Reply(data interface{}) error {
	if data == nil {
		w.ResponseWriter.WriteHeader(http.StatusOK)
		return nil
	}

	var src []byte
	var err error
	switch w.ReplyType {
	case ReplyTypeJSON:
		src, err = json.Marshal(data)
	case ReplyTypeText:
		src = []byte(fmt.Sprint(data))
	case ReplyTypeXML:

	}
	if err != nil {
		Site.Log.Warn(err)
		return err
	}

	w.ResponseWriter.WriteHeader(http.StatusOK)
	w.ResponseWriter.Write(src)

	w.Request = nil
	w.Params = nil
	w.ResponseWriter = nil
	return nil
}
