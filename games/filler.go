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

func fillerStart() (*bot.Instance, bot.GameUpdate) {
	return &bot.Instance{}, bot.GameUpdate{}
}

func fillerUpdate(instance *bot.Instance) bot.GameUpdate {
	return bot.GameUpdate{}
}

func init() {
	bot.AddGame("Filler",
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
		fillerStart,
		fillerUpdate,
	)
}
