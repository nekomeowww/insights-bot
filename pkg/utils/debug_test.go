package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPrintAndPrintJSON(t *testing.T) {
	type testEmbedded struct {
		C int
		D []string
	}
	type testStruct struct {
		A        int
		B        string
		Embedded testEmbedded
	}

	t.Run("Print", func(t *testing.T) {
		assert := assert.New(t)

		assert.NotPanics(func() {
			Print(nil)
		})
		assert.NotPanics(func() {
			Print(testStruct{})
		})
	})

	t.Run("Sprint", func(t *testing.T) {
		assert := assert.New(t)

		assert.NotPanics(func() {
			str := Sprint(nil)
			assert.NotEmpty(str)
		})
		assert.NotPanics(func() {
			str := Sprint(testStruct{})
			assert.NotEmpty(str)
		})
	})

	t.Run("PrintJSON", func(t *testing.T) {
		assert := assert.New(t)

		assert.NotPanics(func() {
			PrintJSON(nil)
		})
		assert.NotPanics(func() {
			PrintJSON(testStruct{})
		})
	})

	t.Run("SprintJSON", func(t *testing.T) {
		assert := assert.New(t)

		assert.NotPanics(func() {
			str := SprintJSON(nil)
			assert.NotEmpty(str)
		})

		assert.NotPanics(func() {
			str := SprintJSON(testStruct{})
			assert.NotEmpty(str)
		})
	})
}
