package gno

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"time"

	"github.com/kere/gno/libs/cache"
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

// SetCacheOption value
func (p *Page) SetCacheOption(mode, headerExp, cacheExp int) {
	p.CacheOption = CacheOption{Mode: mode, HeaderExpires: headerExp, CacheExpires: cacheExp}
}

// GetCacheOption value
func (p *Page) GetCacheOption() *CacheOption {
	return &p.CacheOption
}

// TryCache try to get cache
func TryCache(p IPage) bool {
	if RunMode == ModeDev {
		return false
	}

	var src, srcTmp []byte
	var err error
	var last, key string
	switch p.GetCacheOption().Mode {
	case CacheModePage:
		key = pagecacheKeyPrefix + p.GetDir() + p.GetName()

	case CacheModePagePath:
		key = pagecacheKeyPrefix + p.GetRequest().URL.Path

	case CacheModePageURI:
		key = pagecacheKeyPrefix + p.GetRequest().URL.RequestURI()

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

	req := p.GetRequest()
	w := p.GetResponseWriter()
	// check use cache ?
	// etag0 := req.Header.Get(HeaderIfNoneMatch)
	last0 := req.Header.Get(HeaderIfModifiedSince)

	// h := w.Header()
	// etag := fmt.Sprintf("%x", util.MD5(src))
	// etag := string(util.CRC64Token(src))
	// fmt.Println("etag:", etag, " etag0:", etag0)
	// fmt.Println("last:", last, " last0:", last0)

	if last != "" && last == last0 {
		w.WriteHeader(http.StatusNotModified)
		// h.Set(HeaderIfNoneMatch, "false")
		return true
	}

	// h.Set(HeaderEtag, last)
	// h.Set(HeaderIfNoneMatch, "true")
	setResponseHeader(p, last)

	w.Write(src)
	return true
}

// TrySetCache TrySet cache
func TrySetCache(p IPage, buf *bytes.Buffer) error {
	var err error
	var last string
	switch p.GetCacheOption().Mode {
	case CacheModePage:
		key := pagecacheKeyPrefix + p.GetDir() + p.GetName()
		last = time.Now().Format(time.RFC1123)
		src := append([]byte(last), byte('\n'))
		src = append(src, buf.Bytes()...)
		err = cache.Set(key, string(src), p.GetCacheOption().CacheExpires)

	case CacheModePagePath:
		key := pagecacheKeyPrefix + p.GetRequest().URL.Path
		last = time.Now().Format(time.RFC1123)
		src := append([]byte(last), byte('\n'))
		src = append(src, buf.Bytes()...)
		err = cache.Set(key, string(src), p.GetCacheOption().CacheExpires)

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

	setResponseHeader(p, last)

	return nil
}
func setResponseHeader(p IPage, lastModified string) {
	h := p.GetResponseWriter().Header()
	opt := p.GetCacheOption()
	// set response header
	// Cache-Control: public, max-age=3600
	if opt.HeaderExpires > 0 {
		h.Add(HeaderCacheCtl, "must-revalidate, max-age="+fmt.Sprint(opt.HeaderExpires))
	} else {
		h.Add(HeaderCacheCtl, "must-revalidate")
	}
	h.Set(HeaderLastModified, lastModified)
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
