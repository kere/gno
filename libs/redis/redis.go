package redis

import (
	"context"
	"log"
	"strconv"
	"time"

	"github.com/garyburd/redigo/redis"
	"github.com/kere/goo/libs/conf"
	"github.com/youtube/vitess/go/pools"
)

// ResourceConn adapts a Redigo connection to a Vitess Resource.
type ResourceConn struct {
	redis.Conn
}

// Close connection
func (r ResourceConn) Close() {
	r.Conn.Close()
}

// NewPool new func
func NewPool(config map[string]string) *Pool {
	r := &Pool{}
	c := conf.Conf(config)
	r.MaxCap = c.DefaultInt("max_cap", 200)
	r.Capacity = c.DefaultInt("capacity", 100)
	r.IdleTimeout = time.Duration(c.DefaultInt("idle_timeout", 300)) * time.Second
	r.connectNetwork = c.DefaultString("network", "tcp")
	r.connectString = c.DefaultString("connect", "127.0.0.1:6379")
	r.dbindex = c.DefaultInt("db", 0)

	r.pool = r.newPool()
	return r
}

// Pool struct
type Pool struct {
	// pool                          *redis.Pool
	pool                          *pools.ResourcePool
	dbindex                       int
	MaxCap                        int
	Capacity                      int
	IdleTimeout                   time.Duration
	connectString, connectNetwork string
}

func (r *Pool) newPool() *pools.ResourcePool {
	// return &redis.Pool{
	// 	MaxIdle:     r.MaxIdle,
	// 	MaxActive:   r.MaxActive,
	// 	IdleTimeout: r.IdleTimeout,
	// 	Dial:        r.Connect,
	// 	Wait:        false,
	// 	TestOnBorrow: func(c redis.Conn, t time.Time) error {
	// 		if time.Since(t) < time.Minute {
	// 			return nil
	// 		}
	// 		_, err := c.Do("PING")
	// 		return err
	// 	},
	// }
	// factory Factory, capacity, maxCap int, idleTimeout time.Duration
	return pools.NewResourcePool(func() (pools.Resource, error) {
		c, err := r.Connect()
		return &ResourceConn{c}, err
	}, r.Capacity, r.MaxCap, r.IdleTimeout)
}

//SetDB set db index
func (r *Pool) SetDB(i int) {
	r.dbindex = i
	r.Send("SELECT", r.dbindex)
}

// Pool return redis.Pool
// func (r *Pool) Pool() *redis.Pool {
func (r *Pool) Pool() *pools.ResourcePool {
	return r.pool
}

// Close func
func (r *Pool) Close() {
	r.pool.Close()
}

// Connect to redis server
func (r *Pool) Connect() (redis.Conn, error) {
	c, err := redis.Dial(r.connectNetwork, r.connectString)
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
func (r *Pool) Conn() *ResourceConn {
	// return r.pool.Get()
	ctx := context.TODO()
	c, err := r.pool.Get(ctx)
	if err != nil {
		panic(err)
	}
	cc := c.(*ResourceConn)
	return cc
}

// IsOK is client ok
func (r *Pool) IsOK() bool {
	c := r.Conn()
	defer c.Close()
	err := c.Send("ping")
	if err != nil {
		return false
	}
	return true
}

// Exists key?
func (r *Pool) Exists(bkey []byte) (bool, error) {
	// c := r.Conn()
	// defer c.Close()
	// return redis.Bool(c.Do("EXISTS", string(bkey)))

	return r.DoBool("exists", string(bkey))
}

// Delete key
func (r *Pool) Delete(bkey []byte) error {
	// c := r.Conn()
	// defer c.Close()
	// _, err := c.Do("DEL", string(bkey))
	return r.Send("del", string(bkey))
}

// RecoveryConn func
func (r *Pool) RecoveryConn(c *ResourceConn) {
	r.pool.Put(c)
}

// Get key
func (r *Pool) Get(bkey []byte) (interface{}, error) {
	// c := r.Conn()
	// defer c.Close()
	// return c.Do("GET", string(bkey))
	return r.Do("get", string(bkey))
}

// Send return error
func (r *Pool) Send(m string, args ...interface{}) error {
	c := r.Conn()
	// defer c.Close()
	defer r.RecoveryConn(c)
	_, err := c.Do(m, args...)
	return err
}

// Do somthing
func (r *Pool) Do(m string, args ...interface{}) (interface{}, error) {
	c := r.Conn()
	// defer c.Close()
	defer r.RecoveryConn(c)
	return c.Do(m, args...)
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
