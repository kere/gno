package util

import (
	"crypto/md5"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"time"
)

func AjaxGet(uri string, dat MapData) ([]byte, error) {
	params := ""
	if dat != nil {
		vals := url.Values{}
		for k, v := range dat {
			vals.Add(k, fmt.Sprint(v))
		}
		params = "?" + vals.Encode()
	}
	resq, err := http.Get(uri + params)
	if err != nil {
		return nil, err
	}

	defer resq.Body.Close()
	return ioutil.ReadAll(resq.Body)
}

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

func AjaxSend(uri string, method string, dat MapData) (MapData, error) {
	// data:       {'_src': jsonStr, 'now': now, 'token': md5(str), 'method': method},
	// str = now+method+now+jsonStr+now;
	src, err := json.Marshal(dat)
	if err != nil {
		return nil, err
	}

	ts := fmt.Sprint(time.Now().Unix())
	bstr := []byte(ts + method + ts + string(src) + ts)
	token := fmt.Sprintf("%x", md5.Sum(bstr))

	vals := url.Values{}
	vals.Add("now", ts)
	vals.Add("_src", string(src))
	vals.Add("token", token)
	vals.Add("method", method)

	resq, err := http.PostForm(uri+"/"+method, vals)
	if err != nil {
		return nil, err
	}

	defer resq.Body.Close()

	body, err := ioutil.ReadAll(resq.Body)
	if err != nil {
		return nil, err
	}
	if resq.StatusCode != http.StatusOK {
		return nil, errors.New(string(body) + " " + uri + "/" + method)
	}

	var obj MapData
	err = json.Unmarshal(body, &obj)

	return obj, err
}
