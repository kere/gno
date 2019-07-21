package redis

import (
	"log"
	"strconv"
	"time"

	"github.com/gomodule/redigo/redis"
	"github.com/kere/gno/libs/conf"
)

// NewPool new func
func NewPool(config map[string]string) *Pool {
	r := &Pool{}
	c := conf.Conf(config)

	r.MaxIdle = c.DefaultInt("max_idle", 2)
	r.MaxActive = c.DefaultInt("max_active", 10)
	r.IdleTimeout = time.Duration(c.DefaultInt("idle_timeout", 180)) * time.Second
	r.Network = c.DefaultString("network", "tcp")
	r.connectString = c.DefaultString("connect", "127.0.0.1:6379")
	r.dbindex = c.DefaultInt("db", 0)

	r.pool = &redis.Pool{
		MaxIdle:     r.MaxIdle,
		MaxActive:   r.MaxActive,
		IdleTimeout: r.IdleTimeout,
		Dial: func() (redis.Conn, error) {
			c, err := redis.Dial(r.Network, r.connectString)
			if err != nil {
				return nil, err
			}
			c.Do("SELECT", r.dbindex)
			return c, nil
		},
	}
	return r
}

// Pool struct
type Pool struct {
	// pool                          *redis.Pool
	pool          *redis.Pool
	dbindex       int
	MaxIdle       int
	MaxActive     int
	IdleTimeout   time.Duration
	connectString string
	Network       string
}

//SetDB set db index
func (r *Pool) SetDB(i int) {
	r.dbindex = i
	r.Send("SELECT", r.dbindex)
}

// Pool return redis.Pool
// func (r *Pool) Pool() *redis.Pool {
func (r *Pool) Pool() *redis.Pool {
	return r.pool
}

// Close func
func (r *Pool) Close() {
	r.pool.Close()
}

// Connect to redis server
func (r *Pool) Connect() (redis.Conn, error) {
	c, err := redis.Dial(r.Network, r.connectString)
	if err != nil {
		log.Fatalln("connect error: ", err)
		return nil, err
	}

	if r.dbindex > 0 {
		c.Send("SELECT", r.dbindex)
	}
	return c, nil
}

// Conn return redis.Conn
func (r *Pool) Conn() redis.Conn {
	return r.pool.Get()
}

// IsOK is client ok
func (r *Pool) IsOK() bool {
	c := r.Conn()
	err := c.Send("ping")
	if err != nil {
		c.Close()
		return false
	}
	c.Close()
	return true
}

// Exists key?
func (r *Pool) Exists(key string) (bool, error) {
	return r.DoBool("exists", key)
}

// Delete key
func (r *Pool) Delete(key string) error {
	return r.Send("del", key)
}

// Get key
func (r *Pool) Get(key string) (interface{}, error) {
	return r.Do("get", key)
}

// Send return error
func (r *Pool) Send(m string, args ...interface{}) error {
	c := r.Conn()
	// defer c.Close()
	_, err := c.Do(m, args...)
	c.Close()
	return err
}

// Do somthing
func (r *Pool) Do(m string, args ...interface{}) (interface{}, error) {
	c := r.Conn()
	v, err := c.Do(m, args...)
	c.Close()
	return v, err
}

// DoBool return bool
func (r *Pool) DoBool(m string, args ...interface{}) (bool, error) {
	return redis.Bool(r.Do(m, args...))
}

// DoString return string
func (r *Pool) DoString(m string, args ...interface{}) (string, error) {
	return redis.String(r.Do(m, args...))
}

// DoInt64 return int64
func (r *Pool) DoInt64(m string, args ...interface{}) (int64, error) {
	return redis.Int64(r.Do(m, args...))
}

// DoUint64 return uint64
func (r *Pool) DoUint64(m string, args ...interface{}) (uint64, error) {
	return redis.Uint64(r.Do(m, args...))
}

// DoInt return int
func (r *Pool) DoInt(m string, args ...interface{}) (int, error) {
	return redis.Int(r.Do(m, args...))
}

// DoFloat return float64
func (r *Pool) DoFloat(m string, args ...interface{}) (float64, error) {
	return redis.Float64(r.Do(m, args...))
}

// DoBytes return []byte
func (r *Pool) DoBytes(m string, args ...interface{}) ([]byte, error) {
	return redis.Bytes(r.Do(m, args...))
}

// DoByteSlice return [][]byte
func (r *Pool) DoByteSlice(m string, args ...interface{}) ([][]byte, error) {
	return redis.ByteSlices(r.Do(m, args...))
}

// DoValues return []interface{}
func (r *Pool) DoValues(m string, args ...interface{}) ([]interface{}, error) {
	return redis.Values(r.Do(m, args...))
}

// DoInts return []int
func (r *Pool) DoInts(m string, args ...interface{}) ([]int, error) {
	return redis.Ints(r.Do(m, args...))
}

// DoInt64s []int64
func (r *Pool) DoInt64s(m string, args ...interface{}) ([]int64, error) {
	var ints []int64
	reply, err := r.Do(m, args...)
	values, err := redis.Values(reply, err)
	if err != nil {
		return ints, err
	}

	if err := redis.ScanSlice(values, &ints); err != nil {
		return ints, err
	}

	return ints, nil
}

// DoFloats []float64
func (r *Pool) DoFloats(m string, args ...interface{}) ([]float64, error) {
	s, err := redis.Strings(r.Do(m, args...))
	if err != nil {
		return []float64{}, err
	}
	tmp := make([]float64, len(s))
	for i := range s {
		tmp[i], _ = strconv.ParseFloat(s[i], 64)
	}

	return tmp, nil
}

// DoStrings []string
func (r *Pool) DoStrings(m string, args ...interface{}) ([]string, error) {
	return redis.Strings(r.Do(m, args...))
}

// DoInt64Map map[string]int64
func (r *Pool) DoInt64Map(m string, args ...interface{}) (map[string]int64, error) {
	return redis.Int64Map(r.Do(m, args...))
}

// DoStringMap map[string]string
func (r *Pool) DoStringMap(m string, args ...interface{}) (map[string]string, error) {
	return redis.StringMap(r.Do(m, args...))
}

// DoIntMap map[string]int
func (r *Pool) DoIntMap(m string, args ...interface{}) (map[string]int, error) {
	return redis.IntMap(r.Do(m, args...))
}
