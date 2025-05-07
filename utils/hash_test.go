package utils

import (
	"fmt"
	"testing"
)

func TestHashPwd(t *testing.T) {
	hash := HashPwd("123456")
	fmt.Println(hash)
}
func TestCheckPwd(t *testing.T) {
	ok := CheckPwd("$2a$04$Z7M5TnydUHeCZ6UiQe71qeWoq.a6tOcmtsd1YNLnvpoqFjvNWSNZO", "123456")
	fmt.Println(ok)
}
