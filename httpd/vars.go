package httpd

const (
	//CacheModePage 按照页面名称缓存
	CacheModePage = 1
	//CacheModePagePath 按照URL Path缓存页面
	CacheModePagePath = 2
	//CacheModePageURI 按照URL缓存页面
	CacheModePageURI = 3

	//CacheStoreMem to store in memory
	CacheStoreMem = 0
	//CacheStoreFile to store in file
	CacheStoreFile = 1
	//CacheStoreNone 不缓存页面
	CacheStoreNone = -1

	// pagecacheKeyPrefix = "c:"
	pageCacheSubfix = ".htm"

	delim1 = byte('\n')
	delim2 = "\n"

	//LastModifiedFormat Wed, 21 Oct 2015 07:28:00 GMT
	// Last-Modified: <day-name>, <day> <month> <year> <hour>:<minute>:<second> GMT
	LastModifiedFormat = "Mon, 02 Jan 2006 15:04:05 GMT"

	cacheFileStoreDir = "var/cache/page"

	headSValNoCache       = "no-cache"
	headSValMaxAge        = "max-age="
	headSValContent       = "text/html; charset=utf-8"
	defaultTemplateSubfix = ".htm"
	sAuthURL              = "url"

	// Slash string
	Slash = "/"

	stri16Formart = "%x"
)
