package httpd

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"time"

	"github.com/kere/gno/libs/cache"
	"github.com/valyala/fasthttp"
)

const (
	//CacheModeNone 不缓存页面
	CacheModeNone = 0
	//CacheModePage 按照页面名称缓存
	CacheModePage = 1
	//CacheModePagePath 按照URL Path缓存页面
	CacheModePagePath = 2
	//CacheModePageURI 按照URL缓存页面
	CacheModePageURI = 3

	pagecacheKeyPrefix = "c:"
	pageCacheSubfix    = ".htm"

	delim1 = byte('\n')
	delim2 = "\n"

	// HeaderEtag etag
	HeaderEtag = "ETag"
	// HeaderCacheCtl cache
	HeaderCacheCtl = "Cache-Control"
	// HeaderIfNoneMatch If-None-Match
	HeaderIfNoneMatch = "If-None-Match"
	// HeaderLastModified Last-Modified
	HeaderLastModified = "Last-Modified"
	// HeaderIfModifiedSince = "If-Modified-Since"
	HeaderIfModifiedSince = "If-Modified-Since"
)

// TryCache try to get cache
func TryCache(ctx *fasthttp.RequestCtx, p IPage) bool {
	if RunMode == ModeDev {
		return false
	}

	var src, srcTmp []byte
	var err error
	var last, key string

	mode, _, _ := p.CacheOption()

	switch mode {
	case CacheModePage:
		key = pagecacheKeyPrefix + p.Dir() + p.Name()

	case CacheModePagePath:
		key = pagecacheKeyPrefix + string(ctx.URI().Path())

	case CacheModePageURI:
		key = pagecacheKeyPrefix + string(ctx.URI().RequestURI())

	// case CacheModeFile:
	// 	filename := filepath.Join(WEBROOT, p.GetRequest().URL.Path+pageCacheSubfix)
	// 	src, err = ioutil.ReadFile(filename)

	default:
		return false
	}

	srcTmp, err = cache.GetBytes(key)
	if err != nil {
		return false
	}
	buf := bytes.NewBuffer(srcTmp)
	// 读取第一行时间戳
	last, err = buf.ReadString(delim1)
	if err != nil {
		return false
	}
	last = strings.TrimRight(last, delim2)
	src, err = ioutil.ReadAll(buf)
	if err != nil {
		return false
	}

	// check use cache ?
	// last0 := req.Header.Get(HeaderIfModifiedSince)
	last0 := string(ctx.Request.Header.Peek(HeaderIfModifiedSince))

	if last != "" && last == last0 {
		// w.WriteHeader(http.StatusNotModified)
		ctx.SetStatusCode(http.StatusNotModified)
		return true
	}

	// h.Set(HeaderEtag, last)
	// h.Set(HeaderIfNoneMatch, "true")
	setResponseHeader(ctx, p, last)

	// w.Write(src)
	ctx.Write(src)

	return true
}

// TrySetCache TrySet cache
func TrySetCache(ctx *fasthttp.RequestCtx, p IPage, buf *bytes.Buffer) error {
	var err error
	var last string
	mode, _, expires := p.CacheOption()

	switch mode {
	case CacheModePage:
		key := pagecacheKeyPrefix + p.Dir() + p.Name()
		last = time.Now().Format(time.RFC1123)
		src := append([]byte(last), byte('\n'))
		src = append(src, buf.Bytes()...)
		err = cache.Set(key, string(src), expires)

	case CacheModePagePath:
		key := pagecacheKeyPrefix + string(ctx.URI().Path())
		last = time.Now().Format(time.RFC1123)
		src := append([]byte(last), byte('\n'))
		src = append(src, buf.Bytes()...)
		err = cache.Set(key, string(src), expires)

		// case CacheModeFile:
		// 	filename := filepath.Join(WEBROOT, p.GetRequest().URL.Path+pageCacheSubfix)
		// 	err = os.MkdirAll(filepath.Dir(filename), os.ModeDir)
		// 	if err != nil {
		// 		return err
		// 	}
		//
		// 	f, err := os.OpenFile(filename, os.O_CREATE|os.O_TRUNC, os.ModePerm)
		// 	if err != nil {
		// 		return err
		// 	}
		// 	defer f.Close()
		// 	_, err = f.Write(buf.Bytes())
	}
	if err != nil {
		return err
	}

	setResponseHeader(ctx, p, last)

	return nil
}

func setResponseHeader(ctx *fasthttp.RequestCtx, p IPage, lastModified string) {
	_, expires, _ := p.CacheOption()
	// set response header
	// Cache-Control: public, max-age=3600
	if expires > 0 {
		ctx.Request.Header.Add(HeaderCacheCtl, "must-revalidate, max-age="+fmt.Sprint(expires))
	} else {
		ctx.Request.Header.Add(HeaderCacheCtl, "must-revalidate")
	}
	ctx.Request.Header.Set(HeaderLastModified, lastModified)
}

// ClearCache Clear cache
func ClearCache() error {
	if !cache.IsOK() {
		return nil
	}
	redis := cache.GetRedis()
	keys, err := redis.DoStrings("keys", pagecacheKeyPrefix+"*")
	if err != nil {
		return err
	}
	for i := range keys {
		redis.Do("del", keys[i])
	}
	return nil
}
