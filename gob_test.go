package utils

import (
	"fmt"
	"testing"
)

func TestExecute(t *testing.T) {
	result := make(map[string]string)
	result["aaa"] = "bbbb"
	aa := GobEncode(result)
	fmt.Println(aa)

	result2 := make(map[string]string)
	if err := GobDecode(aa, &result2);err != nil {
		panic(err)
	}
	fmt.Println(result2)
}