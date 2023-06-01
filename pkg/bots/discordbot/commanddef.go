package discordbot

import "github.com/disgoorg/disgo/discord"

var commands = []discord.ApplicationCommandCreate{
	discord.SlashCommandCreate{
		Name:        "smr",
		Description: "对网页进行总结",
		Options: []discord.ApplicationCommandOption{
			discord.ApplicationCommandOptionString{
				Required:    true,
				Name:        "link",
				Description: "网页链接",
			},
		},
	},
}
