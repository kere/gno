package gno

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/panjf2000/ants"
)

var pool *ants.PoolWithFunc

// PoolParams for pool
type PoolParams struct {
	Typ    int
	Page   IPage
	WebAPI IWebAPI
	RW     http.ResponseWriter
	Req    *http.Request
	Params httprouter.Params
	Error  chan error
}

// NewPool new
func NewPool(n int) *ants.PoolWithFunc {
	po, err := ants.NewPoolWithFunc(n, func(a interface{}) {
		ps, ok := a.(*PoolParams)
		var err error
		if !ok {
			ps.Error <- err
			return
		}
		switch ps.Typ {
		case 1:
			err = doPageHandle(ps.Page, ps.RW, ps.Req, ps.Params)
		case 2:
			err = openAPIHandle(ps.RW, ps.Req, ps.Params)
		case 3:
			err = doAPIHandle(ps.WebAPI, ps.RW, ps.Req, ps.Params)
		}

		ps.Error <- err
	})
	if err != nil {
		panic(err)
	}
	return po
}
