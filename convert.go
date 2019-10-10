package utils

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

// GetString convert interface to string.
func GetString(v interface{}) string {
	switch result := v.(type) {
	case string:
		return result
	case []byte:
		return string(result)
	default:
		if v != nil {
			return fmt.Sprint(result)
		}
	}
	return ""
}

// GetInt convert interface to int.
func GetInt(v interface{}) int {
	switch result := v.(type) {
	case int:
		return result
	case int32:
		return int(result)
	case int64:
		return int(result)
	default:
		if d := GetString(v); d != "" {
			value, err := strconv.Atoi(d)
			if err != nil {
				panic(err.Error())
			}
			return value
		}
	}
	return 0
}

// GetInt64 convert interface to int64.
func GetInt64(v interface{}) int64 {
	switch result := v.(type) {
	case int:
		return int64(result)
	case int32:
		return int64(result)
	case int64:
		return result
	default:

		if d := GetString(v); d != "" {
			value, err := strconv.ParseInt(d, 10, 64)
			if err != nil {
				panic(err.Error())
			}
			return value
		}
	}
	return 0
}

// GetFloat64 convert interface to float64.
func GetFloat64(v interface{}) float64 {
	switch result := v.(type) {
	case float64:
		return result
	default:
		if d := GetString(v); d != "" {
			value, err := strconv.ParseFloat(d, 64)
			if err != nil {
				panic(err.Error())
			}
			return value
		}
	}
	return 0
}

// GetBool convert interface to bool.
func GetBool(v interface{}) bool {
	switch result := v.(type) {
	case bool:
		return result
	default:
		if d := GetString(v); d != "" {
			value, err := strconv.ParseBool(d)
			if err != nil {
				panic(err.Error())
			}
			return value
		}
	}
	return false
}

func StrToInt(str string) int64 {
	if len(str) == 0 {
		return 0
	}
	result, err := strconv.ParseInt(str, 10, 64)
	if err != nil {
		return 0
	}
	return result
}

func StrToFloat(str string) float64 {
	if len(str) == 0 {
		return 0
	}
	result, err := strconv.ParseFloat(str, 64)
	if err != nil {
		panic(err.Error())
	}
	return result
}

func TrimHtml(src string) string {
	//将HTML标签全转换成小写
	re, _ := regexp.Compile("\\<[\\S\\s]+?\\>")
	src = re.ReplaceAllStringFunc(src, strings.ToLower)
	//去除STYLE
	re, _ = regexp.Compile("\\<style[\\S\\s]+?\\</style\\>")
	src = re.ReplaceAllString(src, "")
	//去除SCRIPT
	re, _ = regexp.Compile("\\<script[\\S\\s]+?\\</script\\>")
	src = re.ReplaceAllString(src, "")
	//去除所有尖括号内的HTML代码，并换成换行符
	re, _ = regexp.Compile("\\<[\\S\\s]+?\\>")
	src = re.ReplaceAllString(src, "\n")
	//去除连续的换行符
	re, _ = regexp.Compile("\\s{2,}")
	src = re.ReplaceAllString(src, "\n")
	return strings.TrimSpace(src)
}

func Base64Encode(str string) string {
	if str != "" {
		return base64.StdEncoding.EncodeToString([]byte(str))
	}
	return str
}

// Base64Decode base64_decode()
func Base64Decode(str string) string {
	if str != "" {
		data, err := base64.StdEncoding.DecodeString(str)
		if err == nil {
			return string(data)
		}
		return str
	}
	return ""
}

/*
func UTF82GBK(str string) string {
	//将utf-8编码的字符串转换为GBK编码
	ret, _ := simplifiedchinese.GBK.NewEncoder().String(str)
	return ret //如果转换失败返回空字符串
	//如果是[]byte格式的字符串，可以使用Bytes方法
	//b, err := simplifiedchinese.GBK.NewEncoder().Bytes([]byte(str))
	//return string(b)
}

func GBK2UTF8(gbkStr string) string {
	//将GBK编码的字符串转换为utf-8编码
	ret, _ := simplifiedchinese.GBK.NewDecoder().String(gbkStr)
	return ret //如果转换失败返回空字符串

	//如果是[]byte格式的字符串，可以使用Bytes方法
	//b, err := simplifiedchinese.GBK.NewDecoder().Bytes([]byte(gbkStr))
	//return string(b)
}
*/

func Addslashes(str string) string {
	var buf bytes.Buffer
	for _, char := range str {
		switch char {
		case '\'', '"', '\\':
			buf.WriteRune('\\')
		}
		buf.WriteRune(char)
	}
	return buf.String()
}

func HideStar(str string) (result string) {
	if str == "" {
		return ""
	}
	if strings.Contains(str, "@") {
		res := strings.Split(str, "@")
		if len(res[0]) < 3 {
			resString := "***"
			result = resString + "@" + res[1]
		} else {
			res2 := SubStr2(str, 0, 3)
			resString := res2 + "***"
			result = resString + "@" + res[1]
		}
		return result
	} else {
		reg := `^1[0-9]\d{9}$`
		rgx := regexp.MustCompile(reg)
		mobileMatch := rgx.MatchString(str)
		if mobileMatch {
			result = SubStr2(str, 0, 3) + "****" + SubStr2(str, 7, 11)
		} else {
			nameRune := []rune(str)
			lens := len(nameRune)
			if lens <= 1 {
				result = "***"
			} else if lens == 2 {
				result = string(nameRune[:1]) + "*"
			} else if lens == 3 {
				result = string(nameRune[:1]) + "*" + string(nameRune[2:3])
			} else if lens == 4 {
				result = string(nameRune[:1]) + "**" + string(nameRune[lens-1:lens])
			} else if lens > 4 {
				result = string(nameRune[:2]) + "***" + string(nameRune[lens-2:lens])
			}
		}
		return
	}
}

//四舍五入，保留2位
func Round(number float64) float64 {
	numberStr := fmt.Sprintf("%.2f", number) //四舍五入，保留2位
	floatNum, err := strconv.ParseFloat(numberStr, 64)
	if err != nil {
		panic(err.Error())
	}
	return floatNum
}
