package utils

import (
	"bytes"
	"encoding/json"
	"fmt"
	"runtime/debug"
)

func CatchError() {
	if err := recover(); err != nil {
		msg := fmt.Sprintf("error: %s\n%s", err, debug.Stack())
		Wlog("runtime", "罢工", msg) //这里直接写入文件。不用丢到信道了。
		fmt.Println("panic:", msg)
	}
}

//func Pr(data interface{}) {
//	fmt.Printf("%+v\n-------------------------", data)
//	//%#v\n 带结构
//}

func Pr(data interface{}) {
	//fmt.Printf("%+v\n", data)
	//%#v\n 带结构
	b, err := json.Marshal(data)
	if err != nil {
		fmt.Printf("%+v", data)
	}
	var out bytes.Buffer
	err = json.Indent(&out, b, "", "	")
	if err != nil {
		fmt.Printf("%+v", data)
	}
	fmt.Println(out.String())
	fmt.Println("------------------------")
}

func InString(v string, sl []string) bool {
	for _, vv := range sl {
		if vv == v {
			return true
		}
	}
	return false
}
