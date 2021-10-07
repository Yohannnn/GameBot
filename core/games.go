package core

import (
	"fmt"
	"strings"
)

//Game
//The update function, start function, and game information for a game
type Game struct {
	UpdateFunc func(GameInput) (GameUpdate, string)
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
	PlayerState     GameState
	OpponentSate    GameState
	SelectedOptions []Option
	Option          Option
}

//GameUpdate
//An update to a game
type GameUpdate struct {
	Type            string
	State           GameState
	Option          Option
	SelectedOptions []Option
}

//Games
//Map games names to their game struct
var Games = make(map[string]Game)

//AddGame
//Adds a game to the game map
func AddGame(updateFunc func(GameInput) (GameUpdate, string), startFunc func() GameUpdate, gI *GameInfo) {
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

//gameUpdateLocal
//Sends a local game update
func gameUpdateLocal(info *GameInfo, update GameUpdate, playerID string, opponentID string, currentGameID string) error {
	var stats string
	var board string

	//Gets the dm channel for the player
	playerChannel, err := Session.UserChannelCreate(playerID)
	if err != nil {
		return err
	}

	//Gets the current game message
	message, err := Session.ChannelMessage(playerChannel.ID, currentGameID)
	if err != nil {
		return err
	}

	//Creates new embed
	embed := newEmbed()

	//Formats the board
	if update.State.Board == nil {
		board = message.Embeds[0].Fields[0].Value
	} else {
		for _, l := range update.State.Board {
			var line string
			for _, e := range l {
				line += fmt.Sprintf(":%s:", e)
			}
			board += line + "\n"
			line = ""
		}
	}
	embed.addField("Board", board, true)

	//Formats the stats
	if update.State.Stats == nil {
		stats = message.Embeds[0].Fields[1].Value
	} else {
		for stat, value := range update.State.Stats {
			stats += fmt.Sprintf("%s = %s\n", stat, value)
		}
	}
	embed.addField("Game Stats:", stats, true)

	//Formats option field
	switch update.Option.Type {
	case "select":
		embed.addField("Select an option", update.Option.Name, false)
	}

	return nil
}

//gameUpdatePlayerWin
//Sends the updated game to the opponent
func gameUpdatePlayerWin(info *GameInfo, update GameUpdate, playerID string, opponentID string, currentGameID string) error {
	var stats string
	var board string

	//Gets the user struct of each player
	player, err := Session.User(playerID)
	if err != nil {
		return err
	}
	opponent, err := Session.User(opponentID)
	if err != nil {
		return err
	}

	//Gets the dm channel for each player
	playerChannel, err := Session.UserChannelCreate(player.ID)
	if err != nil {
		return err
	}
	opponentChannel, err := Session.UserChannelCreate(opponent.ID)
	if err != nil {
		return err
	}

	//Gets the current game message
	message, err := Session.ChannelMessage(playerChannel.ID, currentGameID)
	if err != nil {
		return err
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
		embed.setColor(Yellow)
		err = Session.ChannelMessageDelete(playerChannel.ID, currentGameID)
		if err != nil {
			return err
		}
		embed.send("You Won!", fmt.Sprintf("You won your %s game against <@%s>", info.Name, opponentID), playerChannel.ID)
		embed.setColor(Red)
		embed.send("You Lost!", fmt.Sprintf("You lost your %s game against <@%s>", info.Name, playerID), opponentChannel.ID)
		return nil

	case "opponentwin":
		embed.setColor(Red)
		err = Session.ChannelMessageDelete(playerChannel.ID, currentGameID)
		if err != nil {
			return err
		}
		embed.send("You Won!", fmt.Sprintf("You won your %s game against %s", info.Name, playerID), opponent.Username)
		embed.setColor(Red)
		embed.send("You Lost!", fmt.Sprintf("You lost your %s game against %s", info.Name, opponentID), player.Username)
		return nil

	case "local":
		embed.setColor(info.Color)
		err = Session.ChannelMessageDelete(playerChannel.ID, currentGameID)
		if err != nil {
			return err
		}
		embed.setFooter(message.Embeds[0].Footer.Text[6:], "", "")
		embed.send(info.Name, fmt.Sprintf("%s game against %s", info.Name, opponent.Username), playerChannel.ID)
	}
	return nil
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
