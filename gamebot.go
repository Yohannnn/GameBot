package main

import (
	bot "GameBot/core"
	_ "GameBot/games"
	"encoding/json"
	"io/ioutil"
	"log"
	"math/rand"
	"os"
	"strings"
	"time"
)

func main() {
	//Seeds rand
	rand.Seed(time.Now().UTC().UnixNano())

	//Unmarshalls instances JSON
	file, err := os.Open("instances.json")
	if err != nil {
		log.Fatal(err.Error())
	}
	defer file.Close()

	byteval, err := ioutil.ReadAll(file)
	if err != nil {
		log.Fatal(err.Error())
	}

	JSONInstances := make(map[string]bot.JSONInstance)
	json.Unmarshal(byteval, &JSONInstances)

	for ID, inst := range JSONInstances {
		bot.Instances[ID] = &bot.Instance{
			ID:               ID,
			Game:             bot.Games[strings.ToLower(inst.GameName)],
			Board:            inst.Board,
			DisplayBoard:     inst.DisplayBoard,
			CurrentInput:     inst.CurrentInput,
			Stats:            inst.Stats,
			CurrentMessageID: inst.CurrentMessageID,
			Players:          inst.Players,
			Turn:             inst.Turn,
		}
	}

	// Start bot
	bot.Start()
}
