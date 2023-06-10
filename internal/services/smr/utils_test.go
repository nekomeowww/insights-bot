package smr

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCheckUrl(t *testing.T) {
	a := assert.New(t)
	err, _ := CheckUrl("")
	a.ErrorIs(err, ErrNoLink)
	err, _ = CheckUrl("not a url")
	a.ErrorIs(err, ErrScheme)
	err, _ = CheckUrl("://test.com")
	a.ErrorIs(err, ErrParse)
	err, _ = CheckUrl("wss://test.com")
	a.ErrorIs(err, ErrScheme)
	err, _ = CheckUrl("https://test.com")
	a.NoError(err)
}
