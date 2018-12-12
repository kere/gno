package util

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
)

// AjaxGet 获取Get
func AjaxGet(uri string, vals url.Values) ([]byte, error) {
	params := ""
	if len(vals) > 0 {
		params = "?" + vals.Encode()
	}

	resq, err := http.Get(uri + params)
	if err != nil {
		return nil, err
	}

	defer resq.Body.Close()
	return ioutil.ReadAll(resq.Body)
}

// AjaxPost 获取Post
func AjaxPost(uri string, dat MapData) ([]byte, error) {
	vals := url.Values{}
	for k, v := range dat {
		vals[k] = []string{fmt.Sprint(v)}
	}

	resq, err := http.PostForm(uri, vals)
	if err != nil {
		return nil, err
	}

	defer resq.Body.Close()
	return ioutil.ReadAll(resq.Body)
}
