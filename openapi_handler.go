package gno

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/kere/gno/libs/util"
)

func doOpenAPIHandle(rw http.ResponseWriter, req *http.Request, ps httprouter.Params) {
	uri := req.URL.Path

	var item openapiItem
	var isok bool
	if item, isok = openapiMap[uri]; !isok {
		doAPIError(errors.New(uri+" openapi not found"), rw, req)
		return
	}

	if isReq, err := item.API.Auth(req, ps); isReq && err != nil {
		doAPIError(err, rw, req)
		return
	}

	var args util.MapData
	str := req.PostFormValue(APIFieldSrc)
	src := []byte(str)
	if str != "" {
		err := json.Unmarshal(src, &args)
		if err != nil {
			doAPIError(err, rw, req)
			return
		}
	}

	err := authAPIToken(req, src)
	if err != nil {
		doAPIError(err, rw, req)
		return
	}

	data, err := item.Exec(req, ps, args)
	if err != nil {
		doAPIError(err, rw, req)
		return
	}

	err = item.API.Reply(rw, data)
	if err != nil {
		doAPIError(err, rw, req)
		return
	}
}
