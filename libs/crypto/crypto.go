package crypto

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/md5"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha1"
	"encoding/base64"
	"errors"
	"fmt"
	"io"
	// "strconv"
)

func Base64Decode(s []byte) ([]byte, error) {
	dbuf := make([]byte, base64.StdEncoding.DecodedLen(len(s)))
	n, err := base64.StdEncoding.Decode(dbuf, s)
	return dbuf[:n], err
}

func Base64Encode(src []byte) []byte {
	buf := make([]byte, base64.StdEncoding.EncodedLen(len(src)))
	base64.StdEncoding.Encode(buf, src)
	return buf
}

func Sha1(str string) string {
	h := sha1.New()
	io.WriteString(h, "abc")
	return fmt.Sprintf("%x", h.Sum(nil))
}

func MD5String(str string) string {
	return fmt.Sprintf("%x", MD5([]byte(str)))
}

func MD5(b []byte) []byte {
	h := md5.New()
	h.Write(b)
	return h.Sum(nil)
}

func AESEncrypt(iv, key, src []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	src = PKCS5Padding(src, block.BlockSize())

	crypted := make([]byte, len(src))
	cipher.NewCBCEncrypter(block, iv).CryptBlocks(crypted, src)
	return crypted, nil
}

func AESDecrypt(iv, key, aesCipher []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	l := len(aesCipher)
	bl := block.BlockSize()
	if l < bl || l%bl != 0 {
		return nil, errors.New("aes cipher is not a multiple of the block size:" + fmt.Sprint(len(aesCipher)))
	}
	plaintext := make([]byte, len(aesCipher))
	cipher.NewCBCDecrypter(block, iv).CryptBlocks(plaintext, aesCipher)

	return PKCS5UnPadding(plaintext)
}

func NewRSAKey() *rsa.PrivateKey {
	key, err := rsa.GenerateKey(rand.Reader, 512)
	if err != nil {
		return nil
	}
	return key
}

func PKCS5Padding(ciphertext []byte, blockSize int) []byte {
	padding := blockSize - len(ciphertext)%blockSize
	padtext := bytes.Repeat([]byte{byte(padding)}, padding)

	return append(ciphertext, padtext...)
}

func PKCS5UnPadding(origData []byte) ([]byte, error) {
	length := len(origData)
	if length%16 > 0 {
		return nil, fmt.Errorf("PKCS5UnPadding error, data length is %d", length)
	}

	unpadding := int(origData[length-1])
	if length < unpadding {
		return origData, fmt.Errorf("PKCS5UnPadding error")
	}
	return origData[:(length - unpadding)], nil
}
