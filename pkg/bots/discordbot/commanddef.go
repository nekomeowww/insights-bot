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
				Required:    true,
				Name:        "link",
				Description: "Article link",
				DescriptionLocalizations: map[discord.Locale]string{
					discord.LocaleChineseCN: "文章链接",
					discord.LocaleEnglishUS: "Article link",
					discord.LocaleEnglishGB: "Article link",
				},
			},
		},
	},
}
