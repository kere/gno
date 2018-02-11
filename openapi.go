package gno

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"reflect"

	"github.com/julienschmidt/httprouter"
	"github.com/kere/gno/libs/util"
)

const (
	methodFieldName = "_method"
	srcField        = "_src"
)

// IOpenAPI interface
type IOpenAPI interface {
	Auth(req *http.Request, ps httprouter.Params) (require bool, err error)
	Reply(rw http.ResponseWriter, data interface{}) error
}

// OpenAPI class
type OpenAPI struct {
	ReplyType int //json, xml, text
}

// Auth page auth
// if require is true then do auth
func (w *OpenAPI) Auth(req *http.Request, ps httprouter.Params) (require bool, err error) {
	return require, nil
}

// Reply response
func (w *OpenAPI) Reply(rw http.ResponseWriter, data interface{}) error {
	if data == nil {
		rw.WriteHeader(http.StatusOK)
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

	// rw.WriteHeader(http.StatusOK)
	rw.Write(src)

	return nil
}

type openAPIExec func(req *http.Request, ps httprouter.Params, args util.MapData) (interface{}, error)

type openapiItem struct {
	Exec openAPIExec
	API  IOpenAPI
}

var openapiMap = make(map[string]openapiItem)

// RegistOpenAPI init open api
func (s *SiteServer) RegistOpenAPI(rule string, openapi IOpenAPI) {
	v := reflect.ValueOf(openapi)
	typ := v.Type()
	l := typ.NumMethod()
	for i := 0; i < l; i++ {
		m := typ.Method(i)
		name := m.Name
		if name == "Auth" || name == "Reply" {
			continue
		}
		f := v.Method(i).Interface().(func(req *http.Request, ps httprouter.Params, args util.MapData) (interface{}, error))
		openapiMap[rule+"/"+name] = openapiItem{Exec: f, API: openapi}

		fmt.Println("regist openapi:", rule+"/"+name)
	}

	// s.Router.GET(rule, doOpenAPIHandle)
	s.Router.POST(rule+"/:"+methodFieldName, doOpenAPIHandle)
}

func doOpenAPIHandle(rw http.ResponseWriter, req *http.Request, ps httprouter.Params) {
	uri := req.URL.RequestURI()
	// method := ps.ByName(methodFieldName)

	var item openapiItem
	var isok bool
	if item, isok = openapiMap[uri]; !isok {
		doAPIError(errors.New(uri+" openapi not found"), rw)
		return
	}

	if isReq, err := item.API.Auth(req, ps); isReq && err != nil {
		doAPIError(err, rw)
		return
	}

	var args util.MapData
	src := req.PostFormValue(srcField)
	if src != "" {
		err := json.Unmarshal([]byte(src), &args)
		if err != nil {
			doAPIError(err, rw)
			return
		}
	}

	data, err := item.Exec(req, ps, args)
	if err != nil {
		doAPIError(err, rw)
		return
	}

	err = item.API.Reply(rw, data)
	if err != nil {
		doAPIError(err, rw)
		return
	}
}
