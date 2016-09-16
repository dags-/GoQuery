package queryutils

import (
	"bytes"
	"encoding/binary"
	"encoding/json"
	"strconv"
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

func ToBytes(input interface{}) []byte {
	buffer := bytes.Buffer{}
	binary.Write(&buffer, binary.BigEndian, input)
	return buffer.Bytes()
}

func ParseInt(input string) int32 {
	num, err := strconv.Atoi(input)
	if err != nil {
		return 0
	}
	return int32(num)
}
