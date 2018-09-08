package main

import (
	"github.com/kere/gno"
	"github.com/kere/gno/example/hello/app/openapi"
	"github.com/kere/gno/example/hello/app/page"
)

func main() {
	site := gno.Init()

	site.RegistGet("/", page.NewDefaultPage)

	gno.Site.RegistOpenAPI("/openapi/app", openapi.NewApp())
	site.Start()
}
