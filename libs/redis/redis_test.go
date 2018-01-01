package redis

import "testing"

var rs *Pool

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

	// redis.Do("flushdb")
	rs = NewPool(conf)
	rs.Send("flushdb")

	key := "data:test:zset0"
	err := rs.Send("zadd", key, 10, 11)
	if err != nil {
		t.Fatal(err)
	}
	_, err = rs.Do("zadd", key, 12, 13)
	if err != nil {
		t.Fatal(err)
	}
}

func Test_Func(t *testing.T) {
	key := "data:test:zset"
	l := 100
	for i := 0; i < l; i++ {
		_, err := rs.Do("zadd", key, i, i)
		if err != nil {
			t.Fatal(err)
		}
	}

	ints, err := rs.DoInt64s("zrange", key, 0, -1)
	if err != nil {
		t.Fatal(err)
	}

	if len(ints) != 100 {
		t.Fatal("ints lens != 100", ints)
	}

	if ints[99] != 99 {
		t.Fatal("last value != 99", ints[99])
	}
}

func Test_M1(t *testing.T) {
	l := 100
	key := "data:test:exists1"
	_, err := rs.Do("hset", key, "field1", 1)
	if err != nil {
		t.Fatal(err)
	}

	for i := 0; i < l; i++ {
		v, err := rs.DoBool("exists", key)
		if err != nil {
			t.Fatal(err)
		}
		if !v {
			t.Fatal("Key is not exists")
			return
		}
	}
}

func Benchmark_M2(b *testing.B) {
	errChan := make(chan error)
	boolChan := make(chan bool)

	l := 300
	key := "data:test:exists2"
	v, err := rs.DoBool("exists", key)
	if err != nil {
		b.Fatal(err)
	}
	if !v {
		rs.Do("hset", key, "field1", 1)
	}

	for i := 0; i < l; i++ {
		go func() {
			key := "data:test:exists2"
			v, err := rs.DoBool("exists", key)
			errChan <- err
			boolChan <- v
		}()
	}

	for i := 0; i < l; i++ {
		err := <-errChan
		v := <-boolChan
		if err != nil {
			b.Fatal(err)
		}

		if !v {
			b.Fatal("key not found!")
		}
	}
}

func Benchmark_Buffer(b *testing.B) {
	errChan := make(chan error, 50)

	l := 3000
	key := "data:test:exists2"
	v, err := rs.DoBool("exists", key)
	if err != nil {
		b.Fatal(err)
	}
	if !v {
		rs.Do("hset", key, "field1", 1)
	}

	for i := 0; i < l; i++ {
		go func() {
			key := "data:test:exists2"
			_, err := rs.DoBool("exists", key)
			errChan <- err
		}()
	}

	for i := 0; i < l; i++ {
		err := <-errChan
		if err != nil {
			b.Fatal(err)
		}
	}
}
