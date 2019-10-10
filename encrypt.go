package utils

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/des"
	"crypto/md5"
	"encoding/base64"
	"encoding/hex"
)

//PKCS7Padding 的填充方式和PKCS5Padding 填充方式一样。只是加密块的字节数不同。PKCS5Padding明确定义了加密块是8字节，PKCS7Padding加密快可以是1-255之间。

func PKCS5Padding(ciphertext []byte, blockSize int) []byte {
	padding := blockSize - len(ciphertext)%blockSize
	padtext := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(ciphertext, padtext...)
}

func PKCS5UnPadding(origData []byte) []byte {
	length := len(origData)
	unpadding := int(origData[length-1])
	return origData[:(length - unpadding)]
}

//补码
func PKCS7Padding(ciphertext []byte, blocksize int) []byte {
	padding := blocksize - len(ciphertext)%blocksize
	padtext := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(ciphertext, padtext...)
}

//去码
func PKCS7UnPadding(origData []byte) []byte {
	length := len(origData)
	unpadding := int(origData[length-1])
	return origData[:(length - unpadding)]
}

func DesEncrypt(orig, key string) (string, error) {
	origData := []byte(orig)
	k := []byte(key)
	block, err := des.NewCipher(k)
	if err != nil {
		return "", err
	}
	origData = PKCS5Padding(origData, block.BlockSize())
	// origData = ZeroPadding(origData, block.BlockSize())
	blockMode := cipher.NewCBCEncrypter(block, k)
	crypted := make([]byte, len(origData))
	// 根据CryptBlocks方法的说明，如下方式初始化crypted也可以
	// crypted := origData
	blockMode.CryptBlocks(crypted, origData)
	return base64.StdEncoding.EncodeToString(crypted), nil
}

func DesDecrypt(crypted, key string) (string, error) {
	crytedByte, _ := base64.StdEncoding.DecodeString(crypted)
	k := []byte(key)
	block, err := des.NewCipher(k)
	if err != nil {
		return "", err
	}
	blockMode := cipher.NewCBCDecrypter(block, k)
	origData := make([]byte, len(crytedByte))
	// origData := crypted
	blockMode.CryptBlocks(origData, crytedByte)
	origData = PKCS5UnPadding(origData)
	// origData = ZeroUnPadding(origData)
	return string(origData), nil
}

func AesEncrypt(orig, key string) (string, error) {
	origData := []byte(orig)
	k := []byte(key)
	block, err := aes.NewCipher(k)
	if err != nil {
		return "", err
	}
	blockSize := block.BlockSize()
	origData = PKCS5Padding(origData, blockSize)
	blockMode := cipher.NewCBCEncrypter(block, k[:blockSize])
	crypted := make([]byte, len(origData))
	blockMode.CryptBlocks(crypted, origData)
	return base64.StdEncoding.EncodeToString(crypted), nil
}

func AesDecrypt(crypted, key string) (string, error) {
	crytedByte, _ := base64.StdEncoding.DecodeString(crypted)
	k := []byte(key)
	block, err := aes.NewCipher(k)
	if err != nil {
		return "", err
	}
	blockSize := block.BlockSize()
	blockMode := cipher.NewCBCDecrypter(block, k[:blockSize])
	origData := make([]byte, len(crytedByte))
	blockMode.CryptBlocks(origData, crytedByte)
	origData = PKCS5UnPadding(origData)
	return string(origData), nil
}

func MD5(text string) string {
	ctx := md5.New()
	ctx.Write([]byte(text))
	return hex.EncodeToString(ctx.Sum(nil))
}
