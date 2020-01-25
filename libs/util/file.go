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

// // DirItems read dir
// func DirItems(dir string, f func(os.FileInfo)) error {
// 	items, err := ioutil.ReadDir(dir)
// 	if err != nil {
// 		return err
// 	}
//
// 	l := len(items)
// 	for i := 0; i < l; i++ {
// 		f(items[i])
// 	}
// 	return nil
// }
