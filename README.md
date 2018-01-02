## Installation
Make sure you have the a working Go environment. See the [install instructions](http://golang.org/doc/install.html).

And install gos:

	go get github.com/kere/gos


gos project
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
