package gno

import (
	"bytes"
	"crypto/md5"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/julienschmidt/httprouter"
	"github.com/kere/gno/libs/util"
)

const (
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
	Auth(req *http.Request, ps httprouter.Params) (require bool, err error)
	Reply(rw http.ResponseWriter, data interface{}) error
}

// OpenAPI class
type OpenAPI struct {
	ReplyType int //json, xml, text
}

// Auth page auth
// if require is true then do auth
func (w OpenAPI) Auth(req *http.Request, ps httprouter.Params) (require bool, err error) {
	return require, nil
}

// Reply response
func (w OpenAPI) Reply(rw http.ResponseWriter, data interface{}) error {
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

func generateAPIToken(req *http.Request, src []byte) string {
	ts := req.Header.Get(APIFieldTS)

	method := req.PostFormValue(APIFieldMethod)
	ptoken := ""
	c, err := req.Cookie(PageAccessTokenField)
	if err == nil {
		ptoken, _ = url.PathUnescape(c.Value)
	}

	// ts + method + jsonStr + ptoken
	s := append([]byte(ts+method+ts), src...)
	if ptoken != "" {
		b64, _ := base64.StdEncoding.DecodeString(ptoken)
		s = append(s, b64...)
	}

	return fmt.Sprintf("%x", md5.Sum(s))
}

func authAPIToken(req *http.Request, src []byte) error {
	token := req.Header.Get(APIFieldToken)
	u32 := generateAPIToken(req, src)
	if u32 != token {
		return errors.New("api token failed")
	}

	return nil
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
