package httpd

import "io/ioutil"

var (
	rqs string
)

const (
	// PageLoadOpen open
	pageLoadOpen = `<div id="pageLoadMask" style="text-align:center;position:fixed;z-index:1000;top:0;right:0;left:0;bottom:0; background:#FFF;opacity:0.8"><span id="pageLoadText" style="position:relative;top:44%"></span></div><script>var __pLoadWd=window.pageLoadWord||'^o^';var __pLoadTxt = document.getElementById('pageLoadText');__pLoadTxt.innerText=__pLoadWd;var __pLoadTid = setInterval(()=>{ if(__pLoadTxt.innerText.length<20){__pLoadTxt.innerText+=" "+__pLoadWd}else{__pLoadTxt.innerText=__pLoadWd}}, 600);</script>`

	// PageLoadClose close
	pageLoadClose = `<script>function closePageLoad(){clearInterval(__pLoadTid);document.getElementById('pageLoadMask').style.display='none';}</script>`
)

// PageOption page
type PageOption struct {
	HasElement  bool
	HasVue      bool
	HasHeader   bool
	HasFooter   bool
	NoMainCSS   bool
	NoPageLoad  bool
	NoRequireJS bool
}

// PageInit page
func PageInit(pa *PageAttr, opt PageOption) {
	siteConf := Site.C.GetConf("site")

	// , user-scalable=no
	viewport := NewStrRender(`<meta name="viewport" content="width=device-width, initial-scale=1.0,minimum-scale=1.0,maximum-scale=1.0">`)

	pa.Head = make([]IRender, 1, 5)
	pa.CSS = make([]IRenderA, 0, 3)
	pa.JS = make([]IRenderA, 0, 5)
	pa.Top = make([]IRender, 0, 2)
	pa.Bottom = make([]IRender, 0, 4)

	pa.Head[0] = viewport
	if !opt.NoPageLoad {
		pa.Top = append(pa.Top, NewStrRender(pageLoadOpen))
	}

	// vue
	if opt.HasVue {
		pa.JS = append(pa.JS, NewJS(siteConf.DefaultString("vuejs", "vue.min.js")))
	}

	// element-ui
	if opt.HasElement {
		pa.CSS = append(pa.CSS, NewCSS(siteConf.DefaultString("elementcss", "element/index.css")))
		pa.JS = append(pa.JS, NewJS(siteConf.DefaultString("elementjs", "element/index.js")))
	}

	if !opt.NoMainCSS {
		pa.CSS = append(pa.CSS, NewCSS("main.css"))
	}

	if !opt.NoRequireJS {
		pa.JS = append(pa.JS, RequireJSWithSrc(pa, ReadRequireJS()))
	}

	pa.JSPosition = JSPositionBottom

	if opt.HasHeader {
		pa.Top = append(pa.Top, NewTemplate("app/view/_header.htm"))
	}

	if !opt.NoPageLoad {
		pa.Bottom = append(pa.Bottom, NewStrRender(pageLoadClose))
	}
	if opt.HasFooter {
		pa.Bottom = append(pa.Bottom, NewTemplate("app/view/_footer.htm"))
	}
}

var requirejs []byte

// ReadRequireJS script src
func ReadRequireJS() []byte {
	if len(requirejs) > 0 {
		return requirejs
	}

	var err error
	requirejs, err = ioutil.ReadFile("./webroot/assets/js/require.js")
	if err != nil {
		panic(err)
	}
	return requirejs
}
