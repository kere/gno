package main

import (
	"github.com/kere/gno/httpd"
	"github.com/kere/gno/httpd/example/app/openapi"
	"github.com/kere/gno/httpd/example/app/page"
	"github.com/kere/gno/httpd/example/app/websock"
)

func main() {
	site := httpd.Init()

	site.RegistGet("/", page.NewDefault())
	site.RegistGet("/abc/:name", page.NewDefault())
	site.RegistOpenAPI("/openapi/app", openapi.NewApp())

	site.RegistWS("/ws", websock.NewWS())

	// httpd.RunMode = httpd.ModePro
	site.Start()
}
