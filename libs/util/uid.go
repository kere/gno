package util

import "hash/crc32"

// IID32 return iid
func IID32(str ...string) int64 {
	ieee := crc32.NewIEEE()
	count := len(str)
	for i := 0; i < count; i++ {
		ieee.Write(Str2Bytes(str[i]))
	}
	return int64(ieee.Sum32())
}
