package hash

import (
	"testing"
)

func TestGeneratePassword(t *testing.T) {
	origin := "123456"
	hashPwd, err := GenerateFromPassword(origin)
	if err != nil {
		t.Error("GeneratePassword Failed: ", err.Error())
	}
	t.Log("test password: ", hashPwd, ",length: ", len(hashPwd))

	if err := CompareHashAndPassword(hashPwd, origin); err != nil {
		t.Error("Unmatched password: ", err.Error())
	}
}

func TestMD5(t *testing.T) {
	origin := "123456"
	hashVal := "e10adc3949ba59abbe56e057f20f883e"
	if v := MD5String(origin); v != hashVal {
		t.Error("Failed to generate MD5 hash: ", v)
	}
}
