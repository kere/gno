package gno

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/kere/gno/libs/cache"
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

	delim1 = byte('\n')
	delim2 = "\n"
)

const (
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
	var src, srcTmp []byte
	var err error
	var last string
	switch p.GetCacheMode() {
	case CacheModePage:
		key := pagecacheKeyPrefix + p.GetDir() + p.GetName()
		srcTmp, err = cache.GetBytes(key)
		buf := bytes.NewBuffer(srcTmp)
		// 读取第一行时间戳
		last, _ = buf.ReadString(delim1)
		last = strings.TrimRight(last, delim2)
		src, _ = ioutil.ReadAll(buf)

	case CacheModePagePath:
		key := pagecacheKeyPrefix + p.GetRequest().URL.Path
		srcTmp, err = cache.GetBytes(key)
		buf := bytes.NewBuffer(srcTmp)
		// 读取第一行时间戳
		last, _ = buf.ReadString(delim1)
		last = strings.TrimRight(last, delim2)
		src, _ = ioutil.ReadAll(buf)

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
	// etag0 := req.Header.Get(HeaderIfNoneMatch)
	last0 := req.Header.Get(HeaderIfModifiedSince)

	h := w.Header()
	// etag := fmt.Sprintf("%x", util.MD5(src))
	// etag := string(util.CRC64Token(src))
	// fmt.Println("etag:", etag, " etag0:", etag0)
	// fmt.Println("last:", last, " last0:", last0)

	if last != "" && last == last0 {
		w.WriteHeader(http.StatusNotModified)
		// h.Set(HeaderIfNoneMatch, "false")
		return true
	}

	// Cache-Control: public, max-age=3600
	h.Add(HeaderCacheCtl, "must-revalidate")
	// h.Set(HeaderEtag, last)
	// h.Set(HeaderIfNoneMatch, "true")
	h.Set(HeaderLastModified, last)

	w.Write(src)
	return true
}

// TrySetCache TrySet cache
func TrySetCache(p IPage, buf *bytes.Buffer) error {
	switch p.GetCacheMode() {
	case CacheModePage:
		key := pagecacheKeyPrefix + p.GetDir() + p.GetName()
		last := time.Now().Format(time.RFC1123)
		src := append([]byte(last), byte('\n'))
		src = append(src, buf.Bytes()...)
		p.GetResponseWriter().Header().Set(HeaderLastModified, last)

		return cache.Set(key, string(src), p.GetExpires())

	case CacheModePagePath:
		key := pagecacheKeyPrefix + p.GetRequest().URL.Path
		last := time.Now().Format(time.RFC1123)
		src := append([]byte(last), byte('\n'))
		src = append(src, buf.Bytes()...)
		p.GetResponseWriter().Header().Set(HeaderLastModified, last)
		return cache.Set(key, string(src), p.GetExpires())

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
