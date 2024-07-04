package gob_test

import (
	"testing"

	"github.com/gopkg-dev/karma/encoding/gob"
	"github.com/stretchr/testify/assert"
)

// 定义一个测试结构体
type User struct {
	Name string
	Age  int
}

func TestGobEncoding(t *testing.T) {
	// 创建一个 User 实例
	user := User{Name: "Alice", Age: 30}

	// 对 user 进行 gob 编码
	data, err := gob.Marshal(user)
	assert.NoError(t, err, "编码时应该没有错误")

	// 创建一个空的 User 实例用于解码
	var decodedUser User

	// 对数据进行 gob 解码
	err = gob.Unmarshal(data, &decodedUser)
	assert.NoError(t, err, "解码时应该没有错误")

	// 验证解码后的数据是否与原始数据一致
	assert.Equal(t, user, decodedUser, "解码后的数据应该与原始数据一致")
}
