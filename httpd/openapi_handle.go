package httpd

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/kere/gno/libs/log"
	"github.com/kere/gno/libs/util"
	"github.com/valyala/fasthttp"
)

// ErrNoAPI err
var ErrNoAPI = errors.New("openapi not found")

func openAPIHandle(ctx *fasthttp.RequestCtx) {
	uri := util.Bytes2Str(ctx.URI().Path())

	itemExec, isok := openapiMap[uri]
	if !isok {
		Site.Log.Error(uri)
		doAPIError(ctx, ErrNoAPI)
		return
	}

	if RunMode == ModePro {
		defer func() {
			if p := recover(); p != nil {
				str, isok := p.(string)
				var err error
				if isok {
					err = errors.New(str)
				} else {
					err = errors.New("panic")
				}
				log.App.Warn(err)
			}
		}()
	}
	// Auth

	// application/x-www-form-urlencoded;charset=UTF-8
	// application/json;charset=UTF-8
	src := ctx.Request.Body()

	var params util.MapData
	if len(src) > 0 {
		if err := json.Unmarshal(src, &params); err != nil {
			doAPIError(ctx, err)
			return
		}
	}

	if !isAPIOK(&ctx.Request, src) {
		doAPIError(ctx, errors.New("api auth failed"))
		return
	}

	data, err := itemExec(ctx, params)
	if err != nil {
		doAPIError(ctx, err)
		return
	}

	if data == nil {
		ctx.SetStatusCode(http.StatusOK)
		return
	}

	result, err := json.Marshal(data)
	if err != nil {
		Site.Log.Warn(err)
		doAPIError(ctx, err)
		return
	}

	_, err = ctx.Write(result)
	if err != nil {
		Site.Log.Error(err)
		doAPIError(ctx, err)
		return
	}
}

// // OpenAPIReply response
// func OpenAPIReply(ctx *fasthttp.RequestCtx, data interface{}) error {
// 	return nil
// }
