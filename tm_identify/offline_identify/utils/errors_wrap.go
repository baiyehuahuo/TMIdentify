package utils

import (
	"runtime"

	"github.com/pkg/errors"
)

// ErrWrapOrWithMessage 给错误附加当前运行的函数名
// wrap -> true 表示err是来自非本人编写的函数返回的
// wrap -> false 表示err是来自本人编写的函数返回的
func ErrWrapOrWithMessage(wrap bool, err error) error {
	pc := make([]uintptr, 1)
	runtime.Callers(2, pc)
	f := runtime.FuncForPC(pc[0])
	message := f.Name() + " fail"
	if wrap {
		return errors.Wrap(err, message)
	}
	return errors.WithMessage(err, message)
}
