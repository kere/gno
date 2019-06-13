package httpd

import (
	"crypto/md5"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"reflect"

	"github.com/kere/gno/libs/log"
	"github.com/kere/gno/libs/myerr"
	"github.com/kere/gno/libs/util"
	"github.com/valyala/fasthttp"
)

const (
	// ReplyTypeJSON reply json
	ReplyTypeJSON = 0
	// ReplyTypeText reply text
	ReplyTypeText = 1

	// APIFieldSrc post field
	APIFieldSrc = "_src"
	// APIFieldTS post field
	APIFieldTS = "Accts"
	// APIFieldMethod post field
	APIFieldMethod = "method"
	// APIFieldToken post field
	APIFieldToken = "Accto"
	// PageAccessTokenField 页面访问token的名称
	PageAccessTokenField = "accpt" //access page token
)

// IOpenAPI interface
type IOpenAPI interface {
	ReplyType() int
	Auth(ctx *fasthttp.RequestCtx) (err error)
}

// OpenAPIReply response
func OpenAPIReply(ctx *fasthttp.RequestCtx, api IOpenAPI, data interface{}) error {
	if data == nil {
		ctx.SetStatusCode(http.StatusOK)
		return nil
	}

	var src []byte
	var err error
	switch api.ReplyType() {
	case ReplyTypeJSON:
		src, err = json.Marshal(data)

	case ReplyTypeText:
		src = []byte(fmt.Sprint(data))

	}

	if err != nil {
		Site.Log.Warn(err)
		return err
	}

	ctx.Write(src)

	return nil
}

type openAPIExec func(args util.MapData) (interface{}, error)

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
		if name == "Auth" || name == "ReplyType" {
			continue
		}

		f := v.Method(i).Interface().(func(args util.MapData) (interface{}, error))
		openapiMap[rule+"/"+name] = openapiItem{Exec: f, API: openapi}

		s.Router.POST(rule+"/"+name, func(ctx *fasthttp.RequestCtx) {
			uri := string(ctx.URI().Path())
			item, isok := openapiMap[uri]
			if !isok {
				// doAPIError(errors.New(uri+" openapi not found"), rw, req)
				doAPIError(ctx, errors.New(uri+" openapi not found"))
				return
			}

			if RunMode == ModePro {
				defer func() {
					if p := recover(); p != nil {
						var err error
						str, isok := p.(string)
						if isok {
							err = errors.New(str)
						} else {
							err = errors.New("panic")
						}
						doAPIError(ctx, err)
					}
				}()
			}

			var args util.MapData
			pArgs := ctx.Request.PostArgs()
			src := pArgs.Peek(APIFieldSrc)

			if len(src) > 0 {
				err := json.Unmarshal(src, &args)
				if err != nil {
					doAPIError(ctx, myerr.New(err, string(src)))
					return
				}
			}

			err := authAPIToken(&ctx.Request, src)
			if err != nil {
				doAPIError(ctx, err)
				return
			}

			data, err := item.Exec(args)
			if err != nil {
				doAPIError(ctx, err)
				return
			}

			err = OpenAPIReply(ctx, item.API, data)
			if err != nil {
				doAPIError(ctx, err)
				return
			}

		})
	}
}

func doAPIError(ctx *fasthttp.RequestCtx, err error) {
	addr := ctx.RemoteAddr()
	log.App.Error(err, addr)
	ctx.Error(err.Error(), http.StatusInternalServerError)
}

func authAPIToken(req *fasthttp.Request, src []byte) error {
	token := req.Header.Peek(APIFieldToken)
	u32 := generateAPIToken(req, src)
	if string(u32) != string(token) {
		return errors.New("api token failed")
	}

	return nil
}

func generateAPIToken(req *fasthttp.Request, src []byte) string {
	ts := req.Header.Peek(APIFieldTS)
	method := req.PostArgs().Peek(APIFieldMethod)

	pageToken := req.Header.Cookie(PageAccessTokenField)

	// ts + method + ts + jsonStr + ptoken
	s := append([]byte{}, ts...)
	s = append(s, method...)
	s = append(s, ts...)
	s = append(s, src...)
	if len(pageToken) > 0 {
		b64, _ := base64.StdEncoding.DecodeString(string(pageToken))
		s = append(s, b64...)
	}

	return fmt.Sprintf("%x", md5.Sum(s))
}
