package httpd

import (
	"bytes"
	"crypto/md5"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"reflect"
	"time"

	"github.com/kere/gno/libs/log"
	"github.com/kere/gno/libs/util"
	"github.com/valyala/bytebufferpool"
	"github.com/valyala/fasthttp"
)

var openapiMap = make(map[string]apiExec)

const (
	// ReplyTypeJSON reply json
	ReplyTypeJSON = 0
	// ReplyTypeText reply text
	ReplyTypeText = 1

	// // APIFieldSrc post field
	// APIFieldSrc = "_src"

	// APIFieldTS post field
	APIFieldTS = "Accts"
	// APIFieldMethod post field
	APIFieldMethod = "Api"

	// APIFieldToken post field
	APIFieldToken = "Accto"
	// APIFieldPageToken post field
	APIFieldPageToken = "AccPage"

	// PageAccessTokenField 页面访问token的名称
	PageAccessTokenField = "accpt" //access page token

	// HeadOrigin http head
	HeadOrigin = "Origin"

	methodAuth  = "Auth"
	methodReply = "ReplyType"
)

// IOpenAPI interface
type IOpenAPI interface {
	Auth(ctx *fasthttp.RequestCtx) (err error)
}

type apiExec func(ctx *fasthttp.RequestCtx, args util.MapData) (interface{}, error)

// RegistOpenAPI init open api
func (s *SiteServer) RegistOpenAPI(rule string, openapi IOpenAPI) {
	v := reflect.ValueOf(openapi)
	typ := v.Type()
	l := typ.NumMethod()
	for i := 0; i < l; i++ {
		m := typ.Method(i)
		if m.Name == methodAuth || m.Name == methodReply {
			continue
		}

		openapiMap[rule+Slash+m.Name] = v.Method(i).Interface().(func(ctx *fasthttp.RequestCtx, args util.MapData) (interface{}, error))

		s.Router.POST(rule+Slash+m.Name, func(ctx *fasthttp.RequestCtx) {
			// if s.Pool == 0 {
			openAPIHandle(ctx)
			// 	return
			// }
			// done := make(chan struct{})
			// if err := pool.Invoke(PoolParams{Typ: invokeAPI, Ctx: ctx, Done: done}); err != nil {
			// 	doAPIError(ctx, errors.New("Throttle limit error"))
			// }
			// // <-done
			// select {
			// case <-done:
			// case <-time.After(s.Timeout):
			// 	ctx.TimeoutError("timeout!")
			// }

		})
	}
}

func doAPIError(ctx *fasthttp.RequestCtx, err error) {
	addr := ctx.RemoteAddr()
	log.App.Error(err, addr)
	ctx.Error(err.Error(), http.StatusInternalServerError)
}

const (
	clientAgent = "Go-http-client"
)

// SendAPI send api method
func SendAPI(uri string, method string, dat util.MapData) (util.MapData, error) {
	// str = ts+method+ts+jsonStr+ts;
	src, err := json.Marshal(dat)
	if err != nil {
		return nil, err
	}

	// vals := url.Values{}
	// vals.Add(APIFieldSrc, string(src))
	// vals.Add(APIFieldMethod, method)

	reader := bytes.NewReader(src)
	req, err := http.NewRequest(http.MethodPost, uri+"/"+method, reader)
	if err != nil {
		return nil, err
	}
	req.Header.Set(APIFieldMethod, method)
	req.Header.Set("Content-Type", "application/json;charset=UTF-8")
	req.Header.Set("User-Agent", clientAgent)

	// write head token
	// method + ts + src + agent + ts + ptoken + hostname
	ts := fmt.Sprint(time.Now().Unix())
	buf := bytebufferpool.Get()
	buf.WriteString(method)
	buf.WriteString(ts)
	buf.Write(src)
	buf.WriteString(clientAgent)
	buf.WriteString(ts)

	// u, err := url.Parse(uri)
	// if err != nil {
	// 	return nil, err
	// }
	// buf.WriteString(u.Hostname())

	// fmt.Println("Client:", buf.String())
	token := fmt.Sprintf("%x", md5.Sum(buf.Bytes()))
	bytebufferpool.Put(buf)

	req.Header.Set(APIFieldTS, ts)
	req.Header.Set(APIFieldToken, token)

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
