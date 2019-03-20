package util

import (
	"fmt"
	"time"
)

// Computation 并行计算
type Computation struct {
	NumProcess int
	dataChan   chan interface{}
	errChan    chan error
	counter    int

	ErrHandler func(err error)
}

// NewComputation f
func NewComputation(num int) *Computation {
	return &Computation{NumProcess: num, dataChan: make(chan interface{}, num), errChan: make(chan error, num)}
}

// Run 并行性计算
// l：循环数量
func (c *Computation) Run(l int, execFunc func(i int) (interface{}, error), resultFunc func(dat interface{})) {
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
	go func() {
		for i := 0; i < l; i++ {
			chanA <- 1
			go func(i int) {
				dat, err := execFunc(i)
				if err != nil {
					c.errChan <- err
					return
				}
				c.dataChan <- dat
			}(i)
		}
	}()

	for c.counter < l {
		select {
		case dat := <-c.dataChan:
			resultFunc(dat)
			c.counter++
			<-chanA

		case err := <-c.errChan:
			c.ErrHandler(err)
			c.counter++
			<-chanA

		}
	}

	fmt.Println("Finished:", time.Now().Sub(startTime))
}

// GetCounter f
func (c *Computation) GetCounter() int {
	return c.counter
}

// // Done f
// func (c *Computation) Done(dat interface{}) {
// 	c.dataChan <- dat
// }

// // DoError err
// func (c *Computation) DoError(err error) {
// 	c.errChan <- err
// }
