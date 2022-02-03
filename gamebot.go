package main

import (
	bot "GameBot/core"
	_ "GameBot/games"
	"math/rand"
	"time"
)

func main() {
	//Seeds rand
	rand.Seed(time.Now().UTC().UnixNano())

	// Start bot
	bot.Start()
}
