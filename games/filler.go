package games

import (
	bot "GameBot/core"
	"math/rand"
)

var colors = []string{"Red", "Brown", "Orange", "Yellow", "Green", "Blue", "Purple"}
var BoardKey = map[string]string{
	"Red":    bot.RedSqr,
	"Brown":  bot.BrownSqr,
	"Orange": bot.OrangeSqr,
	"Yellow": bot.YellowSqr,
	"Green":  bot.GreenSqr,
	"Blue":   bot.BlueSqr,
	"Purple": bot.PurpleSqr,
}

func fillerStart(instance *bot.Instance) {
	var reactions []string

	for i := 0; i < 8; i++ {
		var line []string
		//var PlayerColor string
		lastColor := ""
		for i := 0; i < 8; i++ {
			color := colors[rand.Intn(7)]
			if color == lastColor {
				return
			}
			line = append(line, color)
		}
		instance.Board = append(instance.Board, line)
	}

	instance.Stats["PlayerColors"] = []string{instance.Board[7][0], instance.Board[0][7]}
	instance.Stats["DisallowedColors"] = []string{instance.Stats["PlayerColors"][0], instance.Stats["PlayerColors"][1]}

	instance.Board[7][0] = "Player0"
	instance.Board[0][7] = "Player1"

	for _, c := range bot.RemoveItems(colors, instance.Stats["DisallowedColors"]) {
		reactions = append(reactions, BoardKey[c])
	}

	input := bot.CreateInput("Color", "Select a color to switch to", reactions)

	board := instance.Board
	for l, line := range instance.Board {
		for i, c := range line {
			if c == "Player0" {
				board[l][i] = BoardKey[instance.Stats["PlayerColors"][0]]
			} else if c == "Player1" {
				board[l][i] = BoardKey[instance.Stats["PlayerColors"][1]]
			} else {
				board[l][i] = BoardKey[c]
			}
		}
	}

	bot.StartGame(bot.Games["filler"], instance.Players[0], instance.Players[1], board, input)
}

func fillerUpdate(instance *bot.Instance, output bot.Output) {
}

func init() {
	bot.AddGame("Filler",
		"Try to fill the entire board with your color",
		"I'll write rules soon on god",
		[][]string{
			{bot.OrangeSqr, bot.BrownSqr, bot.BlueSqr, bot.YellowSqr, bot.GreenSqr, bot.BlueSqr, bot.BlueSqr, bot.BlueSqr},
			{bot.YellowSqr, bot.YellowSqr, bot.BrownSqr, bot.BrownSqr, bot.PurpleSqr, bot.BlueSqr, bot.BlueSqr, bot.BlueSqr},
			{bot.BlueSqr, bot.GreenSqr, bot.OrangeSqr, bot.YellowSqr, bot.GreenSqr, bot.BlueSqr, bot.BlueSqr, bot.BlueSqr},
			{bot.RedSqr, bot.BrownSqr, bot.RedSqr, bot.BlueSqr, bot.PurpleSqr, bot.BlueSqr, bot.BlueSqr, bot.BlueSqr},
			{bot.RedSqr, bot.RedSqr, bot.OrangeSqr, bot.PurpleSqr, bot.YellowSqr, bot.BlueSqr, bot.PurpleSqr, bot.RedSqr},
			{bot.RedSqr, bot.RedSqr, bot.RedSqr, bot.YellowSqr, bot.GreenSqr, bot.BlueSqr, bot.GreenSqr, bot.BlueSqr},
			{bot.RedSqr, bot.RedSqr, bot.RedSqr, bot.RedSqr, bot.BlueSqr, bot.PurpleSqr, bot.GreenSqr, bot.RedSqr},
			{bot.RedSqr, bot.RedSqr, bot.RedSqr, bot.YellowSqr, bot.RedSqr, bot.BlueSqr, bot.PurpleSqr, bot.GreenSqr},
		},
		fillerStart,
		fillerUpdate,
	)
}
