package main

import (
	"github.com/kere/goo"
	"github.com/kere/goo/example/hello/app/page"
)

func main() {
	site := goo.Init()

	site.RegistGet("/", page.NewDefaultPage)

	site.Start()
}
