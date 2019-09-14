package httpd

import (
	"crypto/md5"
	"fmt"
	"hash/crc32"
	"io/ioutil"
	"net/http"
	"net/url"
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
func pageCachedKey(ctx *fasthttp.RequestCtx, attr *PageAttr) string {
	var src []byte
	switch attr.CacheOption.PageMode {
	case CacheModePage:
		src = util.Str2Bytes(attr.Dir + attr.Name)

	case CacheModePageURI:
		src = ctx.URI().RequestURI()

	default:
		src = ctx.URI().Path()
	}
	return util.Bytes2Str(cachedKey(src))
}

// cachedKey key
func cachedKey(src []byte) []byte {
	ieee := crc32.NewIEEE()
	// io.WriteString(ieee, str)
	ieee.Write(src)
	v64 := uint64(ieee.Sum32())
	return util.IntZipTo62(v64)
}

// DisablePageCache bool
var DisablePageCache = true

// ClearCache clear page cache
func ClearCache(urlstr []byte, p IPage) {
	if DisablePageCache {
		return
	}

	u, err := url.Parse(util.Bytes2Str(urlstr))
	if err != nil {
		panic(err)
	}

	attr := p.Attr()

	var src []byte
	switch attr.CacheOption.PageMode {
	case CacheModePage:
		src = util.Str2Bytes(attr.Dir + attr.Name)
	case CacheModePageURI:
		src = util.Str2Bytes(u.RequestURI())
	default:
		src = util.Str2Bytes(u.Path)
	}

	key := cachedKey(src)

	switch attr.CacheOption.Store {
	case CacheStoreMem:
		pageCacheMap.Delete(key)

	case CacheStoreFile:
		name := filepath.Join(HomeDir, cacheFileStoreDir, fmt.Sprintf("%x", key))
		_, err := os.Stat(name)
		if os.IsNotExist(err) {
			return
		}
		os.Remove(name)
	}
}

// TryCache try to get cache
func TryCache(ctx *fasthttp.RequestCtx, p IPage) bool {
	if DisablePageCache {
		return false
	}

	attr := p.Attr()

	var src []byte
	var last string

	switch attr.CacheOption.Store {
	case CacheStoreNone:
		return false

	case CacheStoreMem:
		key := pageCachedKey(ctx, attr)
		v, isFound := pageCacheMap.Load(key)
		if !isFound {
			return false
		}
		pe := v.(pCache)
		src = pe.Src
		last = pe.LastModified
		fmt.Println("mem cache found", key)

	case CacheStoreFile:
		key := pageCachedKey(ctx, attr)
		name := filepath.Join(HomeDir, cacheFileStoreDir, fmt.Sprintf("%x", key))
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

		fmt.Println("file cache found", key)

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
	if DisablePageCache {
		return nil
	}
	attr := p.Attr()

	var last string
	switch attr.CacheOption.Store {
	case CacheStoreNone:
		setHeaderCache(p, ctx, "")
		return nil

	case CacheStoreMem:
		key := pageCachedKey(ctx, attr)
		last = gmtNowTime(time.Now())
		pageCacheMap.Store(key, pCache{LastModified: last, Src: body})
		fmt.Println("set mem cache")

	case CacheStoreFile:
		key := pageCachedKey(ctx, attr)
		name := filepath.Join(cacheFileStoreDir, fmt.Sprintf("%x", key))
		f, err := os.OpenFile(name, os.O_CREATE|os.O_TRUNC, os.ModePerm)
		if err != nil {
			return err
		}
		defer f.Close()
		_, err = f.Write(body)
		if err != nil {
			return err
		}
		info, err := f.Stat()
		if err != nil {
			return err
		}
		last = gmtNowTime(info.ModTime())
		fmt.Println("set file cache")
	}

	setHeaderCache(p, ctx, last)
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
	// set response header
	// Cache-Control: public, max-age=3600
	// must-revalidate

	mode := p.Attr().CacheOption.HTTPHead
	switch {
	case mode == 1:
		ctx.Response.Header.Set(fasthttp.HeaderETag, fmt.Sprintf("%x", md5.Sum(util.Str2Bytes(lastModified))))
		ctx.Response.Header.Set(fasthttp.HeaderLastModified, lastModified)
	case mode > 1:
		ctx.Response.Header.Set(fasthttp.HeaderCacheControl, fmt.Sprint(headSValMaxAge, mode))
		ctx.Response.Header.Set(fasthttp.HeaderLastModified, lastModified)
	default:
		ctx.Response.Header.Set(fasthttp.HeaderCacheControl, headSValNoCache)
	}
}
