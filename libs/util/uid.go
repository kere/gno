package util

import (
	"hash/crc64"
	"time"
)

// Unique 获得一个稀有字符
func Unique() string {
	t := time.Now().UTC().UnixNano()

	tab := crc64.MakeTable(uint64(t))
	i := crc64.Checksum(RandStr(12), tab)
	return string(IntZipTo62(i))
}
