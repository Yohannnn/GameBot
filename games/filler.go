package games

import (
	bot "GameBot/core"
	"math/rand"
)

var colors = []string{"Red", "Brown", "Orange", "Yellow", "Green", "Blue", "Purple"}
var boardKey = map[string]string{
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

	instance.Board[7][0] = "0"
	instance.Board[0][7] = "1"

	for _, c := range bot.RemoveItems(colors, instance.Stats["DisallowedColors"]) {
		reactions = append(reactions, boardKey[c])
	}

	input := bot.CreateInput("Color", "Select a color to switch to", reactions)

	instance.DisplayBoard = [][]string{
		{"", "", "", "", "", "", "", ""},
		{"", "", "", "", "", "", "", ""},
		{"", "", "", "", "", "", "", ""},
		{"", "", "", "", "", "", "", ""},
		{"", "", "", "", "", "", "", ""},
		{"", "", "", "", "", "", "", ""},
		{"", "", "", "", "", "", "", ""},
		{"", "", "", "", "", "", "", ""},
	}

	for l, line := range instance.Board {
		for i, c := range line {
			if c == "0" {
				instance.DisplayBoard[l][i] = boardKey[instance.Stats["PlayerColors"][0]]
			} else if c == "1" {
				instance.DisplayBoard[l][i] = boardKey[instance.Stats["PlayerColors"][1]]
			} else {
				instance.DisplayBoard[l][i] = boardKey[c]
			}
		}
	}

	bot.StartGame(instance, input)
}

func fillerUpdate(instance *bot.Instance, output bot.Output) {
	var searchVal string
	var color string
	var reactions []string

	//Errors
	if len(output.SelOptions) > 1 {
		bot.EditGame(instance, bot.CreateInput("Color", "You can only select 1 color", instance.CurrentInput.Options))
		return
	}

	//Gets color value of output
	switch output.SelOptions[0] {
	case "🟥":
		color = "Red"
	case "🟫":
		color = "Brown"
	case "🟧":
		color = "Orange"
	case "🟨":
		color = "Yellow"
	case "🟩":
		color = "Green"
	case "🟦":
		color = "Blue"
	case "🟪":
		color = "Purple"
	}

	//Gets player value to search for
	if instance.Turn == 0 {
		searchVal = "0"
	} else if instance.Turn == 1 {
		searchVal = "1"
	}

	//Converts adjacent colors to player colors
	for l, line := range instance.Board {
		for i, c := range line {
			if c == searchVal {
				//Checks horiz colors
				if i > 0 {
					if instance.Board[l][i-1] == color {
						instance.Board[l][i-1] = searchVal
					}
				}
				if i < 7 {
					if instance.Board[l][i+1] == color {
						instance.Board[l][i+1] = searchVal
					}
				}
				//Checks vert colors
				if l > 0 {
					if instance.Board[l-1][i] == color {
						instance.Board[l-1][i] = searchVal
					}
				}
				if l < 7 {
					if instance.Board[l+1][i] == color {
						instance.Board[l+1][i] = searchVal
					}
				}
			}
		}
	}

	//Defines new game stats
	instance.Stats["PlayerColors"][instance.Turn] = color
	instance.Stats["DisallowedColors"] = []string{instance.Stats["PlayerColors"][0], instance.Stats["PlayerColors"][1]}

	//Renders new display board
	for l, line := range instance.Board {
		for i, c := range line {
			if c == "0" {
				instance.DisplayBoard[l][i] = boardKey[instance.Stats["PlayerColors"][0]]
			} else if c == "1" {
				instance.DisplayBoard[l][i] = boardKey[instance.Stats["PlayerColors"][1]]
			} else {
				instance.DisplayBoard[l][i] = boardKey[c]
			}
		}
	}

	//Checks for win
out:
	for i, l := range instance.Board {
		var Count int
		for _, c := range l {
			if c == "0" {
				Count++
			} else if c != "1" {
				break out
			}
		}
		if i == 7 {
			if Count > 64 {
				bot.EndGame(instance, instance.Players[0], instance.Players[1])
				return
			} else if Count < 64 {
				bot.EndGame(instance, instance.Players[1], instance.Players[0])
				return
			}
		}
	}

	for _, c := range bot.RemoveItems(colors, instance.Stats["DisallowedColors"]) {
		reactions = append(reactions, boardKey[c])
	}

	input := bot.CreateInput("Color", "Select a color to switch to", reactions)

	//Sends update
	bot.UpdateGame(instance, input)

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
