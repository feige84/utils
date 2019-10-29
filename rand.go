package utils

import (
	"math"
	"math/rand"
)

//生成随机字符串
func GetRandomString(length int) string {
	str := "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	bytes := []byte(str)
	result := []byte{}
	r := rand.New(rand.NewSource(GetNow().UnixNano()))
	for i := 0; i < length; i++ {
		result = append(result, bytes[r.Intn(len(bytes))])
	}
	return string(result)
}

//生成十六进制随机
func GetRandomHex(length int) string {
	str := "0123456789abcdef"
	bytes := []byte(str)
	result := []byte{}
	r := rand.New(rand.NewSource(GetNow().UnixNano()))
	for i := 0; i < length; i++ {
		result = append(result, bytes[r.Intn(len(bytes))])
	}
	return string(result)
}

//生成随机数字,用于验证码
func GetRandomNum(length int) string {
	str := "0123456789"
	bytes := []byte(str)
	result := []byte{}
	r := rand.New(rand.NewSource(GetNow().UnixNano()))
	for i := 0; i < length; i++ {
		result = append(result, bytes[r.Intn(len(bytes))])
	}
	return string(result)
}

//取一个范围的随机数
func Rand(min, max int) int {
	if min > max {
		panic("min: min cannot be greater than max")
	}
	r := rand.New(rand.NewSource(GetNow().UnixNano()))
	n := r.Intn(math.MaxInt32)
	return n/((math.MaxInt32+1)/(max-min+1)) + min
}
