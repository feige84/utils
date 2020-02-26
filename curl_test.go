package utils

import (
	"fmt"
	"testing"
)

func TestCURL(t *testing.T) {
	req := &HttpSend{}
	req.ConnectTimeout = 5
	req.Header = map[string]string{
		"User-Agent":      "Mozilla/5.0 (Linux; U; Android 4.4.1; zh-cn; NX507J Build/KVT49L) AppleWebKit/533.1 (KHTML, like Gecko)Version/4.0 MQQBrowser/5.4 TBS/025469 Mobile Safari/533.1 MicroMessenger/6.2.5.53_r2565f18.621 NetType/WIFI Language/zh_CN",
		"Connection":      "Keep-Alive",
		"Cache-Control":   "max-stale=0",
		"Accept-Encoding": "gzip",
		"referer":         "http://www.baidu.com",
	}
	req.RequestUrl = "http://www.baidu.com"

	ret, err := HttpHandle(req)
	if err != nil {
		panic(err)
	}
	fmt.Println("String:", ret.String())
	fmt.Println("Bytes:", ret.Bytes())
	fmt.Println("Cookie:", ret.Cookie())
	fmt.Println("Redirect:", ret.Redirect())
}
