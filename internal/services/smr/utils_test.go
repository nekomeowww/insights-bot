package smr

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCheckUrl(t *testing.T) {
	a := assert.New(t)
	err := CheckUrl("")
	a.Equal(err.Error(), ErrNoLink.Error())
	err = CheckUrl("not a url")
	a.Equal(err.Error(), ErrParse.Error())
	err = CheckUrl("://test.com")
	a.Equal(err.Error(), ErrScheme.Error())
	err = CheckUrl("wss://test.com")
	a.Equal(err.Error(), ErrScheme.Error())
	err = CheckUrl("https://test.com")
	a.NoError(err)
}
