package games

import (
	bot "GameBot/core"
)

const red = "red_square"
const brown = "brown_square"
const orange = "orange_square"
const yellow = "yellow_square"
const green = "green_square"
const blue = "blue_square"
const purple = "purple_square"
const black = "black_large_square"
const white = "white_large_square"

var fillerInfo = bot.CreateGameInfo(
	"Filler",
	"Try to fill the entire board with your color",
	"I'll write rules soon on god",
	0,
	[][]string{
		{red, brown, orange, yellow, green, blue, purple, black},
		{red, brown, orange, yellow, green, blue, purple, white},
		{red, brown, orange, yellow, green, blue, purple, white},
		{red, brown, orange, yellow, green, blue, purple, white},
		{red, brown, orange, yellow, green, blue, purple, white},
		{red, brown, orange, yellow, green, blue, purple, white},
		{red, brown, orange, yellow, green, blue, purple, white},
		{red, brown, orange, yellow, green, blue, purple, white},
	},
)

func fillerStart() *bot.GameUpdate {
	return nil
}

func init() {
	bot.AddGame(fillerStart, nil, fillerInfo)
}
