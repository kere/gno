package main

import (
	"github.com/kere/gno/httpd"
	"github.com/kere/gno/httpd/example/app/openapi"
	"github.com/kere/gno/httpd/example/app/page"
	"github.com/kere/gno/httpd/example/app/upload"
	"github.com/kere/gno/httpd/example/app/websock"
)

func main() {
	httpd.Init("app/app.conf")
	// httpd.RunMode = httpd.ModePro

	httpd.Site.RegistGet("/", page.NewDefault())
	httpd.Site.RegistGet("/abc/:name", page.NewDefault())
	httpd.Site.RegistOpenAPI("/openapi/app", openapi.NewApp())

	httpd.Site.RegistWS("/ws", websock.NewWS())
	httpd.Site.RegistUpload("/upload/app", upload.NewImage())

	// httpd.RunMode = httpd.ModePro
	httpd.Site.Start()
}
