package smrutils

import (
	"errors"
	smr2 "github.com/nekomeowww/insights-bot/internal/models/smr"
	"github.com/nekomeowww/insights-bot/internal/services/smr"
	"net/url"

	"github.com/samber/lo"
)

func CheckUrl(urlString string) (error, error) {
	if urlString == "" {
		return smr.ErrNoLink, nil
	}

	parsedURL, err := url.Parse(urlString)
	if err != nil {
		return smr.ErrParse, err
	}
	if parsedURL.Scheme == "" || !lo.Contains([]string{"http", "https"}, parsedURL.Scheme) {
		return smr.ErrScheme, err
	}

	return nil, nil
}

func FormatUrlCheckError(err error, platform smr2.FromPlatform) string {
	switch {
	case errors.Is(err, smr.ErrNoLink):
		switch platform {
		case smr2.FromPlatformTelegram:
			return "没有找到链接，可以发送一个有效的链接吗？用法：<code>/smr <链接></code>"
		case smr2.FromPlatformDiscord, smr2.FromPlatformSlack:
			return "没有找到链接，可以发送一个有效的链接吗？用法：`/smr <链接>`"
		default:
			return err.Error()
		}
	case errors.Is(err, smr.ErrParse), errors.Is(err, smr.ErrScheme):
		switch platform {
		case smr2.FromPlatformTelegram:
			return "你发来的链接无法被理解，可以重新发一个试试。用法：<code>/smr <链接></code>"
		case smr2.FromPlatformDiscord, smr2.FromPlatformSlack:
			return "你发来的链接无法被理解，可以重新发一个试试。用法：`/smr <链接>`"
		default:
			return err.Error()
		}
	default:
		return err.Error()
	}
}

func IsUrlCheckError(err error) bool {
	return errors.Is(err, smr.ErrNoLink) || errors.Is(err, smr.ErrParse) || errors.Is(err, smr.ErrScheme)
}
