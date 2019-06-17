package httpd

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/kere/gno/libs/log"
	"github.com/kere/gno/libs/util"
	"github.com/valyala/fasthttp"
)

func openAPIHandle(ctx *fasthttp.RequestCtx) error {
	uri := string(ctx.URI().Path())
	itemExec, isok := openapiMap[uri]
	if !isok {
		return errors.New(uri + " openapi not found")
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

	pArgs := ctx.Request.PostArgs()
	src := pArgs.Peek(APIFieldSrc)

	var params util.MapData
	if len(src) > 0 {
		err := json.Unmarshal(src, &params)
		if err != nil {
			return err
		}
	}

	if !isAPIOK(&ctx.Request, src) {
		return errors.New("api auth failed")
	}

	data, err := itemExec(ctx, params)
	if err != nil {
		return err
	}

	if data == nil {
		ctx.SetStatusCode(http.StatusOK)
		return nil
	}

	result, err := json.Marshal(data)
	if err != nil {
		Site.Log.Warn(err)
		return err
	}

	_, err = ctx.Write(result)
	return err
}

// OpenAPIReply response
func OpenAPIReply(ctx *fasthttp.RequestCtx, data interface{}) error {

	return nil
}
