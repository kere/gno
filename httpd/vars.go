package httpd

var (
	// DefaultTemplateSubfix for render html template
	DefaultTemplateSubfix = ".htm"
	// BDote bytes
	BDote = []byte(".")
	// BDote bytes
	BComma = []byte(",")
)

const (
	//CacheModePage 按照页面名称缓存
	CacheModePage = 1
	//CacheModePagePath 按照URL Path缓存页面
	CacheModePagePath = 2
	//CacheModePageURI 按照URL缓存页面
	CacheModePageURI = 3

	//CacheStoreMem to store in memory
	CacheStoreMem = 1
	//CacheStoreFile to store in file
	CacheStoreFile = 2
	//CacheStoreNone 不缓存页面
	CacheStoreNone = 0

	// pagecacheKeyPrefix = "c:"
	pageCacheSubfix = ".htm"

	delim1 = byte('\n')
	delim2 = "\n"

	//LastModifiedFormat Wed, 21 Oct 2015 07:28:00 GMT
	// Last-Modified: <day-name>, <day> <month> <year> <hour>:<minute>:<second> GMT
	LastModifiedFormat = "Mon, 02 Jan 2006 15:04:05 GMT"

	cacheFileStoreDir = "var/cache/page"

	headSValNoCache = "no-cache"
	headSValMaxAge  = "max-age="
	headSValContent = "text/html; charset=utf-8"
	sAuthURL        = "url"

	// Slash string
	Slash = "/"
	// Comma string
	Comma = ","

	stri16Formart = "%x"
)
