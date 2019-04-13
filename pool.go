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
	// Error  chan error
}

// InvokeExec by http
func InvokeExec(a interface{}) {
	ps := a.(PoolParams)

	var err error
	switch ps.Typ {
	case 1:
		err = pageHandle(ps.Page)
		if err != nil {
			doPageError(Site.ErrorURL, err, ps.RW, ps.Req)
		}
	case 2:
		err = openAPIHandle(ps.RW, ps.Req, ps.Params)
		if err != nil {
			doAPIError(err, ps.RW, ps.Req)
		}
	case 3:
		err = webAPIHandle(ps.WebAPI, ps.RW, ps.Req, ps.Params)
		if err != nil {
			doAPIError(err, ps.RW, ps.Req)
		}
	}
}

// NewPool new
func NewPool(n int) *ants.PoolWithFunc {
	po, err := ants.NewPoolWithFunc(n, InvokeExec)
	if err != nil {
		panic(err)
	}
	return po
}
