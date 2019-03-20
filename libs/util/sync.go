package util

import (
	"fmt"
	"time"
)

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
func (c *Computation) Run(l int, execFunc func(i int)) {
	if l == 0 {
		return
	}

	if c.ErrHandler == nil {
		c.ErrHandler = func(err error) {
			fmt.Println(err)
		}
	}

	startTime := time.Now()
	chanA := make(chan int8, c.NumProcess)
	for i := 0; i < l; i++ {
		chanA <- 1
		go func(i int) {
			execFunc(i)
			<-chanA
		}(i)
	}

	fmt.Println("Finished:", time.Now().Sub(startTime))
}
