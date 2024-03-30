package discordbot

import "github.com/disgoorg/disgo/discord"

var commands = []discord.ApplicationCommandCreate{
	discord.SlashCommandCreate{
		Name:        "smr",
		Description: "对网页进行总结",
		DescriptionLocalizations: map[discord.Locale]string{
			discord.LocaleChineseCN: "对网页进行总结",
			discord.LocaleEnglishUS: "Summarize a web article",
			discord.LocaleEnglishGB: "Summarize a web article",
		},
		Options: []discord.ApplicationCommandOption{
			discord.ApplicationCommandOptionString{
				Required: true,
				Name:     "link",
				NameLocalizations: map[discord.Locale]string{
					discord.LocaleChineseCN: "文章链接",
					discord.LocaleEnglishUS: "link",
					discord.LocaleEnglishGB: "link",
				},
				Description: "The link of web article",
				DescriptionLocalizations: map[discord.Locale]string{
					discord.LocaleChineseCN: "需要被总结的文章的链接",
					discord.LocaleEnglishUS: "The link of web article",
					discord.LocaleEnglishGB: "The link of web article",
				},
			},
		},
	},
}
