package gno

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/kere/gno/libs/myerr"
	"github.com/kere/gno/libs/util"
)

func openAPIHandle(rw http.ResponseWriter, req *http.Request, ps httprouter.Params) (err error) {
	uri := req.URL.Path

	var item openapiItem
	var isok bool
	if item, isok = openapiMap[uri]; !isok {
		// doAPIError(errors.New(uri+" openapi not found"), rw, req)
		return errors.New(uri + " openapi not found")
	}

	if RunMode == ModePro {
		defer func() {
			if p := recover(); p != nil {
				str, ok := p.(string)
				if ok {
					err = errors.New(str)
				} else {
					err = errors.New("panic")
				}
			}
		}()
	}

	prepareDat, err := item.API.Prepare(req, ps)
	if err != nil {
		return err
	}

	var args util.MapData
	str := req.PostFormValue(APIFieldSrc)
	src := []byte(str)
	if str != "" {
		err := json.Unmarshal(src, &args)
		if err != nil {
			return myerr.New(err, str)
		}
	}

	err = authAPIToken(req, src)
	if err != nil {
		return err
	}

	data, err := item.Exec(args, prepareDat)
	if err != nil {
		return err
	}

	return item.API.Reply(rw, data)
}
