package smrutils

import (
	"github.com/nekomeowww/insights-bot/internal/services/smr"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCheckUrl(t *testing.T) {
	a := assert.New(t)
	err, _ := CheckUrl("")
	a.ErrorIs(err, smr.ErrNoLink)
	err, _ = CheckUrl("not a url")
	a.ErrorIs(err, smr.ErrScheme)
	err, _ = CheckUrl("://test.com")
	a.ErrorIs(err, smr.ErrParse)
	err, _ = CheckUrl("wss://test.com")
	a.ErrorIs(err, smr.ErrScheme)
	err, _ = CheckUrl("https://test.com")
	a.NoError(err)
}
