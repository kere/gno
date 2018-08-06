package util

import (
	"fmt"
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

// Run f
func (c *Computation) Run(l int, execFunc func(i int, c *Computation), resultFunc func(dat interface{}, c *Computation)) {
	if l == 0 {
		return
	}
	if c.ErrHandler == nil {
		c.ErrHandler = func(err error) {
			fmt.Println(err)
		}
	}

	go func() {
		for i := 0; i < l; i++ {
			go execFunc(i, c)
		}
	}()

	for c.counter < l {
		select {
		case dat := <-c.dataChan:
			resultFunc(dat, c)
			c.counter++

		case err := <-c.errChan:
			c.ErrHandler(err)
			c.counter++

		}
	}
}

// Done f
func (c *Computation) Done(dat interface{}) {
	c.dataChan <- dat
}

// DoError err
func (c *Computation) DoError(err error) {
	c.errChan <- err
}
