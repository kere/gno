module github.com/kere/gno/httpd

go 1.16

replace (
	github.com/kere/gno/db => D:/Work/go/project/src/github.com/kere/gno/db
	github.com/kere/gno/libs/cache => D:/Work/go/project/src/github.com/kere/gno/libs/cache
	github.com/kere/gno/libs/conf => D:/Work/go/project/src/github.com/kere/gno/libs/conf
	github.com/kere/gno/libs/crypto => D:/Work/go/project/src/github.com/kere/gno/libs/crypto
	github.com/kere/gno/libs/i18n => D:/Work/go/project/src/github.com/kere/gno/libs/i18n
	github.com/kere/gno/libs/log => D:/Work/go/project/src/github.com/kere/gno/libs/log
	github.com/kere/gno/libs/myerr => D:/Work/go/project/src/github.com/kere/gno/libs/myerr
	github.com/kere/gno/libs/redis => D:/Work/go/project/src/github.com/kere/gno/libs/redis
	github.com/kere/gno/libs/util => D:/Work/go/project/src/github.com/kere/gno/libs/util
)

require (
	github.com/buaazp/fasthttprouter v0.1.1
	github.com/fasthttp/websocket v1.4.3
	github.com/gomodule/redigo v1.8.5 // indirect
	github.com/kere/gno/db v0.0.0-00010101000000-000000000000
	github.com/kere/gno/libs/cache v0.0.0-00010101000000-000000000000
	github.com/kere/gno/libs/conf v0.0.0-00010101000000-000000000000
	github.com/kere/gno/libs/i18n v0.0.0-00010101000000-000000000000
	github.com/kere/gno/libs/log v0.0.0-00010101000000-000000000000
	github.com/kere/gno/libs/redis v0.0.0-00010101000000-000000000000 // indirect
	github.com/kere/gno/libs/util v0.0.0-00010101000000-000000000000
	github.com/valyala/bytebufferpool v1.0.0
	github.com/valyala/fasthttp v1.27.0
)
