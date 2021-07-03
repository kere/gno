module github.com/kere/gno/db

go 1.16

replace (
	github.com/kere/gno/httpd => D:/Work/go/project/src/github.com/kere/gno/httpd
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
	github.com/djimenez/iconv-go v0.0.0-20160305225143-8960e66bd3da
	github.com/kere/gno/libs/conf v0.0.0-00010101000000-000000000000
	github.com/kere/gno/libs/log v0.0.0-00010101000000-000000000000
	github.com/kere/gno/libs/myerr v0.0.0-00010101000000-000000000000
	github.com/kere/gno/libs/util v0.0.0-00010101000000-000000000000
	github.com/lib/pq v1.10.2
	github.com/valyala/bytebufferpool v1.0.0
)
