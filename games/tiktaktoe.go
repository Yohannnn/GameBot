package games

import (
	bot "GameBot/core"
)

func tiktaktoeStart(instance *bot.Instance) {
}

func tiktaktoeUpdate(instance *bot.Instance, output bot.Output) {
}

func init() {
	bot.AddGame("TikTakToe",
		"Try to get 3 in a row",
		"Each player gets turns placing mark on the board\nFirst to 3 in a row wins",
		[][]string{
			{bot.BlackSqr, bot.LetCord[1], bot.LetCord[2], bot.LetCord[3]},
			{bot.NumCord[1], bot.WhiteSqr, bot.WhiteSqr, bot.WhiteSqr},
			{bot.NumCord[2], bot.WhiteSqr, bot.WhiteSqr, bot.WhiteSqr},
			{bot.NumCord[3], bot.WhiteSqr, bot.WhiteSqr, bot.WhiteSqr},
		},
		tiktaktoeStart,
		tiktaktoeUpdate,
	)
}
