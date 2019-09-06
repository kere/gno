package httpd

var (
	// DefaultTemplateSubfix for render html template
	DefaultTemplateSubfix = ".htm"
)

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

	headSValNoCache = "no-cache"
	headSValMaxAge  = "max-age="
	headSValContent = "text/html; charset=utf-8"
	sAuthURL        = "url"

	// Slash string
	Slash = "/"

	stri16Formart = "%x"

	// PageLoadOpen open
	PageLoadOpen = `<div id="pageLoadMask" style="text-align:center;position:fixed;z-index:1000;top:0;right:0;left:0;bottom:0; background:#FFF;opacity:0.8"><span id="pageLoadText" style="position:relative;top:44%"></span></div><script>var __pLoadWd=window.pageLoadWord||'^o^';var __pLoadTxt = document.getElementById('pageLoadText');__pLoadTxt.innerText=__pLoadWd;var __pLoadTid = setInterval(()=>{ if(__pLoadTxt.innerText.length<20){__pLoadTxt.innerText+=" "+__pLoadWd}else{__pLoadTxt.innerText=__pLoadWd}}, 600);</script>`

	// PageLoadClose close
	PageLoadClose = `<script>function closePageLoad(){clearInterval(__pLoadTid);document.getElementById('pageLoadMask').style.display='none';}</script>`
)
