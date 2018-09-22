package cache

import (
	"fmt"

	"github.com/kere/gno/libs/conf"
	"github.com/kere/gno/libs/redis"
)

// ICache interface
type ICache interface {
	Init(conf.Conf) error
	Set(string, string, int) error
	Get(string) (interface{}, error)

	GetDriver() string
	GetString(string) (string, error)
	GetInt(string) (int, error)
	GetInt64(string) (int64, error)
	GetFloat(string) (float64, error)
	IsExists(string) (bool, error)
	Delete(string) error
}

var (
	cache ICache
)

// Init func
func Init(config conf.Conf) error {
	switch config.GetString("driver") {
	case "redis":
		cache = NewRedisCache()
	default:
		return fmt.Errorf("no cache driver found:%s", config["driver"])
	}

	return cache.Init(config)
}

// IsOK  is cache ok.
func IsOK() bool {
	return cache != nil
}

// CurrentCache icache
func CurrentCache() ICache {
	return cache
}

// IsEnable bool
func IsEnable() bool {
	return cache != nil
}

// IsExists key
func IsExists(key string) (bool, error) {
	return cache.IsExists(key)
}

// Get key
func Get(key string) (interface{}, error) {
	return cache.Get(key)
}

// GetString key
func GetString(key string) (string, error) {
	return cache.GetString(key)
}

// GetInt int
func GetInt(key string) (int, error) {
	return cache.GetInt(key)
}

// GetInt64 int64
func GetInt64(key string) (int64, error) {
	return cache.GetInt64(key)
}

// GetFloat float64
func GetFloat(key string) (float64, error) {
	return cache.GetFloat(key)
}

// Set func
func Set(key string, value string, expire int) error {
	return cache.Set(key, value, expire)
}

// Delete key
func Delete(key string) error {
	return cache.Delete(key)
}

// GetRedis return client
func GetRedis() *redis.Pool {
	if CurrentCache() == nil {
		panic("redis is not initalized")
	}
	if CurrentCache().GetDriver() != "redis" {
		return nil
	}
	c := CurrentCache().(*RedisCache)
	return c.GetRedis()
}
