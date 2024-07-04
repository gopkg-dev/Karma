package gob

import (
	"bytes"
	"encoding/gob"
)

// Marshal 使用 gob 编码
// 将给定的值 v 编码为字节切片并返回
func Marshal(v interface{}) ([]byte, error) {
	var (
		buffer bytes.Buffer
	)

	// 使用 gob 编码器将值 v 编码到 buffer 中
	err := gob.NewEncoder(&buffer).Encode(v)
	return buffer.Bytes(), err
}

// Unmarshal 使用 gob 解码
// 将字节切片 data 解码到给定的值 value 中
func Unmarshal(data []byte, value interface{}) error {
	// 使用 gob 解码器将字节切片 data 解码到值 value 中
	err := gob.NewDecoder(bytes.NewReader(data)).Decode(value)
	if err != nil {
		return err
	}
	return nil
}
