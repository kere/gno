package util

import (
	"fmt"
	"time"
)

// CptResult class
type CptResult struct {
	Index int
	Data  interface{}
	Error error
}

// Computation 并行计算
type Computation struct {
	NumProcess int
	ErrHandler func(err error)
}

// NewComputation f
func NewComputation(num int) *Computation {
	return &Computation{NumProcess: num}
}

// Run 并行性计算
// l：循环数量
func (c *Computation) Run(l int, execFunc func(i int) (interface{}, error), resultFunc func(i int, dat interface{})) {
	if l == 0 {
		return
	}

	if c.ErrHandler == nil {
		c.ErrHandler = func(err error) {
			fmt.Println(err)
		}
	}

	startTime := time.Now()
	chanA := make(chan CptResult, c.NumProcess)
	go func() {
		for i := 0; i < l; i++ {
			go func(k int) {
				dat, err := execFunc(k)
				chanA <- CptResult{k, dat, err}
			}(i)
		}
	}()

	var result CptResult
	for i := 0; i < l; i++ {
		fmt.Println("--==--", i, l)
		result = <-chanA
		if result.Error != nil && c.ErrHandler != nil {
			c.ErrHandler(result.Error)
			continue
		}
		resultFunc(result.Index, result.Data)
	}

	fmt.Println("Finished:", time.Now().Sub(startTime))
}
