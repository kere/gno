package main

import (
	"github.com/kere/gno"
	"github.com/kere/gno/example/hello/app/openapi"
	"github.com/kere/gno/example/hello/app/page"
	wb "github.com/kere/gno/example/hello/app/websock"
)

func main() {
	site := gno.Init()

	site.RegistGet("/", page.NewDefaultPage)

	site.RegistOpenAPI("/openapi/app", openapi.NewApp())

	// site.RegistMessageSocket("/ws", wb.NewMessage())
	site.RegistWebSocket("/ws", wb.NewUser())

	site.Start()
}
