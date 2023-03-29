package utils

import (
	"os"
	"path/filepath"
	"runtime"
)

// RelativePathOf 获取基于调用函数的调用对象相对位置的相对路径
func RelativePathOf(fp string) string {
	_, file, _, ok := runtime.Caller(1)
	if !ok {
		return ""
	}
	callerDir := filepath.Dir(filepath.FromSlash(file))
	return filepath.FromSlash(filepath.Join(callerDir, fp))
}

func RelativePathBasedOnPwdOf(fp string) string {
	dir, err := os.Getwd()
	if err != nil {
		return ""
	}

	return filepath.FromSlash(filepath.Join(dir, fp))
}
