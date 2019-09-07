package httpd

import (
	"fmt"
	"hash/crc32"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/kere/gno/libs/log"
	"github.com/kere/gno/libs/util"
	"github.com/valyala/fasthttp"
)

// PageCacheOption option
type PageCacheOption struct {
	PageMode int // page, path, uri
	HTTPHead int // 页面缓存模式 0:不缓存  1: etag  >1: 过期模式
	Store    int // 0: mem 1:file
}

// pCache class
type pCache struct {
	LastModified string // Last-Modified: Fri, 12 May 2006 18:53:33 GMT
	Src          []byte
}

var pageCacheMap = &sync.Map{}

// pageCachedKey key
func pageCachedKey(opt PageCacheOption, ctx *fasthttp.RequestCtx, p IPage) []byte {
	var src []byte
	switch opt.PageMode {
	case CacheModePage:
		pdata := p.Data()
		src = []byte(pdata.Dir + pdata.Name)

	case CacheModePageURI:
		src = ctx.URI().RequestURI()

	default:
		src = ctx.URI().Path()
	}
	ieee := crc32.NewIEEE()
	// io.WriteString(ieee, str)
	ieee.Write(src)
	v64 := uint64(ieee.Sum32())
	return util.IntZipTo62(v64)
}

// DisablePageCache bool
var DisablePageCache = true

// TryCache try to get cache
func TryCache(ctx *fasthttp.RequestCtx, p IPage) bool {
	if DisablePageCache {
		return false
	}

	opt := p.Data().CacheOption
	if opt.Store == CacheStoreNone {
		return false
	}

	key := pageCachedKey(opt, ctx, p)

	var src []byte
	var last string

	switch opt.Store {
	case CacheStoreMem:
		v, isCached := pageCacheMap.Load(string(key))
		if !isCached {
			return false
		}
		pe := v.(pCache)
		src = pe.Src
		last = pe.LastModified

	case CacheStoreFile:
		name := filepath.Join(cacheFileStoreDir, fmt.Sprintf("%x", key))
		file, err := os.OpenFile(name, os.O_RDONLY, os.ModePerm)
		if os.IsNotExist(err) {
			return false
		}
		defer file.Close()
		src, _ = ioutil.ReadAll(file)
		stat, err := file.Stat()
		if err != nil {
			log.App.Error(err)
			return false
		}
		last = gmtNowTime(stat.ModTime())
	}

	// check use cache ?
	last0 := util.Bytes2Str(ctx.Request.Header.Peek(fasthttp.HeaderIfModifiedSince))

	if last != "" && last == last0 {
		ctx.SetStatusCode(http.StatusNotModified)
		return true
	}

	setHeaderCache(p, ctx, last)
	ctx.SetBody(src)
	return true
}

// TrySetCache TrySet cache
func TrySetCache(ctx *fasthttp.RequestCtx, p IPage, body []byte) error {
	opt := p.Data().CacheOption
	if opt.Store == CacheStoreNone {
		setHeaderCache(p, ctx, "")
		return nil
	}

	key := pageCachedKey(opt, ctx, p)

	var last string

	switch opt.Store {
	case CacheStoreMem:
		last = gmtNowTime(time.Now())
		pageCacheMap.Store(string(key), pCache{LastModified: last, Src: body})

	case CacheStoreFile:
		name := filepath.Join(cacheFileStoreDir, fmt.Sprintf("%x", key))
		f, err := os.OpenFile(name, os.O_CREATE|os.O_TRUNC, os.ModePerm)
		if err != nil {
			return err
		}
		_, err = f.Write(body)
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

	setHeaderCache(p, ctx, last)
	return nil
}

// ClearCache Clear cache
func ClearCache() error {
	return nil
}

func gmtNowTime(d time.Time) string {
	lc, err := time.LoadLocation("GMT")
	if err != nil {
		panic(err)
	}
	gmt := d.In(lc)

	return gmt.Format(LastModifiedFormat)
}

func setHeaderCache(p IPage, ctx *fasthttp.RequestCtx, lastModified string) {
	ctx.SetContentTypeBytes(contentTypePage)
	// set response header
	// Cache-Control: public, max-age=3600
	// must-revalidate

	mode := p.Data().CacheOption.HTTPHead
	switch {
	case mode == 1:
		ctx.Response.Header.Set(fasthttp.HeaderETag, fmt.Sprintf("%x", util.MD5(lastModified)))
		ctx.Response.Header.Set(fasthttp.HeaderLastModified, lastModified)
	case mode > 1:
		ctx.Response.Header.Set(fasthttp.HeaderCacheControl, fmt.Sprint(headSValMaxAge, mode))
		ctx.Response.Header.Set(fasthttp.HeaderLastModified, lastModified)
	default:
		ctx.Response.Header.Set(fasthttp.HeaderCacheControl, headSValNoCache)
	}
}
