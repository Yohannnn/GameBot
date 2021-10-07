package core

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
	"strings"
)

//Game
//The update function, start function, and game information for a game
type Game struct {
	UpdateFunc func(GameInput) GameUpdate
	StartFunc  func() GameUpdate
	Info       *GameInfo
}

//GameState
//The state of a game
type GameState struct {
	Board [][]string
	Stats map[string]string
}

//GameInfo
//Information about a game
type GameInfo struct {
	Name         string
	Description  string
	Rules        string
	Color        int
	ExampleBoard [][]string
}

//GameInput
//Information about a games state and the players input
type GameInput struct {
	PlayerState  GameState
	OpponentSate GameState
	OptionType   string
	Option       []string
}

//GameUpdate
//An update to a game
type GameUpdate struct {
	Type   string
	State  GameState
	Option Option
}

//Option
//An option for a game update
type Option struct {
	Type      string
	Name      string
	Reactions []string
}

//Games
//Map games names to their game struct
var Games = make(map[string]Game)

//AddGame
//Adds a game to the game map
func AddGame(updateFunc func(GameInput) GameUpdate, startFunc func() GameUpdate, gI *GameInfo) {
	game := Game{
		UpdateFunc: updateFunc,
		StartFunc:  startFunc,
		Info:       gI,
	}
	Games[strings.ToLower(gI.Name)] = game
}

//CreateGameInfo
//Creates game info for a game
func CreateGameInfo(name string, description string, rules string, color int, exampleBoard [][]string) *GameInfo {
	gI := &GameInfo{
		Name:         name,
		Description:  description,
		Rules:        rules,
		ExampleBoard: exampleBoard,
		Color:        color,
	}
	return gI
}

//sendGameUpdate
//Sends the updated game to the opponent
func sendGameUpdate(info *GameInfo, update GameUpdate, playerID string, opponentID string) {
	var stats string
	var board string

	//Gets the dm channel for the opponent
	opponentChannel, err := Session.UserChannelCreate(opponentID)
	if err != nil {
		Log.Error(err.Error())
		return
	}

	//Gets the dm channel for the player
	playerChannel, err := Session.UserChannelCreate(playerID)
	if err != nil {
		Log.Error(err.Error())
		return
	}

	//Creates a new embed
	embed := newEmbed()

	//Formats the game stats
	for stat, value := range update.State.Stats {
		stats += fmt.Sprintf("%s = %s\n", stat, value)
	}
	embed.addField("Game Stats:", stats, true)

	//Formats the game board
	for _, l := range update.State.Board {
		var line string
		for _, e := range l {
			line += fmt.Sprintf(":%s:", e)
		}
		board += line + "\n"
		line = ""
	}
	embed.addField("Board", board, true)

	//Checks the type of update
	switch update.Type {
	case "playerwin":
		embed.setTitle("You Won!")
		embed.setDescription(fmt.Sprintf("You won your %s game against <@%s>", info.Name, playerID))
		embed.setColor(Yellow)
	}

	//Adds option field
	switch update.Option.Type {
	case "select":
		embed.addField("Input", "Select an option", false)
	case "Coordinate":
		embed.addField("Input", "Select a coordinate", false)
	}

	//Sends message
	m := embed.send(channel.ID, info.Name, fmt.Sprintf("Playing against <@%s>", playerID))

	//Adds the reactions
	for _, e := range update.Options {
		err = Session.MessageReactionAdd(channel.ID, m.ID, e)
		if err != nil {
			Log.Error(err.Error())
			return
		}
	}
	err = Session.MessageReactionAdd(channel.ID, m.ID, "✅")
	if err != nil {
		Log.Error(err.Error())
		return
	}
}

//formatBoard
//Formats a game board into a string
func formatBoard(board [][]string) string {
	var BoardString string
	var LineString string
	for _, l := range board {
		for _, e := range l {
			LineString += fmt.Sprintf(":%s:", e)
		}
		BoardString += LineString + "\n"
		LineString = ""
	}
	return BoardString
}
