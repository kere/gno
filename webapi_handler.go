package gno

import (
	"encoding/json"
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/kere/gno/libs/log"
	"github.com/kere/gno/libs/util"
)

func doAPIHandle(webapi IWebAPI, rw http.ResponseWriter, req *http.Request, ps httprouter.Params) error {
	if isReq, err := webapi.Auth(req, ps); isReq && err != nil {
		return err
	}

	var args util.MapData
	str := req.PostFormValue(APIFieldSrc)
	src := []byte(str)
	if str != "" {
		err := json.Unmarshal(src, &args)
		if err != nil {
			return err
		}
	}

	if !webapi.IsSkipToken(req.Method) {
		err := authAPIToken(req, src)
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

func doAPIError(err error, rw http.ResponseWriter, req *http.Request) {
	addr := req.Header.Get("X-Forwarded-For")
	if addr == "" {
		addr = req.Header.Get("X-Real-IP")
	}
	log.App.Error(err, addr)
	rw.WriteHeader(http.StatusInternalServerError)
	rw.Write([]byte(err.Error()))
}
