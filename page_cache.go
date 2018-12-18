package gno

import (
	"bytes"
	"io"

	"github.com/kere/gno/libs/cache"
)

const (
	//CacheModeNone 不缓存页面
	CacheModeNone = 0
	//CacheModePage 缓存页面
	CacheModePage = 1
	//CacheModePagePath 缓存页面
	CacheModePagePath = 2

	pagecacheKeyPrefix = "c:"
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
func TryCache(p IPage, w io.Writer) bool {
	var key string
	switch p.GetCacheMode() {
	case CacheModePage:
		key = pagecacheKeyPrefix + p.GetDir() + p.GetName()
	case CacheModePagePath:
		key = pagecacheKeyPrefix + p.GetRequest().URL.Path
	default:
		return false
	}

	src, err := cache.GetString(key)
	if err != nil {
		return false
	}

	w.Write([]byte(src))
	return true
}

// TrySetCache TrySet cache
func TrySetCache(p IPage, buf *bytes.Buffer) error {
	var key string
	switch p.GetCacheMode() {
	case CacheModePage:
		key = pagecacheKeyPrefix + p.GetDir() + p.GetName()
	case CacheModePagePath:
		key = pagecacheKeyPrefix + p.GetRequest().URL.Path
	default:
		return nil
	}
	return cache.Set(key, buf.String(), p.GetExpires())
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
