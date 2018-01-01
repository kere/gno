package i18n

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"strings"
	"text/template"
)

var (
	trunslation   = make(map[string]TrunsData)
	defaultLocale = "en-US"
)

type TrunsData map[string]string

func (t TrunsData) T(k string, args ...interface{}) string {
	if v, ok := t[k]; ok {
		if strings.ContainsAny(v, "%") {
			return fmt.Sprintf(v, args...)
		} else {
			return fmt.Sprint(v)
		}
	}
	return "undefined"
}

func (t TrunsData) TT(k string, data interface{}) string {
	if v, ok := t[k]; ok {
		t, _ := template.New("").Parse(v)

		b := bytes.Buffer{}
		t.Execute(&b, data)
		return b.String()
	}
	return "undefined"
}

func GetDefault() string {
	return defaultLocale
}

func SetDefault(locale string) {
	defaultLocale = locale
}

func Load(locale, file string) (TrunsData, error) {
	if v, ok := trunslation[file]; ok {
		return v, nil
	}

	var data TrunsData
	fileStr, err := ioutil.ReadFile(file)
	if err != nil {
		if locale == defaultLocale {
			return nil, err
		}
		// retry default
		data, err = Load(defaultLocale, file)
		if err != nil {
			return nil, err
		}
		return data, nil
	}

	err = json.Unmarshal(fileStr, &data)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	trunslation[file] = data

	return data, nil
}

func TransFunc(locale, file string) func(s string, args ...interface{}) string {
	data, err := Load(locale, file)
	if err != nil {
		data, err = Load(defaultLocale, file)
		if err != nil {
			return EmptyTransFunc
		}
		return data.T
	}

	return data.T
}

func EmptyTransFunc(k string, args ...interface{}) string {
	return "undefined"
}
