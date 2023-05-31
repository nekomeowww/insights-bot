package smrutils

import (
	"errors"
	"github.com/nekomeowww/insights-bot/internal/services/smr"
	"net/url"

	"github.com/samber/lo"
)

func CheckUrl(urlString string) error {
	if urlString == "" {
		return smr.ErrNoLink
	}

	parsedURL, err2 := url.Parse(urlString)
	if err2 != nil {
		return smr.ErrParse
	}
	if parsedURL.Scheme == "" || !lo.Contains([]string{"http", "https"}, parsedURL.Scheme) {
		return smr.ErrScheme
	}

	return nil
}

func IsUrlCheckError(err error) bool {
	return errors.Is(err, smr.ErrNoLink) || errors.Is(err, smr.ErrParse) || errors.Is(err, smr.ErrScheme)
}
