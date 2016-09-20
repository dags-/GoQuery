package goquery

import (
	"bytes"
	"encoding/json"
)

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

type Data map[string] interface{}

func (data Data) Put(key string, val interface{}) {
	if key != "" && val != "" && val != nil {
		data[key] = val
	}
}

func (data Data) Retain(keys ...string) Data {
	result := Data{}
	for _, k := range keys {
		result.Put(k, data[k])
	}
	return result
}