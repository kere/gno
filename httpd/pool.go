package httpd

import (
	"github.com/panjf2000/ants"
	"github.com/valyala/fasthttp"
)

const (
	invokePage = 1
	invokeAPI  = 3
)

var pool *ants.PoolWithFunc

// PoolParams for pool
type PoolParams struct {
	Typ  int
	Page IPage
	Ctx  *fasthttp.RequestCtx
	Done chan bool
}

// InvokeExec by http
func InvokeExec(dat interface{}) {
	param := dat.(PoolParams)
	switch param.Typ {
	case invokePage:
		pageHandle(param.Page, param.Ctx)
	case invokeAPI:
		openAPIHandle(param.Ctx)
	}
	param.Done <- true
}

// initPool new
func initPool(n int) {
	var err error
	pool, err = ants.NewPoolWithFunc(n, InvokeExec)
	if err != nil {
		panic(err)
	}
}
