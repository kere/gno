package cache

import (
	"fmt"
	"sync"
	"testing"
)

type cacheMapSync struct {
	Map
}

var counter int
var counterSync int

func (c *cacheMapSync) Build(args ...interface{}) (interface{}, error) {
	str := args[0].(string)
	counterSync++
	fmt.Println("build cache ", str)
	return str, nil
}

func Test_MapSync(t *testing.T) {
	c := &cacheMapSync{}
	c.Init(c, 0)
	var wg sync.WaitGroup
	f := func() {
		c.Get("a")
		wg.Done()
	}
	for i := 0; i < 20; i++ {
		go func() {
			wg.Add(1)
			f()
		}()
	}
	wg.Wait()
	if counterSync != 1 {
		t.Fatal()
	}
}

func Test_func(t *testing.T) {
	// driver=redis
	// network=tcp
	// connect=127.0.0.1:6379
	// db=3

	conf := make(map[string]string, 0)
	conf["driver"] = "redis"
	conf["network"] = "tcp"
	conf["db"] = "5"
	conf["connect"] = "127.0.0.1:6379"
	// conf["max_idle"] = "10"
	// conf["max_active"] = "0"
	// conf["idle_timeout"] = "0"

	Init(conf)

	err := Set("001", "abc1", 100)
	if err != nil {
		t.Fatal(err)
	}
	v, err := Get("001")
	if err != nil {
		t.Fatal(err)
	}
	if v != "abc1" {
		t.Fatal("v=", v)
	}

	rds := GetRedis()

	v, err = rds.DoString("get", "001")
	if err != nil {
		t.Fatal(err)
	}
	if v != "abc1" {
		t.Fatal("v=", v)
	}

	fmt.Println(v, err)
}
