package main

import (
	"github.com/kere/gno/httpd"
	"github.com/kere/gno/httpd/example/app/openapi"
	"github.com/kere/gno/httpd/example/app/page"
	"github.com/kere/gno/httpd/example/app/websock"
)

func main() {
	httpd.Init("app/app.conf")

	httpd.Site.RegistGet("/", page.NewDefault())
	httpd.Site.RegistGet("/abc/:name", page.NewDefault())
	httpd.Site.RegistOpenAPI("/openapi/app", openapi.NewApp())

	httpd.Site.RegistWS("/ws", websock.NewWS())

	// httpd.RunMode = httpd.ModePro
	httpd.Site.Start()
}
