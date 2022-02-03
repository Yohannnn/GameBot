package games

import (
	bot "GameBot/core"
	"math/rand"
)

var colors = []string{"red", "Brown", "Orange", "Yellow", "Green", "Blue", "Purple", "White"}

func fillerStart(instance *bot.Instance) {
	for i := 0; i < 8; i++ {
		for i := 0; i < 8; i++ {
			lastColor := ""
			var line []string
			color := colors[rand.Intn(7)]
			if color == lastColor {
				return
			}
			line = append(line, color)
		}
		line := []string{colors[rand.Intn(7)], colors[rand.Intn(7)], colors[rand.Intn(7)], colors[rand.Intn(7)], colors[rand.Intn(7)], colors[rand.Intn(7)], colors[rand.Intn(7)], colors[rand.Intn(7)]}
		instance.Board = append(instance.Board, line)
	}
	instance.Stats["Player0Color"] = instance.Board[7][0]
	instance.Stats["Player1Color"] = instance.Board[0][7]

	instance.Board[7][0] = "Player0"
	instance.Board[0][7] = "Player1"

	bot.StartGame()

}

func fillerUpdate(instance *bot.Instance) {
}

func init() {
	bot.AddGame("Filler",
		"Try to fill the entire board with your color",
		"I'll write rules soon on god",
		[][]string{
			{bot.OrangeSqr, bot.BrownSqr, bot.BlueSqr, bot.YellowSqr, bot.GreenSqr, bot.BlueSqr, bot.BlueSqr, bot.BlueSqr},
			{bot.YellowSqr, bot.YellowSqr, bot.WhiteSqr, bot.BrownSqr, bot.PurpleSqr, bot.BlueSqr, bot.BlueSqr, bot.BlueSqr},
			{bot.BlueSqr, bot.GreenSqr, bot.OrangeSqr, bot.YellowSqr, bot.GreenSqr, bot.BlueSqr, bot.BlueSqr, bot.BlueSqr},
			{bot.RedSqr, bot.BrownSqr, bot.WhiteSqr, bot.BlueSqr, bot.PurpleSqr, bot.BlueSqr, bot.BlueSqr, bot.BlueSqr},
			{bot.RedSqr, bot.RedSqr, bot.OrangeSqr, bot.PurpleSqr, bot.WhiteSqr, bot.BlueSqr, bot.PurpleSqr, bot.RedSqr},
			{bot.RedSqr, bot.RedSqr, bot.RedSqr, bot.YellowSqr, bot.GreenSqr, bot.BlueSqr, bot.PurpleSqr, bot.WhiteSqr},
			{bot.RedSqr, bot.RedSqr, bot.RedSqr, bot.RedSqr, bot.BlueSqr, bot.WhiteSqr, bot.WhiteSqr, bot.RedSqr},
			{bot.RedSqr, bot.RedSqr, bot.RedSqr, bot.YellowSqr, bot.RedSqr, bot.BlueSqr, bot.PurpleSqr, bot.GreenSqr},
		},
		fillerStart,
		fillerUpdate,
	)
}
