package httpd

import (
	"github.com/panjf2000/ants"
	"github.com/valyala/fasthttp"
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
	case 1:
		pageHandle(param.Page, param.Ctx)
	case 3:
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
