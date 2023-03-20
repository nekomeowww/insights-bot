package utils

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/davecgh/go-spew/spew"
)

// Print 格式化输出所有传入的值的字段、值、类型、大小
func Print(any ...interface{}) {
	fmt.Println(Sprint(any))
}

// Sprint 格式化输出所有传入的值的字段、值、类型、大小，并返回字符串
// NOTICE: 包含换行符
func Sprint(any ...interface{}) string {
	return spew.Sdump(any)
}

// PrintJSON 格式化输出 JSON 格式
func PrintJSON(any ...interface{}) {
	fmt.Println(SprintJSON(any))
}

// SprintJSON 格式化输出 JSON 格式并返回字符串
// NOTICE: 包含换行符
func SprintJSON(any ...interface{}) string {
	strSlice := make([]string, 0)
	for _, v := range any {
		b, _ := json.MarshalIndent(v, "", "  ")
		strSlice = append(strSlice, string(b))
	}
	return strings.Join(strSlice, "\n")
}
