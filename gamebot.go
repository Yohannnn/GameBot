package main

import (
	bot "GameBot/core"
	_ "GameBot/games"
)

func main() {
	// Start bot
	err := bot.Start()
	if err != nil {
		bot.Log.Error(err.Error())
	}
}
