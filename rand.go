package utils

import (
	"fmt"
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

//生成随机字符串
func GetRandomSmallStr(length int) string {
	str := "abcdefghijklmnopqrstuvwxyz"
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

//生成订单号，传入content可传用户ID或邀请码等可代表身份的标识。也可留空自动生成
func GenerateOrderNo(content string) string {
	now := GetNow()
	if content == "" {
		nanoStr := fmt.Sprint(now.UnixNano())
		content = nanoStr[8:10] + nanoStr[12:14]
	}
	lastStr := ""
	contentLen := len(content)
	if contentLen >= 4 { //邀请码取后位数，不取全是为了让别人猜不到完整邀请码
		lastStr = content[contentLen-4:]
	} else {
		//不足4位用0补齐4位
		lastStr = fmt.Sprintf("%04v", content)
	}
	nano := fmt.Sprint(now.Nanosecond())[:6]
	//年月日时分秒+纳秒前6位+邀请码后4位
	//14位+6位+4位=24位
	return fmt.Sprintf("%s%s%s", now.Format("20060102150405"), nano, lastStr) //生成自己的订单号
}
