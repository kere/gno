package cache

import (
	"log"
	"sync"
	"time"

	"github.com/garyburd/redigo/redis"
	"github.com/kere/goo/libs/conf"
	libRedis "github.com/kere/goo/libs/redis"
)

// RedisCache redis
type RedisCache struct {
	Driver      string
	client      *libRedis.Pool
	EnableMutex bool
	mu          sync.Mutex
}

// NewRedisCache new
func NewRedisCache() *RedisCache {
	return &RedisCache{Driver: "redis"}
}

// GetDriver string
func (r *RedisCache) GetDriver() string {
	return r.Driver
}

// Init err
func (r *RedisCache) Init(c conf.Conf) error {
	r.EnableMutex = true

	r.client = libRedis.NewPool(c)

	return nil
}

// GetRedis return redis client
func (r *RedisCache) GetRedis() *libRedis.Pool {
	return r.client
}

// Conn connection
func (r *RedisCache) Conn() redis.Conn {
	return r.Conn()
}

// Delete remove key
func (r *RedisCache) Delete(bkey []byte) error {
	return r.client.Delete(bkey)
}

// Set func
func (r *RedisCache) Set(bkey []byte, value interface{}, expire int) error {
	var err error
	conn := r.Conn()
	defer conn.Close()
	key := string(bkey)
	r.mu.Lock()
	_, err = conn.Do("SET", key, value)
	if expire > 0 {
		conn.Do("EXPIRE", key, expire)
	}
	r.mu.Unlock()
	if err != nil {
		log.Fatalln(err)
		return err
	}
	return nil
}

// Get func
func (r *RedisCache) Get(bkey []byte) (interface{}, error) {
	conn := r.Conn()
	defer conn.Close()
	var value interface{}
	var err error
	r.mu.Lock()
	key := string(bkey)
	if r.EnableMutex {
		var exists bool
		if exists, _ = r.Exists(bkey); exists {
			value, err = conn.Do("GET", key)
		} else {
			mxPre := []byte("MUTEX-")
			if exists, _ = r.Exists(append(mxPre, bkey...)); exists {
				for i := 0; i < 10; i++ {
					time.Sleep(10 * time.Millisecond)
					if exists, _ = r.Exists(bkey); exists {
						value, err = conn.Do("GET", key)
						break
					} else {
						continue
					}
				}
			} else {
				r.Set(append(mxPre, bkey...), true, 1)
				value = nil
			}
		}
	} else {
		value, err = conn.Do("GET", key)
	}
	r.mu.Unlock()

	if err != nil {
		log.Fatalln(err)
		return nil, err
	}
	return value, nil
}

// Exists key
func (r *RedisCache) Exists(bkey []byte) (bool, error) {
	return r.client.Exists(bkey)
}

// GetString string
func (r *RedisCache) GetString(bkey []byte) (string, error) {
	return redis.String(r.Get(bkey))
}

// GetInt int
func (r *RedisCache) GetInt(bkey []byte) (int, error) {
	return redis.Int(r.Get(bkey))
}

// GetInt64 int64
func (r *RedisCache) GetInt64(bkey []byte) (int64, error) {
	return redis.Int64(r.Get(bkey))
}

// GetFloat float64
func (r *RedisCache) GetFloat(bkey []byte) (float64, error) {
	return redis.Float64(r.Get(bkey))
}
