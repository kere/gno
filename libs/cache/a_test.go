package cache

import (
	"fmt"
	"testing"
	"time"
)

type cacheMap struct {
	Map
}

var counter int

func (c *cacheMap) Build(args ...interface{}) (interface{}, error) {
	str := args[0].(string)
	counter++
	fmt.Println("build cache ", str, counter)
	return counter, nil
}

func Test_map(t *testing.T) {
	c := &cacheMap{}
	c.Init(c, 3)

	v := c.Get("a")
	if v.(int) != 1 {
		fmt.Println(v)
	}

	time.Sleep(2 * time.Second)
	v = c.Get("a")
	if v.(int) != 1 {
		fmt.Println(v)
	}

	time.Sleep(1 * time.Second)
	v = c.Get("a")
	if v.(int) != 2 {
		fmt.Println(v)
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
