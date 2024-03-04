package i18n

import (
	"os"
	"path/filepath"

	"github.com/nicksnyder/go-i18n/v2/i18n"
	"go.uber.org/zap"
	"golang.org/x/text/language"
	yaml "gopkg.in/yaml.v3"

	"github.com/nekomeowww/insights-bot/pkg/logger"
)

type M = map[string]any

type I18n struct {
	Bundle *i18n.Bundle
	logger *logger.Logger
}

type newI18nOptions struct {
	localesDir string
	logger     *logger.Logger
}

type NewI18nOption func(*newI18nOptions)

func WithLocalesDir(dir string) NewI18nOption {
	return func(o *newI18nOptions) {
		o.localesDir = dir
	}
}

func WithLogger(logger *logger.Logger) NewI18nOption {
	return func(o *newI18nOptions) {
		o.logger = logger
	}
}

func NewI18n(options ...NewI18nOption) (*I18n, error) {
	opts := newI18nOptions{}
	for _, o := range options {
		o(&opts)
	}
	if opts.logger == nil {
		logger, err := logger.NewLogger(zap.DebugLevel, "i18n", "", nil)
		if err != nil {
			return nil, err
		}

		opts.logger = logger
	}

	bundle := i18n.NewBundle(language.English)
	bundle.RegisterUnmarshalFunc("yaml", yaml.Unmarshal)

	dirs, err := os.ReadDir(opts.localesDir)
	if err != nil {
		return nil, err
	}

	for _, v := range dirs {
		if v.IsDir() {
			continue
		}

		localeFilePath := filepath.Join(opts.localesDir, v.Name())

		file, err := bundle.LoadMessageFile(localeFilePath)
		if err != nil {
			return nil, err
		}

		err = bundle.AddMessages(language.SimplifiedChinese, file.Messages...)
		if err != nil {
			return nil, err
		}
	}

	return &I18n{
		Bundle: bundle,
		logger: opts.logger,
	}, nil
}

func (i *I18n) TWithLanguage(lang string, key string, args ...any) string {
	message := i.t(language.Make(lang), key, args...)

	i.logger.Debug("localized message",
		zap.String("lang", lang),
		zap.String("language", language.Make(lang).String()),
		zap.String("key", key),
		zap.Any("args", args),
		zap.String("message", message),
	)

	return message
}

func (i *I18n) TWithTag(lang language.Tag, key string, args ...any) string {
	message := i.t(lang, key, args...)

	i.logger.Debug("localized message",
		zap.String("lang", lang.String()),
		zap.String("key", key),
		zap.Any("args", args),
		zap.String("message", message),
	)

	return message
}

func (i *I18n) t(lang language.Tag, key string, args ...any) string {
	localizer := i18n.NewLocalizer(i.Bundle, lang.String(), language.English.String())

	config := &i18n.LocalizeConfig{
		MessageID: key,
	}
	if len(args) > 0 {
		config.TemplateData = args[0]
	}

	str, err := localizer.Localize(config)
	if err != nil {
		i.logger.Error("failed to localize message",
			zap.String("lang", lang.String()),
			zap.String("key", key),
			zap.Any("args", args),
			zap.Error(err),
		)

		return key
	}

	return str
}
