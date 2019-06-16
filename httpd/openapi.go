package httpd

import (
	"bytes"
	"crypto/md5"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"reflect"
	"strings"
	"time"

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
	// APIFieldPageToken post field
	APIFieldPageToken = "AccPage"

	// PageAccessTokenField 页面访问token的名称
	PageAccessTokenField = "accpt" //access page token

	// HeadOrigin http head
	HeadOrigin = "Origin"
)

// IOpenAPI interface
type IOpenAPI interface {
	Auth(ctx *fasthttp.RequestCtx) (err error)
}

// OpenAPIReply response
func OpenAPIReply(ctx *fasthttp.RequestCtx, data interface{}) error {
	if data == nil {
		ctx.SetStatusCode(http.StatusOK)
		return nil
	}
	src, err := json.Marshal(data)
	if err != nil {
		Site.Log.Warn(err)
		return err
	}

	ctx.Write(src)

	return nil
}

type apiExec func(ctx *fasthttp.RequestCtx, args util.MapData) (interface{}, error)

var openapiMap = make(map[string]apiExec)

// RegistOpenAPI init open api
func (s *SiteServer) RegistOpenAPI(rule string, openapi IOpenAPI) {
	v := reflect.ValueOf(openapi)
	typ := v.Type()
	l := typ.NumMethod()
	for i := 0; i < l; i++ {
		m := typ.Method(i)
		if m.Name == "Auth" || m.Name == "ReplyType" {
			continue
		}

		openapiMap[rule+"/"+m.Name] = v.Method(i).Interface().(func(ctx *fasthttp.RequestCtx, args util.MapData) (interface{}, error))

		s.Router.POST(rule+"/"+m.Name, func(ctx *fasthttp.RequestCtx) {
			uri := string(ctx.URI().Path())
			itemExec, isok := openapiMap[uri]
			if !isok {
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

			pArgs := ctx.Request.PostArgs()
			src := pArgs.Peek(APIFieldSrc)

			var params util.MapData
			if len(src) > 0 {
				err := json.Unmarshal(src, &params)
				if err != nil {
					doAPIError(ctx, myerr.New(err, string(src)))
					return
				}
			}

			if !isAPIOK(s, &ctx.Request, src) {
				doAPIError(ctx, errors.New("api auth failed"))
				return
			}

			data, err := itemExec(ctx, params)
			if err != nil {
				doAPIError(ctx, err)
				return
			}

			err = OpenAPIReply(ctx, data)
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

// SendAPI send api method
func SendAPI(uri string, method string, dat util.MapData) (util.MapData, error) {
	// data:       {'_src': jsonStr, 'now': now, 'token': md5(str), 'method': method},
	// str = now+method+now+jsonStr+now;
	src, err := json.Marshal(dat)
	if err != nil {
		return nil, err
	}

	vals := url.Values{}
	vals.Add(APIFieldSrc, string(src))
	vals.Add(APIFieldMethod, method)

	// ts+method+jsonStr + token;
	ts := fmt.Sprint(time.Now().Unix())

	buf := bytes.NewBufferString(ts + method + ts)
	buf.Write(src)
	token := fmt.Sprintf("%x", md5.Sum(buf.Bytes()))

	req, err := http.NewRequest(http.MethodPost, uri+"/"+method, strings.NewReader(vals.Encode()))
	if err != nil {
		return nil, err
	}
	req.Header.Set(APIFieldTS, ts)
	req.Header.Set(APIFieldToken, token)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	client := &http.Client{}

	resq, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	defer resq.Body.Close()

	body, err := ioutil.ReadAll(resq.Body)
	if err != nil {
		return nil, err
	}

	if resq.StatusCode != http.StatusOK {
		return nil, errors.New(string(body) + " " + uri + "/" + method)
	}

	var obj util.MapData
	err = json.Unmarshal(body, &obj)

	return obj, err
}
