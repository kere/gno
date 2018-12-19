package gno

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"

	"github.com/kere/gno/libs/cache"
	"github.com/kere/gno/libs/util"
)

const (
	//CacheModeNone 不缓存页面
	CacheModeNone = 0
	//CacheModePage 缓存页面
	CacheModePage = 1
	//CacheModePagePath 缓存页面
	CacheModePagePath = 2
	//CacheModeFile 文件缓存页面
	CacheModeFile = 3

	pagecacheKeyPrefix = "c:"
	pageCacheSubfix    = ".htm"
)

const (
	// HeaderEtag etag
	HeaderEtag = "ETag"
	// HeaderCacheCtl cache
	HeaderCacheCtl = "Cache-Control"
	// HeaderIfNoneMatch If-None-Match
	HeaderIfNoneMatch = "If-None-Match"
)

// SetPageCache value
func (p *Page) SetPageCache(mode, expires int) {
	p.CacheMode = mode
	p.Expires = expires
}

// GetExpires value
func (p *Page) GetExpires() int {
	return p.Expires
}

// GetCacheMode value
func (p *Page) GetCacheMode() int {
	return p.CacheMode
}

// TryCache try to get cache
func TryCache(p IPage) bool {
	var src []byte
	var err error
	switch p.GetCacheMode() {
	case CacheModePage:
		key := pagecacheKeyPrefix + p.GetDir() + p.GetName()
		src, err = cache.GetBytes(key)

	case CacheModePagePath:
		key := pagecacheKeyPrefix + p.GetRequest().URL.Path
		src, err = cache.GetBytes(key)

	case CacheModeFile:
		filename := filepath.Join(WEBROOT, p.GetRequest().URL.Path+pageCacheSubfix)
		src, err = ioutil.ReadFile(filename)

	default:
		return false
	}

	if err != nil {
		return false
	}

	req := p.GetRequest()
	w := p.GetResponseWriter()
	// check use cache ?
	etag := req.Header.Get(HeaderIfNoneMatch)

	// for k, v := range req.Header {
	// 	fmt.Println(k, "=", v)
	// }
	// fmt.Println(etag)

	token := fmt.Sprintf("%x", util.MD5(src))
	// token := string(util.CRC64Token(src))
	if etag == token {
		w.WriteHeader(http.StatusNotModified)
		return true
	}

	h := w.Header()
	// Cache-Control: public, max-age=3600
	h.Add(HeaderCacheCtl, "public")
	h.Set(HeaderEtag, token)
	// fmt.Println("response::")
	// for k, v := range h {
	// 	fmt.Println(k, "=", v)
	// }
	// fmt.Println()

	w.Write(src)
	return true
}

// TrySetCache TrySet cache
func TrySetCache(p IPage, buf *bytes.Buffer) error {
	switch p.GetCacheMode() {
	case CacheModePage:
		key := pagecacheKeyPrefix + p.GetDir() + p.GetName()
		return cache.Set(key, buf.String(), p.GetExpires())

	case CacheModePagePath:
		key := pagecacheKeyPrefix + p.GetRequest().URL.Path
		return cache.Set(key, buf.String(), p.GetExpires())

	case CacheModeFile:
		filename := filepath.Join(WEBROOT, p.GetRequest().URL.Path+pageCacheSubfix)
		err := os.MkdirAll(filepath.Dir(filename), os.ModeDir)
		if err != nil {
			return err
		}

		f, err := os.OpenFile(filename, os.O_CREATE|os.O_TRUNC, os.ModePerm)
		if err != nil {
			return err
		}
		defer f.Close()
		_, err = f.Write(buf.Bytes())
		return err

	default:
		return nil
	}
}

// TryClearCache Clear cache
func TryClearCache() error {
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
