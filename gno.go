package gno

import (
	"github.com/kere/gno/httpd"
	"github.com/kere/gno/libs/conf"
)

var (
	// HomeDir home
	HomeDir = ""

	// RunMode home
	RunMode = ""
)

// GetConfig return Configuration
func GetConfig() conf.Configuration {
	return httpd.Site.C
}

// Init gno
func Init(name string) {
	httpd.Init(name)
	HomeDir = httpd.HomeDir
	RunMode = httpd.RunMode
}
