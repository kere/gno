package cache

import (
	"fmt"

	"github.com/kere/goo/libs/conf"
	"github.com/kere/goo/libs/redis"
)

// ICache interface
type ICache interface {
	Init(conf.Conf) error
	GetDriver() string
	Get([]byte) (interface{}, error)
	GetString([]byte) (string, error)
	GetInt([]byte) (int, error)
	GetInt64([]byte) (int64, error)
	GetFloat([]byte) (float64, error)
	Set([]byte, interface{}, int) error
	Exists([]byte) (bool, error)
	Delete([]byte) error
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

// Exists key
func Exists(bkey []byte) (bool, error) {
	return cache.Exists(bkey)
}

// Get key
func Get(bkey []byte) (interface{}, error) {
	return cache.Get(bkey)
}

// GetString key
func GetString(bkey []byte) (string, error) {
	return cache.GetString(bkey)
}

// GetInt int
func GetInt(bkey []byte) (int, error) {
	return cache.GetInt(bkey)
}

// GetInt64 int64
func GetInt64(bkey []byte) (int64, error) {
	return cache.GetInt64(bkey)
}

// GetFloat float64
func GetFloat(bkey []byte) (float64, error) {
	return cache.GetFloat(bkey)
}

// Set func
func Set(bkey []byte, value interface{}, expire int) error {
	return cache.Set(bkey, value, expire)
}

// Delete key
func Delete(bkey []byte) error {
	return cache.Delete(bkey)
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
