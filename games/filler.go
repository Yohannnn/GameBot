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
		Guide: "Test Guide",
		Rules: "Test Rules",
		Color: 14495300,
	})
}
