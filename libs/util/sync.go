package util

import (
	"fmt"
	"sync"

	"github.com/panjf2000/ants/v2"
)

// SyncObj class
type SyncObj struct {
	Locker sync.RWMutex
	Data   interface{}
}

// TryDo something
func (s *SyncObj) TryDo(trydo func() bool, doit func()) bool {
	s.Locker.RLock()
	if !trydo() {
		s.Locker.RUnlock()
		return false
	}
	s.Locker.RUnlock()

	s.Locker.Lock()
	if !trydo() {
		s.Locker.Unlock()
		return false
	}
	doit()
	s.Locker.Unlock()

	return true
}

// SyncMap class
type SyncMap struct {
	SyncObj
	Data map[interface{}]interface{}
}

// TryStore something
func (s *SyncMap) TryStore(key interface{}, build func() interface{}) bool {
	return s.TryDo(func() bool {
		_, isok := s.Data[key]
		return !isok
	}, func() {
		if len(s.Data) == 0 {
			s.Data = make(map[interface{}]interface{})
		}
		s.Data[key] = build()
	})
}

// TryGetAndStore something
// bool: is set
func (s *SyncMap) TryGetAndStore(key interface{}, build func() interface{}) (interface{}, bool) {
	var obj interface{}
	isdo := s.TryDo(func() bool {
		var isok bool
		obj, isok = s.Data[key]
		return !isok
	}, func() {
		if len(s.Data) == 0 {
			s.Data = make(map[interface{}]interface{})
		}
		obj = build()
		s.Data[key] = obj
	})
	return obj, isdo
}

// Load something
func (s *SyncMap) Load(key interface{}) (interface{}, bool) {
	s.Locker.RLock()
	v, isok := s.Data[key]
	s.Locker.RUnlock()
	return v, isok
}

// Store something
func (s *SyncMap) Store(key, v interface{}) {
	s.Locker.Lock()
	if len(s.Data) == 0 {
		s.Data = make(map[interface{}]interface{})
	}
	s.Data[key] = v
	s.Locker.Unlock()
}

// cptResult class
type cptResult struct {
	Index int
	Data  interface{}
	Error error
}

// Computation 并行计算
type Computation struct {
	IsSync     bool
	NumProcess int
	ErrHandler func(err error)
}

// NewComputation f
func NewComputation(n int) *Computation {
	f := func(err error) {
		fmt.Println(err)
	}
	return &Computation{NumProcess: n, ErrHandler: f}
}

// Run 并行性计算
// l：循环数量
func (c *Computation) Run(l int, execFunc func(i int)) {
	if c.IsSync {
		for k := 0; k < l; k++ {
			execFunc(k)
		}
		return
	}

	var wg sync.WaitGroup
	p, _ := ants.NewPoolWithFunc(c.NumProcess, func(i interface{}) {
		k := i.(int)
		execFunc(k)
		wg.Done()
	})

	for i := 0; i < l; i++ {
		wg.Add(1)
		p.Invoke(i)
	}
	wg.Wait()
	p.Release()
}

// RunA 并行性计算，返回运算结果
// l：循环数量
func (c *Computation) RunA(l int, execFunc func(i int) (interface{}, error), resultFunc func(i int, dat interface{})) {
	if c.IsSync {
		fmt.Println("RunA not in Pool")
		for i := 0; i < l; i++ {
			dat, err := execFunc(i)
			if err != nil {
				c.ErrHandler(err)
				continue
			}
			resultFunc(i, dat)
		}
		return
	}

	chanA := make(chan cptResult, c.NumProcess)
	p, _ := ants.NewPoolWithFunc(c.NumProcess, func(i interface{}) {
		k := i.(int)
		dat, err := execFunc(k)
		chanA <- cptResult{k, dat, err}
	})
	go func() {
		for i := 0; i < l; i++ {
			p.Invoke(i)
		}
	}()
	defer p.Release()

	var result cptResult
	for i := 0; i < l; i++ {
		result = <-chanA
		if result.Error != nil && c.ErrHandler != nil {
			c.ErrHandler(result.Error)
			continue
		}
		if resultFunc == nil {
			continue
		}
		resultFunc(result.Index, result.Data)
	}
	p.Release()
}

// // RunA 并行性计算
// // l：循环数量
// func (c *Computation) RunA(l int, execFunc func(i int) (interface{}, error), resultFunc func(i int, dat interface{})) {
// 	if l == 0 {
// 		return
// 	}
//
// 	startTime := time.Now()
// 	chanA := make(chan cptResult, c.NumProcess)
// 	go func() {
// 		for i := 0; i < l; i++ {
// 			go func(k int) {
// 				dat, err := execFunc(k)
// 				chanA <- cptResult{k, dat, err}
// 			}(i)
// 		}
// 	}()
//
// 	var result cptResult
// 	for i := 0; i < l; i++ {
// 		result = <-chanA
// 		if result.Error != nil && c.ErrHandler != nil {
// 			c.ErrHandler(result.Error)
// 			continue
// 		}
// 		if resultFunc == nil {
// 			continue
// 		}
// 		resultFunc(result.Index, result.Data)
// 	}
//
// 	fmt.Println("Finished:", time.Now().Sub(startTime))
// }
