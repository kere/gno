package httpd

// const (
// 	invokePage = 1
// 	invokeAPI  = 3
// )
//
// var invokeDoneOK = struct{}{}
// var pool *ants.PoolWithFunc
//
// // PoolParams for pool
// type PoolParams struct {
// 	Typ  int
// 	Page IPage
// 	Ctx  *fasthttp.RequestCtx
// 	Done chan struct{}
// }
//
// // InvokeExec by http
// func InvokeExec(dat interface{}) {
// 	param := dat.(PoolParams)
// 	switch param.Typ {
// 	case invokePage:
// 		pageHandle(param.Page, param.Ctx)
// 	case invokeAPI:
// 		openAPIHandle(param.Ctx)
// 	}
// 	close(param.Done)
// 	// param.Done <- invokeDoneOK
// }
//
// // initPool new
// func initPool(n int) {
// 	if n == 0 {
// 		return
// 	}
// 	var err error
// 	pool, err = ants.NewPoolWithFunc(n, InvokeExec)
// 	if err != nil {
// 		panic(err)
// 	}
// }
