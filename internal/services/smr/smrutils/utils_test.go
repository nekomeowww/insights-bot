package smrutils

import (
	"github.com/nekomeowww/insights-bot/internal/services/smr"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCheckUrl(t *testing.T) {
	a := assert.New(t)
	err := CheckUrl("")
	a.Equal(err.Error(), smr.ErrNoLink.Error())
	err = CheckUrl("not a url")
	a.Equal(err.Error(), smr.ErrParse.Error())
	err = CheckUrl("://test.com")
	a.Equal(err.Error(), smr.ErrScheme.Error())
	err = CheckUrl("wss://test.com")
	a.Equal(err.Error(), smr.ErrScheme.Error())
	err = CheckUrl("https://test.com")
	a.NoError(err)
}
