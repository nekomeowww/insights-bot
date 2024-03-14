package smr

import (
	"errors"
	"net/url"

	"github.com/nekomeowww/insights-bot/pkg/i18n"
	"github.com/nekomeowww/insights-bot/pkg/types/bot"
	"github.com/samber/lo"
)

func CheckUrl(urlString string) (error, error) {
	if urlString == "" {
		return ErrNoLink, nil
	}

	parsedURL, err := url.Parse(urlString)
	if err != nil {
		return ErrParse, err
	}
	if parsedURL.Scheme == "" || !lo.Contains([]string{"http", "https"}, parsedURL.Scheme) {
		return ErrScheme, err
	}

	return nil, nil
}

func FormatUrlCheckError(err error, platform bot.FromPlatform, language string, i18n *i18n.I18n) string {
	switch {
	case errors.Is(err, ErrNoLink):
		switch platform {
		case bot.FromPlatformTelegram:
			return i18n.TWithLanguage(language, "commands.groups.summarization.commands.smr.noLinksFound")
			// TODO: support i18n for discord and slack
		case bot.FromPlatformDiscord, bot.FromPlatformSlack:
			return "没有找到链接，可以发送一个有效的链接吗？用法：`/smr <链接>`"
		default:
			return err.Error()
		}
	case errors.Is(err, ErrParse), errors.Is(err, ErrScheme):
		switch platform {
		case bot.FromPlatformTelegram:
			return i18n.TWithLanguage(language, "commands.groups.summarization.commands.smr.invalidLink")
		case bot.FromPlatformDiscord, bot.FromPlatformSlack:
			// TODO: support i18n for discord and slack
			return "你发来的链接无法被理解，可以重新发一个试试。用法：`/smr <链接>`"
		default:
			return err.Error()
		}
	default:
		return err.Error()
	}
}

func IsUrlCheckError(err error) bool {
	return errors.Is(err, ErrNoLink) || errors.Is(err, ErrParse) || errors.Is(err, ErrScheme)
}
