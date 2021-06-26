package httpd

import (
	"bytes"
	"crypto/md5"
	"encoding/base64"
	"errors"
	"fmt"

	"github.com/kere/gno/db"
	"github.com/kere/gno/libs/util"
	"github.com/valyala/bytebufferpool"
	"github.com/valyala/fasthttp"
)

var (
	//TableUsers users
	TableUsers = "users"
	// ErrDataNotFound err
	ErrDataNotFound = errors.New("data not found")
	// ErrLogin err
	ErrLogin = errors.New("auth failed")
	// ErrUserNotFound err
	ErrUserNotFound = errors.New("user not found")

	// FieldToken string
	FieldToken = "token"
	// FieldUserID id
	FieldUserID = "user_id"
	// FieldNick nick
	FieldNick = "nick"
	// FieldID id
	FieldID = "id"
)

const (
	// CookieMaxAge max cookie age
	CookieMaxAge = 3600 * 24 * 30
	// CookieUIID uiid
	CookieUIID = "_uiid"
	// CookieAccess _token
	CookieAccess = "_token"
	// CookieNick _nick
	CookieNick = "_nick"
)

func queryUserByNick(nick string) db.DBRow {
	// row, _ := db.NewQuery(TableUsers).Select("id,iid,nick,token,status").Where("nick=?", nick).QueryOne()
	q := db.Current().NewQuery(TableUsers)
	row, _ := q.Select("id,iid,nick,token,status").Where("nick=$1", nick).QueryOne()
	return row
}

// Auth page
func Auth(ctx *fasthttp.RequestCtx) error {
	nick := ctx.Request.Header.Cookie(CookieNick)
	uiid := ctx.Request.Header.Cookie(CookieUIID)
	val := ctx.Request.Header.Cookie(CookieAccess)
	if len(nick) == 0 {
		return ErrLogin
	}

	row := queryUserByNick(util.Bytes2Str(nick))
	if row.IsEmpty() {
		return ErrLogin
	}

	token := UnDasit(util.Bytes2Str(val))
	arr := bytes.Split(token, BDote)
	if len(arr) != 2 {
		return ErrLogin
	}

	dbToken := row.Bytes(FieldToken)

	if util.Bytes2Str(arr[1]) != accessToken(arr[0], uiid, dbToken) {
		return ErrLogin
	}
	ctx.SetUserValue(FieldUserID, row.Int(FieldID))
	ctx.SetUserValue(FieldNick, row.String(FieldNick))
	return nil
}

// LoginInOption login
type LoginInOption struct {
	TokenHTTPOnly bool
	CookieMaxAge  int
}

// DoLogin user
// uiid unique uid in client
func DoLogin(ctx *fasthttp.RequestCtx, nick string, srcb, signb []byte, opt LoginInOption) (string, error) {
	ts := ctx.Request.Header.Peek(APIFieldTS)
	pageToken := ctx.Request.Header.Peek(APIFieldPageToken)
	uiid := ctx.Request.Header.Cookie(CookieUIID)
	if len(uiid) == 0 {
		return "", ErrLogin
	}

	src := UnDasit(util.Bytes2Str(srcb))
	sign := UnDasit(util.Bytes2Str(signb))

	if len(src) == 0 {
		return "", ErrLogin
	}

	// 判断签名是否正确
	signNew := md5.Sum(append(src, ts...))
	if fmt.Sprintf("%x", signNew) != util.Bytes2Str(sign) {
		return "", ErrLogin
	}

	row := queryUserByNick(nick)
	if row.IsEmpty() {
		return "", ErrUserNotFound
	}

	dbToken := row.Bytes(FieldToken)
	tokenNew := authToken(ts, pageToken, uiid, dbToken)
	if tokenNew != util.Bytes2Str(src) {
		return "", ErrLogin
	}

	// expire := time.Now().AddDate(0, 1, 0)
	// set cookie access token
	accToken := accessToken(ts, uiid, dbToken)
	// println("addcookie:", string(ts), string(uiid), string(dbToken))

	acc := Dasit(util.Bytes2Str(ts) + "." + accToken)
	cook := fasthttp.Cookie{}
	cook.SetKey(CookieAccess)
	cook.SetValue(acc)
	cook.SetMaxAge(opt.CookieMaxAge)
	// cook.SetExpire(expire)
	if opt.TokenHTTPOnly {
		cook.SetHTTPOnly(true)
	}
	cook.SetPath("/")
	cook.SetSameSite(fasthttp.CookieSameSiteLaxMode)
	ctx.Response.Header.SetCookie(&cook)

	cook2 := fasthttp.Cookie{}
	cook2.SetKey(CookieNick)
	cook2.SetValue(nick)
	cook2.SetMaxAge(opt.CookieMaxAge)
	// cook.SetExpire(expire)
	// if opt.CookieHTTPOnly {
	// 	cook2.SetHTTPOnly(true)
	// }
	cook2.SetPath("/")
	cook2.SetSameSite(fasthttp.CookieSameSiteLaxMode)
	ctx.Response.Header.SetCookie(&cook2)

	return accToken, nil
}

func accessToken(ts, uiid, dbToken []byte) string {
	buf := bytebufferpool.Get()
	defer bytebufferpool.Put(buf)
	buf.Write(ts)
	buf.Write(util.Str2Bytes(Site.SiteData.Secret))
	buf.Write(dbToken)
	buf.Write(ts)
	buf.Write(uiid)

	b := md5.Sum(buf.Bytes())
	return fmt.Sprintf("%x", b)
}

func authToken(ts, pageToken, uiid, md5pwd []byte) string {
	// ts + md5(pwd) + ts + pageToken + uiid + ts
	buf := bytebufferpool.Get()
	defer bytebufferpool.Put(buf)
	buf.Write(ts)
	buf.Write(md5pwd)
	buf.Write(ts)
	buf.Write(pageToken)
	buf.Write(uiid)
	buf.Write(ts)

	b := md5.Sum(buf.Bytes())
	return fmt.Sprintf("%x", b)
}

// Dasit 加密
func Dasit(str string) string {
	if str == "" {
		return ""
	}
	b64 := make([]byte, base64.StdEncoding.EncodedLen(len(str)))
	base64.StdEncoding.Encode(b64, util.Str2Bytes(str))

	l := len(b64)
	for i := 0; i < l; i++ {
		v := b64[i]
		b64[i] = byte(int(v) ^ ((i % 7 << 4) + (i % 15)))
	}

	return base64.StdEncoding.EncodeToString(b64)
}

// UnDasit 解密
func UnDasit(s string) []byte {
	b, _ := base64.StdEncoding.DecodeString(s)
	l := len(b)
	for i := 0; i < l; i++ {
		b[i] = byte(int(b[i]) ^ ((i % 7 << 4) + (i % 15)))
	}
	dst := make([]byte, base64.StdEncoding.DecodedLen(l))
	base64.StdEncoding.Decode(dst, b)
	l = len(dst)
	n := 0
	for i := l - 1; i > -1; i-- {
		if dst[i] != 0 {
			break
		}
		n++
	}
	if n == 0 {
		return dst
	}
	return dst[:l-n]
}
