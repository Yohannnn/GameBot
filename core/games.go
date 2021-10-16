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
	Type   string
	State  GameState
	Option Option
}

//Option
//An input for a game update
type Option struct {
	Type     string
	Name     string
	Message  string
	Rollback bool
	Cord     struct {
		x     [2]int
		y     [2]int
		value [2]int
	}
	Select struct {
		Options []string
		value   string
	}
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

//CreateGameUpdate
//Creates a game update to be sent to a player
func CreateGameUpdate(Type string, State GameState, Option Option) GameUpdate {
	//Checks type
	if !Contains([]string{"local", "global", "error", "playerwin", "opponentwin"}, Type) {
		Log.Error("Invalid type for game update")
		return GameUpdate{}
	}

	gU := GameUpdate{
		Type:   Type,
		State:  State,
		Option: Option,
	}
	return gU
}

//gameUpdateLocal
//Sends a local game update
func gameUpdateLocal(info *GameInfo, update GameUpdate, playerID string, opponentID string, currentMessageID string) {
	var stats string
	var board string

	//Gets the dm channel for the player
	playerChannel, err := Session.UserChannelCreate(playerID)
	if err != nil {
		Log.Error(err.Error())
		return
	}

	//Gets the old game message
	message, err := Session.ChannelMessage(playerChannel.ID, currentMessageID)
	if err != nil {
		Log.Error(err.Error())
		return
	}

	//Parses the old message for the embed
	oldEmbed := message.Embeds[0]

	//Creates new embed
	embed := newEmbed()

	//Formats the board
	if update.State.Board == nil {
		board = oldEmbed.Fields[0].Value
	} else {
		board = formatBoard(update.State.Board)
	}
	embed.addField("Board", board, true)

	//Formats the stats
	if update.State.Stats == nil {
		stats = oldEmbed.Fields[1].Value
	} else {
		for stat, value := range update.State.Stats {
			stats += fmt.Sprintf("%s = %s\n", stat, value)
		}
	}
	embed.addField("Game Stats:", stats, true)

	//Adds the selected options field
	for _, field := range oldEmbed.Fields[2:] {
		if field.Name == "Selected Options" {
			embed.addField("Selected Options", field.Value, false)
			break
		}
	}

	//Formats option field
	embed.addField(update.Option.Message, fmt.Sprintf("%s:%s", update.Option.Name, update.Option.Type), false)

	//Adds gameID to footer
	embed.setFooter(oldEmbed.Footer.Text, "", "")

	//Sends the new message and deletes the old one
	embed.send(info.Name, fmt.Sprintf("%s game against <@%s>", info.Name, opponentID), playerChannel.ID)
	err = Session.ChannelMessageDelete(playerChannel.ID, currentMessageID)
	if err != nil {
		Log.Error(err.Error())
		return
	}
}

//gameUpdate
//Sends the game update
func gameUpdate(info *GameInfo, update GameUpdate, playerID string, opponentID string, currentGameID string) {
	var stats string
	var board string

	//Gets the user struct of each player
	player, err := Session.User(playerID)
	if err != nil {
		Log.Error(err.Error())
	}
	opponent, err := Session.User(opponentID)
	if err != nil {
		Log.Error(err.Error())
	}

	//Gets the dm channel for each player
	playerChannel, err := Session.UserChannelCreate(player.ID)
	if err != nil {
		Log.Error(err.Error())
	}
	opponentChannel, err := Session.UserChannelCreate(opponent.ID)
	if err != nil {
		Log.Error(err.Error())
	}

	//Gets the current game message
	message, err := Session.ChannelMessage(playerChannel.ID, currentGameID)
	if err != nil {
		Log.Error(err.Error())
	}

	//Gets the current game embed
	currentEmbed := message.Embeds[0]

	//Creates a new embed
	embed := newEmbed()

	//Formats the board
	if update.State.Board == nil {
		board = currentEmbed.Fields[0].Value
	} else {
		board = formatBoard(update.State.Board)
	}
	embed.addField("Board", board, true)

	//Formats the stats
	if update.State.Stats == nil {
		stats = currentEmbed.Fields[1].Value
	} else {
		for stat, value := range update.State.Stats {
			stats += fmt.Sprintf("%s = %s\n", stat, value)
		}
	}

	//Checks the type of update
	switch update.Type {
	case "playerwin":
		embed.setColor(Yellow)
		err = Session.ChannelMessageDelete(playerChannel.ID, currentGameID)
		if err != nil {
			Log.Error(err.Error())
		}
		embed.send("You Won!", fmt.Sprintf("You won your %s game against <@%s>", info.Name, opponentID), playerChannel.ID)
		embed.setColor(Red)
		embed.send("You Lost!", fmt.Sprintf("You lost your %s game against <@%s>", info.Name, playerID), opponentChannel.ID)
		return

	case "opwin":
		embed.setColor(Red)
		err = Session.ChannelMessageDelete(playerChannel.ID, currentGameID)
		if err != nil {
			Log.Error(err.Error())
		}
		embed.send("You Won!", fmt.Sprintf("You won your %s game against %s", info.Name, playerID), opponent.Username)
		embed.setColor(Red)
		embed.send("You Lost!", fmt.Sprintf("You lost your %s game against %s", info.Name, opponentID), player.Username)
		return

	case "local":
		embed.setColor(info.Color)

		//Adds the selected options field
		for _, field := range currentEmbed.Fields[2:] {
			if field.Name == "Selected Options" {
				embed.addField("Selected Options", field.Value, false)
				break
			}
		}

		//Adds option field and gameID
		embed.addField(update.Option.Message, fmt.Sprintf("%s:%s", update.Option.Name, update.Option.Type), false)
		embed.setFooter(currentEmbed.Footer.Text, "", "")

		//Sends the new message and deletes the old one
		embed.send(info.Name, fmt.Sprintf("%s game against <@%s>", info.Name, opponentID), playerChannel.ID)
		err = Session.ChannelMessageDelete(playerChannel.ID, currentGameID)
		if err != nil {
			Log.Error(err.Error())
			return
		}

	case "err":

	case "global":

	}
}

//formatBoard
//Formats a game board into a string
func formatBoard(board [][]string) string {
	var BoardString string
	var LineString string
	for _, l := range board {
		for _, e := range l {
			emoji, err := Session.State.Emoji("806048328973549578", e)
			if err != nil {
				Log.Error(err.Error())
			}
			LineString += emoji.MessageFormat()
		}
		BoardString += LineString + "\n"
		LineString = ""
	}
	return BoardString
}
