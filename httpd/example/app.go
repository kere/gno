package main

import (
	"github.com/kere/gno/httpd"
	"github.com/kere/gno/httpd/example/app/page"
)

func main() {
	site := httpd.Init()

	site.RegistGet("/", page.NewDefault())

	site.Start()
}
