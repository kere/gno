package main

import (
	"github.com/kere/gno/httpd"
	"github.com/kere/gno/httpd/example/app/openapi"
	"github.com/kere/gno/httpd/example/app/page"
)

func main() {
	site := httpd.Init()

	site.RegistGet("/", page.NewDefault())
	site.RegistGet("/abc", page.NewDefault())
	site.RegistOpenAPI("/openapi/app", openapi.NewApp())

	site.Start()
}
