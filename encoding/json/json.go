package json

import (
	"bytes"

	jsoniter "github.com/json-iterator/go"
	"github.com/tidwall/gjson"
)

// 定义JSON操作
var (
	json          = jsoniter.ConfigCompatibleWithStandardLibrary
	Marshal       = json.Marshal
	Unmarshal     = json.Unmarshal
	MarshalIndent = json.MarshalIndent
	NewDecoder    = json.NewDecoder
	NewEncoder    = json.NewEncoder
)

// Valid reports whether data is a valid JSON encoding.
func Valid(data []byte) bool {
	return json.Valid(data)
}

// ValidString reports whether data is a valid JSON encoding.
func ValidString(str string) bool {
	return gjson.Valid(str)
}

// GetStringFromJson get the string value from json path
func GetStringFromJson(json, path string) string {
	return gjson.Get(json, path).String()
}

// UnmarshalString 解码字符串为 JSON
func UnmarshalString(data string, v interface{}) error {
	return json.Unmarshal([]byte(data), v)
}

// MarshalToString JSON编码为字符串
func MarshalToString(v interface{}) string {
	s, err := json.MarshalToString(v)
	if err != nil {
		return ""
	}
	return s
}

// MarshalToBytes JSON编码为字节数组
func MarshalToBytes(v interface{}) []byte {
	s, err := json.Marshal(v)
	if err != nil {
		return []byte{}
	}
	return s
}

// MarshalIndentToString JSON编码为格式化字符串
func MarshalIndentToString(v interface{}) string {
	bf := bytes.NewBuffer([]byte{})
	encoder := NewEncoder(bf)
	encoder.SetEscapeHTML(false)
	encoder.SetIndent("", "\t")
	_ = encoder.Encode(v)
	return bf.String()
}
