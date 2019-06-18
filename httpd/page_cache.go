package httpd

import (
	"bytes"
	"encoding/base64"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/kere/gno/libs/log"
	"github.com/valyala/fasthttp"
)

const ()

// PageCacheOption option
type PageCacheOption struct {
	PageMode     int // page, path, uri
	HTTPHead     int // 页面缓存模式 0:不缓存  1: etag  >1: 过期模式
	Store        int // 0: mem 1:file
	StoreExpires int
}

// pCacheElem class
type pCacheElem struct {
	LastModified string // Last-Modified: Fri, 12 May 2006 18:53:33 GMT
	Src          []byte
}

var pageCached = &sync.Map{}

// pageCachedKey key
func pageCachedKey(opt PageCacheOption, ctx *fasthttp.RequestCtx, p IPage) string {
	switch opt.PageMode {
	case CacheModePage:
		pdata := p.Data()
		return pagecacheKeyPrefix + pdata.Dir + pdata.Name

	case CacheModePagePath:
		return pagecacheKeyPrefix + string(ctx.URI().Path())

	case CacheModePageURI:
		return pagecacheKeyPrefix + string(ctx.URI().RequestURI())

	default:
		return ""
	}
}

// TryCache try to get cache
func TryCache(ctx *fasthttp.RequestCtx, p IPage) bool {
	if RunMode == ModeDev {
		return false
	}

	opt := p.Data().CacheOption
	key := pageCachedKey(opt, ctx, p)

	var src []byte
	var last string

	switch opt.Store {
	case CacheStoreMem:
		if v, isCached := pageCached.Load(key); isCached {
			pe := v.(pCacheElem)
			src = pe.Src
			last = pe.LastModified
		} else {
			return false
		}

	case CacheStoreFile:
		bstr, err := base64.StdEncoding.DecodeString(key)
		if err != nil {
			log.App.Alert(err)
			return false
		}
		name := filepath.Join(cacheFileStoreDir, string(bstr))
		file, err := os.OpenFile(name, os.O_RDONLY, os.ModePerm)
		if os.IsNotExist(err) {
			return false
		}
		src, _ = ioutil.ReadAll(file)
		stat, err := file.Stat()
		if err != nil {
			last = gmtNowTime(stat.ModTime())
		}
		file.Close()

	}

	// check use cache ?
	last0 := string(ctx.Request.Header.PeekBytes(HeaderIfModifiedSince))

	if last != "" && last == last0 {
		ctx.SetStatusCode(http.StatusNotModified)
		return true
	}

	setHeader(p, ctx, last)

	ctx.Write(src)
	return true
}

// TrySetCache TrySet cache
func TrySetCache(ctx *fasthttp.RequestCtx, p IPage, buf *bytes.Buffer) error {
	opt := p.Data().CacheOption
	key := pageCachedKey(opt, ctx, p)

	var last string

	switch opt.Store {
	case CacheStoreMem:
		last = gmtNowTime(time.Now())
		pageCached.Store(key, pCacheElem{LastModified: last, Src: buf.Bytes()})

	case CacheStoreFile:
		bstr, err := base64.StdEncoding.DecodeString(key)
		if err != nil {
			return err
		}

		name := filepath.Join(cacheFileStoreDir, string(bstr))
		f, err := os.OpenFile(name, os.O_CREATE|os.O_TRUNC, os.ModePerm)
		if err != nil {
			return err
		}
		_, err = f.Write(buf.Bytes())
		if err != nil {
			f.Close()
			return err
		}
		info, err := f.Stat()
		if err != nil {
			return err
		}
		last = gmtNowTime(info.ModTime())
		f.Close()

	}

	setHeader(p, ctx, last)
	return nil
}

// ClearCache Clear cache
func ClearCache() error {
	return nil
}
