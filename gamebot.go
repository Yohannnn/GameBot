package main

import (
	bot "GameBot/core"
)


func main(){
	// Start bot
	err := bot.Start()
	if err != nil{
		bot.Log.Error(err.Error())
	}
}