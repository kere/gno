package httpd

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/kere/gno/libs/log"
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

	//CacheStoreMem to store in memory
	CacheStoreMem = 0
	//CacheStoreFile to store in file
	CacheStoreFile = 1

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

	//LastModifiedFormat Wed, 21 Oct 2015 07:28:00 GMT
	// Last-Modified: <day-name>, <day> <month> <year> <hour>:<minute>:<second> GMT
	LastModifiedFormat = "Mon, 02 Jan 2006 15:04:05 GMT"

	cacheFileStoreDir = "var/cache/page"
)

// pageCachedKey key
func pageCachedKey(opt PageCacheOption, ctx *fasthttp.RequestCtx, p IPage) string {
	switch opt.Mode {
	case CacheModePage:
		return pagecacheKeyPrefix + p.Dir() + p.Name()

	case CacheModePagePath:
		return pagecacheKeyPrefix + string(ctx.URI().Path())

	case CacheModePageURI:
		return pagecacheKeyPrefix + string(ctx.URI().RequestURI())

	default:
		return ""
	}
}

// TryGetCache try to get cache
func TryGetCache(ctx *fasthttp.RequestCtx, p IPage) bool {
	if RunMode == ModeDev {
		return false
	}

	opt := p.CacheOption()
	key := pageCachedKey(opt, ctx, p)

	var src []byte
	var last string

	switch opt.StoreMode {
	case CacheStoreMem:
		if v, isCached := pageCached.Load(key); isCached {
			pageElem := v.(PageCacheElem)
			last = pageElem.LastModified
		}

	case CacheStoreFile:
		bstr, err := base64.StdEncoding.DecodeString(key)
		if err != nil {
			log.App.Alert(err)
			return false
		}
		name := filepath.Join(cacheFileStoreDir, string(bstr))
		file, err := os.OpenFile(name, os.O_RDONLY, os.ModePerm)
		if os.IsExist(err) {
			file.Close()
			src, _ = ioutil.ReadAll(file)
			stat, err := file.Stat()
			if err != nil {
				last = stat.ModTime().Format(LastModifiedFormat)
			}
		}

	}

	// check use cache ?
	last0 := string(ctx.Request.Header.Peek(HeaderIfModifiedSince))

	if last != "" && last == last0 {
		// w.WriteHeader(http.StatusNotModified)
		ctx.SetStatusCode(http.StatusNotModified)
		return true
	}

	setResponseHeader(opt, ctx, p, last)

	// w.Write(src)
	ctx.Write(src)

	return true
}

// TrySetCache TrySet cache
func TrySetCache(ctx *fasthttp.RequestCtx, p IPage, buf *bytes.Buffer) error {
	opt := p.CacheOption()
	key := pageCachedKey(opt, ctx, p)
	last := time.Now().Format(LastModifiedFormat)

	switch opt.StoreMode {
	case CacheStoreMem:

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
			return err
		}
		f.Close()

	}

	setResponseHeader(opt, ctx, p, last)

	return nil
}

func setResponseHeader(opt PageCacheOption, ctx *fasthttp.RequestCtx, p IPage, lastModified string) {
	// set response header
	// Cache-Control: public, max-age=3600
	if opt.HeadExpires > 0 {
		ctx.Request.Header.Add(HeaderCacheCtl, fmt.Sprint("must-revalidate, max-age=", opt.HeadExpires))
	} else {
		ctx.Request.Header.Add(HeaderCacheCtl, "must-revalidate")
	}
	ctx.Request.Header.Set(HeaderLastModified, lastModified)
}

// ClearCache Clear cache
func ClearCache() error {
	return nil
}

// PageCacheOption option
type PageCacheOption struct {
	Mode         int // 0:
	StoreMode    int // 0: mem 1:file
	HeadExpires  int // http head expires
	CacheExpires int
}

// PageCacheElem class
type PageCacheElem struct {
	LastModified string // Last-Modified: Fri, 12 May 2006 18:53:33 GMT
	Src          []byte
}

var pageCached = &sync.Map{}
