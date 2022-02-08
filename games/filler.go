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

	//Generates a game board
	for l := 0; l < 8; l++ {
		var line []string
		for i := 0; i < 8; i++ {
			possibleColors := colors

			//Removes adjacent colors
			if l > 0 {
				possibleColors = bot.RemoveItems(possibleColors, []string{instance.Board[l-1][i]})
			}
			if i > 0 {
				possibleColors = bot.RemoveItems(possibleColors, []string{line[i-1]})
			}

			//Removes player1 color if its player0 color select
			if l == 7 && i == 0 {
				possibleColors = bot.RemoveItems(possibleColors, []string{instance.Board[0][7]})
			}

			//Picks a random color
			color := possibleColors[rand.Intn(len(possibleColors))]
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

	instance.DisplayBoard = instance.Board
	for l, line := range instance.Board {
		for i, c := range line {
			if c == "Player0" {
				instance.DisplayBoard[l][i] = BoardKey[instance.Stats["PlayerColors"][0]]
			} else if c == "Player1" {
				instance.DisplayBoard[l][i] = BoardKey[instance.Stats["PlayerColors"][1]]
			} else {
				instance.DisplayBoard[l][i] = BoardKey[c]
			}
		}
	}

	bot.StartGame(instance, instance.DisplayBoard, input)
}

func fillerUpdate(instance *bot.Instance, output bot.Output) {
	//Errors
	if len(output.SelOptions) > 1 {
		bot.EditGame(instance, instance.DisplayBoard, bot.CreateInput("Color", "You can only select 1 color", instance.CurrentInput.Options))
		return
	}

	//Checks for win
	for i, l := range instance.Board {
		var Count int
		for _, c := range l {
			if c == "Player0" {
				Count++
			} else if c != "Player1" {
				break
			}
		}
		if i == 7 {
			if Count > 64 {
				bot.EndGame(instance, instance.Players[0], instance.Players[1], instance.DisplayBoard)
			} else if Count < 64 {
				bot.EndGame(instance, instance.Players[1], instance.Players[0], instance.DisplayBoard)
			}
		}
	}
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
