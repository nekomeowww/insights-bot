package smr

import (
	"net/url"

	"github.com/samber/lo"
)

func CheckUrl(urlString string) error {
	if urlString == "" {
		return ErrNoLink
	}

	parsedURL, err2 := url.Parse(urlString)
	if err2 != nil {
		return ErrParse
	}
	if parsedURL.Scheme == "" || !lo.Contains([]string{"http", "https"}, parsedURL.Scheme) {
		return ErrScheme
	}

	return nil
}
