package conf

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

var (
	config   Configuration
	bComment = []byte{'#'}
	bEmpty   = []byte{}
	bEqual   = []byte{'='}
	bDQuote  = []byte{'"'}
	bBracket = []byte{'['}
)

// Load config
func Load(f string) Configuration {
	c, err := loadConfig(f)
	if err != nil {
		if os.IsNotExist(err) {

			panic(fmt.Sprint("WARN ****** conf file is not found ****** ", f, " \nCWD:", filepath.Dir(os.Args[0])))

		} else {
			println(err)
			os.Exit(1)
		}
	}
	config = c
	return c
}

// Configuration conf file
type Configuration map[string]Conf

// Conf sub item
type Conf map[string]string

func loadConfig(name string) (Configuration, error) {
	file, err := os.Open(name)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	cfg := Configuration{}

	buf := bufio.NewReader(file)
	var group string

	for {
		line, _, err := buf.ReadLine()
		if err == io.EOF {
			break
		}
		if bytes.Equal(line, bEmpty) {
			continue
		}

		if bytes.HasPrefix(line, bBracket) {
			group = string(line[1 : len(line)-1])
			continue
		}

		if bytes.HasPrefix(line, bComment) {
			continue
		}

		val := bytes.SplitN(line, bEqual, 2)

		if bytes.HasPrefix(val[1], bDQuote) {
			val[1] = bytes.Trim(val[1], `"`)
		}

		key := strings.TrimSpace(string(val[0]))

		if cfg[group] == nil {
			cfg[group] = make(map[string]string)
		}

		cfg[group][key] = strings.TrimSpace(string(val[1]))
	}

	return cfg, nil
}

// IsSet check isset
func (c Configuration) IsSet(key string) bool {
	_, ok := c[key]
	if ok && len(c[key]) > 0 {
		return true
	}
	return false
}

// GetConf return sub conf
func (c Configuration) GetConf(key string) Conf {
	if c.IsSet(key) {
		return c[key]
	}
	return Conf{}
}

// MapData return MapData
func (c Conf) MapData() map[string]string {
	return map[string]string(c)
}

// IsSet check is exists key
func (c Conf) IsSet(key string) bool {
	_, ok := c[key]
	return ok
}

// GetBool return bool
// Bool returns the boolean value for a given key.
func (c Conf) GetBool(key string) bool {
	if !c.IsSet(key) {
		return false
	}
	value, _ := strconv.ParseBool(c[key])
	return value
}

// DefaultBool bool
func (c Conf) DefaultBool(key string, v bool) bool {
	if !c.IsSet(key) {
		return v
	}

	return c.GetBool(key)
}

// GetInt return int
// Int returns the integer value for a given key.
func (c Conf) GetInt(key string) int {
	if !c.IsSet(key) {
		return 0
	}
	value, _ := strconv.Atoi(c[key])
	return value
}

// DefaultInt bool
func (c Conf) DefaultInt(key string, v int) int {
	if !c.IsSet(key) {
		return v
	}

	return c.GetInt(key)
}

// GetFloat float64
// Float returns the float value for a given key.
func (c Conf) GetFloat(key string) float64 {
	if !c.IsSet(key) {
		return 0
	}
	value, _ := strconv.ParseFloat(c[key], 64)
	return value
}

// DefaultFloat bool
func (c Conf) DefaultFloat(key string, v float64) float64 {
	if !c.IsSet(key) {
		return v
	}

	return c.GetFloat(key)
}

// GetString return string
// String returns the string value for a given key.
func (c Conf) GetString(key string) string {
	return c.Get(key)
}

// DefaultString bool
func (c Conf) DefaultString(key string, v string) string {
	if !c.IsSet(key) {
		return v
	}

	return c.GetString(key)
}

// GetStrings return []string
// returns the string list value for a given key, split by , .
// a,b,c,d
func (c Conf) GetStrings(key string) []string {
	s := c.Get(key)
	return strings.Split(s, ",")
}

// DefaultStrings []string
func (c Conf) DefaultStrings(key string, v []string) []string {
	if !c.IsSet(key) {
		return v
	}

	return c.GetStrings(key)
}

// GetInt64s return []int64
func (c Conf) GetInt64s(key string) []int64 {
	s := c.GetStrings(key)
	arr := make([]int64, len(s))
	for i, str := range s {
		arr[i], _ = strconv.ParseInt(str, 10, 64)
	}
	return arr
}

// DefaultInt64s []int64
func (c Conf) DefaultInt64s(key string, v []int64) []int64 {
	if !c.IsSet(key) {
		return v
	}

	return c.GetInt64s(key)
}

// GetInts return []int
func (c Conf) GetInts(key string) []int {
	s := c.GetStrings(key)
	l := len(s)
	arr := make([]int, l)
	for i := 0; i < l; i++ {
		arr[i], _ = strconv.Atoi(s[i])
	}
	return arr
}

// DefaultInts []int64
func (c Conf) DefaultInts(key string, v []int) []int {
	if !c.IsSet(key) {
		return v
	}

	return c.GetInts(key)
}

// GetFloats return []float64
func (c Conf) GetFloats(key string) []float64 {
	s := c.GetStrings(key)
	l := len(s)
	arr := make([]float64, l)
	for i := 0; i < l; i++ {
		arr[i], _ = strconv.ParseFloat(s[i], 64)
	}
	return arr
}

// DefaultFloats []float64
func (c Conf) DefaultFloats(key string, v []float64) []float64 {
	if !c.IsSet(key) {
		return v
	}

	return c.GetFloats(key)
}

// GetBytes return []byte
func (c Conf) GetBytes(key string) []byte {
	return []byte(c.GetString(key))
}

// DefaultBytes []int64
func (c Conf) DefaultBytes(key string, v []byte) []byte {
	if !c.IsSet(key) {
		return v
	}

	return c.GetBytes(key)
}

// JSONParse parse json to v
func (c Conf) JSONParse(field string, v interface{}) error {
	if err := json.Unmarshal(c.GetBytes(field), v); err != nil {
		return err
	}
	return nil
}

// Get func
// String returns the string value for a given key.
func (c Conf) Get(key string) string {
	if c.IsSet(key) {
		return c[key]
	}
	return ""
}

// GetRunMode return framwork run mode
func (c Configuration) GetRunMode() string {
	var m string

	if c.IsSet("app") && c["app"].IsSet("mode") {
		m = c["app"]["mode"]

		switch m {
		case "pro":
		case "production":
			m = "pro"
		default:
			m = "dev"
		}
	} else {
		m = "dev"
	}

	return m
}
