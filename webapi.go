package gno

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/kere/gno/libs/util"
)

const (
	// ReplyTypeJSON reply json
	ReplyTypeJSON = 0
	// ReplyTypeText reply text
	ReplyTypeText = 1
	// ReplyTypeXML reply xml
	ReplyTypeXML = 2
)

// IWebAPI interface
type IWebAPI interface {
	// Init(rw http.ResponseWriter, req *http.Request, params httprouter.Params) (IWebAPI, error)
	Auth(req *http.Request, params httprouter.Params) (require bool, err error)
	Exec(args util.MapData) (interface{}, error)
	Reply(data interface{}) error
	IsSkipToken(string) bool
}

// WebAPI class
type WebAPI struct {
	Request        *http.Request
	Params         httprouter.Params
	ResponseWriter http.ResponseWriter
	ReplyType      int //json, xml, text

	IsSkipTokenMethodGet  bool // 忽略token检查
	IsSkipTokenMethodPost bool
}

// Init api
// func (w WebAPI) Init() (IWebAPI, error) {
// func (w WebAPI) Init(rw http.ResponseWriter, req *http.Request, params httprouter.Params) (IWebAPI, error) {
// w.ResponseWriter = rw
// w.Request = req
// w.Params = params
// return w, nil
// }

// IsSkipToken 是否忽略token 检查
func (w WebAPI) IsSkipToken(method string) bool {
	if method == http.MethodGet {
		return w.IsSkipTokenMethodGet
	}
	return w.IsSkipTokenMethodPost
}

// Auth page auth
// if require is true then do auth
func (w WebAPI) Auth(req *http.Request, params httprouter.Params) (require bool, err error) {
	return require, nil
}

// Exec api
func (w WebAPI) Exec(args util.MapData) (interface{}, error) {
	return nil, nil
}

// Reply response
func (w WebAPI) Reply(data interface{}) error {
	if data == nil {
		w.ResponseWriter.WriteHeader(http.StatusOK)
		return nil
	}

	var src []byte
	var err error
	switch w.ReplyType {
	case ReplyTypeText:
		src = []byte(fmt.Sprint(data))
	case ReplyTypeJSON:
		src, err = json.Marshal(data)
	case ReplyTypeXML:

	}
	if err != nil {
		Site.Log.Warn(err)
		return err
	}

	w.ResponseWriter.WriteHeader(http.StatusOK)
	w.ResponseWriter.Write(src)

	// w.Request = nil
	// w.Params = nil
	// w.ResponseWriter = nil
	return nil
}
