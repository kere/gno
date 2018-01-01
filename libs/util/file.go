package util

import (
	"io"
	"os"
)

// CopyFile func
func CopyFile(srcName, dstName string) (written int64, err error) {
	src, err := os.Open(srcName)
	if err != nil {
		return 0, err
	}
	defer src.Close()
	dst, err := os.OpenFile(dstName, os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		return 0, err
	}
	defer dst.Close()

	return io.Copy(dst, src)
}
