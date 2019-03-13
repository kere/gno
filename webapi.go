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
	Auth(req *http.Request, params httprouter.Params) error
	Exec(req *http.Request, params httprouter.Params, args util.MapData) (interface{}, error)
	Reply(rw http.ResponseWriter, data interface{}) error
	IsSkipToken(string) bool
}

// WebAPI class
type WebAPI struct {
	ReplyType int //json, xml, text

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
func (w WebAPI) Auth(req *http.Request, params httprouter.Params) error {
	return nil
}

// Exec api
func (w WebAPI) Exec(req *http.Request, params httprouter.Params, args util.MapData) (interface{}, error) {
	return nil, nil
}

// Reply response
func (w WebAPI) Reply(rw http.ResponseWriter, data interface{}) error {
	if data == nil {
		rw.WriteHeader(http.StatusOK)
		return nil
	}

	var src []byte
	var err error
	switch w.ReplyType {
	default:
		src = []byte(fmt.Sprint(data))
	case ReplyTypeJSON:
		src, err = json.Marshal(data)
	case ReplyTypeXML:

	}
	if err != nil {
		Site.Log.Warn(err)
		return err
	}

	rw.WriteHeader(http.StatusOK)
	rw.Write(src)

	return nil
}
