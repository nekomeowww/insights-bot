package utils

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/davecgh/go-spew/spew"
)

// Print 格式化输出所有传入的值的字段、值、类型、大小。
func Print(inputs ...interface{}) {
	fmt.Println(Sprint(inputs))
}

// Sprint 格式化输出所有传入的值的字段、值、类型、大小，并返回字符串。
//
// NOTICE: 包含换行符。
func Sprint(inputs ...interface{}) string {
	return spew.Sdump(inputs)
}

// PrintJSON 格式化输出 JSON 格式。
func PrintJSON(inputs ...interface{}) {
	fmt.Println(SprintJSON(inputs))
}

// SprintJSON 格式化输出 JSON 格式并返回字符串。
//
// NOTICE: 包含换行符。
func SprintJSON(inputs ...interface{}) string {
	strSlice := make([]string, 0)

	for _, v := range inputs {
		b, _ := json.MarshalIndent(v, "", "  ")
		strSlice = append(strSlice, string(b))
	}

	return strings.Join(strSlice, "\n")
}
