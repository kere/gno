package httpd

import (
	"crypto/md5"
	"crypto/sha1"
	"fmt"
	"net/url"

	"github.com/kere/gno/libs/util"
	"github.com/valyala/bytebufferpool"
	"github.com/valyala/fasthttp"
)

func isPageOK(site *SiteServer, req *fasthttp.Request) bool {
	pToken := req.Header.Peek(APIFieldPageToken)
	bPath := req.Header.Referer()
	u, _ := url.Parse(string(bPath))

	pToken2 := buildToken([]byte(u.Path), site.Secret, site.Nonce)

	l := len(pToken)
	if l != len(pToken2) {
		return false
	}
	for i := 0; i < l; i++ {
		if pToken[i] != pToken2[i] {
			return false
		}
	}
	return true
}

func isAPIOK(req *fasthttp.Request, src []byte) bool {
	apiToken := req.Header.Peek(APIFieldToken)
	pToken := req.Header.Peek(APIFieldPageToken)
	u32 := buildAPIToken(req, src, pToken)
	// auth api token
	return u32 == util.Bytes2Str(apiToken)
}

// ts+method+ts+jsonStr + token;
func buildAPIToken(req *fasthttp.Request, src, pToken []byte) string {
	ts := req.Header.Peek(APIFieldTS)
	method := req.Header.Peek(APIFieldMethod)

	// method := req.PostArgs().Peek(APIFieldMethod)
	// method + ts + src + agent + ts + ptoken + hostname

	buf := bytebufferpool.Get()
	buf.Write(method)
	buf.Write(ts)
	buf.Write(src)
	buf.Write(req.Header.UserAgent())
	buf.Write(ts)
	buf.Write(pToken)

	origin := req.Header.Peek(HeadOrigin)
	u, err := url.Parse(string(origin))
	if err != nil {
		return ""
	}
	buf.WriteString(u.Hostname())

	// fmt.Println(buf.String())
	// req.Header.VisitAll(func(key []byte, value []byte) {
	// 	fmt.Println(string(key), string(value))
	// })

	str := fmt.Sprintf(stri16Formart, md5.Sum(buf.Bytes()))
	bytebufferpool.Put(buf)
	return str
}

// userAgent + hostname
func buildWSSign(req *fasthttp.Request) string {
	agent := req.Header.UserAgent()
	src := append(agent, req.URI().Host()...)
	return fmt.Sprintf(stri16Formart, md5.Sum(src))
}

// buildToken 生成 用户令牌
func buildToken(src []byte, sn, salt string) string {
	buf := bytebufferpool.Get()
	buf.Write(src)
	buf.WriteString(salt)
	buf.WriteString(sn)
	str := fmt.Sprintf(stri16Formart, sha1.Sum(buf.Bytes()))
	bytebufferpool.Put(buf)
	return str
}

// ts+blob.name +  blob.size + blob.lastModified + blob.type+navigator.userAgent+ts+ptoken + window.location.hostname;
func buildUploadToken(req *fasthttp.Request, name, size, last, typ, pToken []byte) string {
	ts := req.Header.Peek(APIFieldTS)
	// method + ts + src + agent + ts + ptoken + hostname
	buf := bytebufferpool.Get()
	buf.Write(ts)
	buf.Write(name)
	buf.Write(size)
	buf.Write(last)
	buf.Write(typ)
	buf.Write(req.Header.UserAgent())
	buf.Write(ts)
	buf.Write(pToken)

	origin := req.Header.Peek(HeadOrigin)
	u, err := url.Parse(string(origin))
	if err != nil {
		return ""
	}
	buf.WriteString(u.Hostname())

	str := fmt.Sprintf(stri16Formart, md5.Sum(buf.Bytes()))
	bytebufferpool.Put(buf)
	return str
}
