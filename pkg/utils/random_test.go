package utils

import (
	"fmt"
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRandomBytes(t *testing.T) {
	t.Run("no args", func(t *testing.T) {
		assert := assert.New(t)

		rand1, _ := RandBytes()
		rand2, _ := RandBytes()

		assert.NotEqual(rand1, rand2)
		assert.Equal(len(rand1), 32)
		assert.Equal(len(rand2), 32)
	})
	t.Run("with args", func(t *testing.T) {
		assert := assert.New(t)

		arg := 123
		rand1, _ := RandBytes(arg)
		rand2, _ := RandBytes(arg)

		assert.NotEqual(rand1, rand2)
		assert.Equal(len(rand1), arg)
		assert.Equal(len(rand2), arg)
	})
}

func TestRandomBase64Token(t *testing.T) {
	t.Run("no args", func(t *testing.T) {
		assert := assert.New(t)

		rand1, _ := RandomBase64Token()
		rand2, _ := RandomBase64Token()

		assert.NotEqual(rand1, rand2)
		assert.Equal(len(rand1), 44)
		assert.Equal(len(rand2), 44)
	})
	t.Run("with args", func(t *testing.T) {
		assert := assert.New(t)

		arg := 123
		rand1, _ := RandomBase64Token(arg)
		rand2, _ := RandomBase64Token(arg)

		assert.NotEqual(rand1, rand2)
	})
}

func TestRandomHashString(t *testing.T) {
	assert := assert.New(t)

	hashString := RandomHashString()
	assert.NotEmpty(hashString)
	assert.Len(hashString, 64)

	hashString2 := RandomHashString(32)
	assert.NotEmpty(hashString2)
	assert.Len(hashString2, 32)
}

func TestRandomString(t *testing.T) {
	data, err := RandBytes(32)
	assert.NoError(t, err)
	assert.Len(t, data, 32)

	assert.Len(t, fmt.Sprintf("%x", data), 64)
}

func TestRandomInt64(t *testing.T) {
	assert := assert.New(t)

	assert.NotZero(RandomInt64())
}

func TestRandomInt64InRange(t *testing.T) {
	assert := assert.New(t)

	randomNumber := RandomInt64InRange(100000, 999999)
	assert.NotZero(randomNumber)
	assert.Len(strconv.FormatInt(randomNumber, 10), 6)
}

func TestRandomInt64String(t *testing.T) {
	assert := assert.New(t)

	for i := 0; i < 100; i++ {
		t.Run(fmt.Sprintf("Try#%d", i+1), func(t *testing.T) {
			numberStr := RandomInt64String(6)
			assert.NotEmpty(numberStr)
			assert.Len(numberStr, 6)

			numberStr = RandomInt64String(32)
			assert.NotEmpty(numberStr)
			assert.Len(numberStr, 32)
		})
	}
}
