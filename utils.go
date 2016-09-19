package goquery

import (
	"bytes"
	"encoding/json"
)

type DataHolder map[string] interface{}

func ToJson(input interface{}, pretty bool) string {
	var prefix, indent = "", ""
	if pretty {
		indent = "    "
	}
	buffer := bytes.Buffer{}
	encoder := json.NewEncoder(&buffer)
	encoder.SetIndent(prefix, indent)
	err := encoder.Encode(input)
	if err != nil {
		return "{}"
	}
	return string(buffer.Bytes())
}

func (data DataHolder) GetFirstChild(key string, index int) map[string]interface{} {
	val, ok := data[key].([]interface{})
	if ok {
		arr, ok := val[index].(map[string]interface{})
		if ok {
			return arr
		}
	}
	return DataHolder{}
}

func (data DataHolder) GetStrings(key string) []string {
	val, ok := data[key].([]string)
	if ok {
		return val
	}
	return []string{}
}

func (data DataHolder) GetString(key string) string {
	val, ok := data[key].(string)
	if ok {
		return val
	}
	return ""
}

func (data DataHolder) GetInt(key string) int {
	val, ok := data[key].(int)
	if ok {
		return val
	}
	return 0
}