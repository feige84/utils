package utils

import "strconv"

func Int64Unique(inputInt64 []int64) (outputInt64 []int64) {
	tmpIds := make(map[int64]int64)
	for _, v := range inputInt64 {
		tmpIds[v] = v
	}
	for _, v := range tmpIds {
		outputInt64 = append(outputInt64, v)
	}
	return
}

func Int64Join(inputInt64 []int64) string {
	var ids, comm string
	for _, v := range inputInt64 {
		ids += comm + strconv.FormatInt(v, 10)
		comm = ","
	}
	return ids
}

func Int64Abs(n int64) int64 {
	y := n >> 63
	return (n ^ y) - y
}
