package utils

import (
	"bytes"
	"encoding/json"
)

func JsonMarshal(v interface{}) string {
	buffer := bytes.NewBuffer([]byte{})
	jsonEncoder := json.NewEncoder(buffer)
	jsonEncoder.SetEscapeHTML(false)
	err := jsonEncoder.Encode(v)
	if err != nil {
		panic(err.Error())
	}
	return string(buffer.String())
}
