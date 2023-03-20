package options

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestApplyCallOptions(t *testing.T) {
	assert := assert.New(t)

	type optionsA struct {
		a string
		b string
	}

	callOpts := NewCallOptions(func(o *optionsA) {
		o.a = "a"
	})

	opts := ApplyCallOptions([]CallOptions[optionsA]{callOpts})
	assert.Equal("a", opts.a)
	assert.Empty(opts.b)

	callOpts2 := NewCallOptions(func(o *optionsA) {
		o.b = "b"
	})

	opts2 := ApplyCallOptions([]CallOptions[optionsA]{callOpts2})
	assert.Empty(opts2.a)
	assert.Equal("b", opts2.b)
}
