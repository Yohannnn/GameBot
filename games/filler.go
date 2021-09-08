package games

import (
	bot "GameBot/core"
	"fmt"
)

func filler(game bot.GameInfo) {
	fmt.Println(game)
}

func init() {
	bot.AddGame("filler", bot.Game{
		Name:     "Filler",
		Function: filler,
		PreviewBoard: [][]string{
			{""},
		},
		Description: "Test Description",
		Rules:       "Test Rules",
		Color:       14495300,
	})
}
