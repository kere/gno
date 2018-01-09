## Installation
Make sure you have the a working Go environment. See the [install instructions](http://golang.org/doc/install.html).

And install gno:

	go get github.com/kere/gno


gno project
```
|- [app]
	└── app.conf
	├── app.go
	├── [api]
	├── [page]
	├── [view]
|- [webroot]
	├── [assets]
		├── [js]
		├── [css]
```
##
go get github.com/youtube/vitess/go/pools

## quick start
app.go
```go
package main

import (
	"github.com/kere/goo"
	"github.com/kere/goo/example/hello/app/page"
)

func main() {
	site := gno.Init()

	site.RegistGet("/", page.NewDefaultPage)

	site.Start()
}
```
