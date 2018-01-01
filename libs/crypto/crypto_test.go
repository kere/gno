package crypto

import (
	"fmt"
	"testing"
)

func Test_crypto(t *testing.T) {
	iv := MD5([]byte("1"))
	key := MD5([]byte("a"))
	src := []byte("hello world!")

	aesCipher, err := AesEncrypt(iv, key, src)
	if err != nil {
		toError(err.Error())
	}

	text, err := AESDecrypt(iv, key, aesCipher)
	if err != nil {
		toError(err.Error())
	}
	fmt.Println(string(text))
}

//--------------------------- Benchmark ---------------------------

func Benchmark_encrypto(b *testing.B) {
	iv := MD5([]byte("1"))
	key := MD5([]byte("a"))
	src := []byte("hello world!")

	for i := 0; i < b.N; i++ {
		_, err := AesEncrypt(iv, key, src)
		if err != nil {
			toError(err.Error())
		}
	}
}

func Benchmark_decrypto(b *testing.B) {
	iv := MD5([]byte("1"))
	key := MD5([]byte("a"))
	src := []byte("hello world!")

	aesCipher, err := AesEncrypt(iv, key, src)
	if err != nil {
		toError(err.Error())
	}
	for i := 0; i < b.N; i++ {
		_, err := AESDecrypt(iv, key, aesCipher)
		if err != nil {
			toError(err.Error())
		}
	}
}

//---------------------------

func toError(err interface{}) {
	if err != nil {
		panic(err)
	}
}
