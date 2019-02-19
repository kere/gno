package gno

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/kere/gno/libs/util"
)

func webAPIHandle(webapi IWebAPI, rw http.ResponseWriter, req *http.Request, ps httprouter.Params) (err error) {
	if isReq, err := webapi.Auth(req, ps); isReq && err != nil {
		return err
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

	var args util.MapData
	str := req.PostFormValue(APIFieldSrc)
	src := []byte(str)
	if str != "" {
		err = json.Unmarshal(src, &args)
		if err != nil {
			return err
		}
	}

	if !webapi.IsSkipToken(req.Method) {
		err = authAPIToken(req, src)
		if err != nil {
			return err
		}
	}

	data, err := webapi.Exec(req, ps, args)
	if err != nil {
		return err
	}

	return webapi.Reply(rw, data)
}
