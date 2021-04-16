package sco

import (
	"errors"
	"fmt"
)

// 构造错误字符串
func MakeErrorMsg(l1, l2 string) string {
	return fmt.Sprintf("%s:%s", l1, l2)
}

// 返回一个错误
func MakeError(l1, l2 string) error {
	return errors.New(MakeErrorMsg(l1, l2))
}
