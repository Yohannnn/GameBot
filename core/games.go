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
	Type        string
	PlayerState GameState
	OpState     GameState
	Option      Option
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
func CreateGameUpdate(Type string, PlayerState GameState, OpState GameState, Option Option) GameUpdate {
	//Checks type
	if !Contains([]string{"local", "global", "error", "playerwin", "opponentwin"}, Type) {
		Log.Error("Invalid type for game update")
		return GameUpdate{}
	}

	gU := GameUpdate{
		Type:        Type,
		PlayerState: PlayerState,
		OpState:     OpState,
		Option:      Option,
	}
	return gU
}

//gameUpdate
//Sends the game update
func gameUpdate(info *GameInfo, update GameUpdate, playerID string, opponentID string, currentGameID string) {
	var stats string
	var board string
	var opEmbed embed
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
	opChannel, err := Session.UserChannelCreate(opponent.ID)
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
	playerEmbed := newEmbed()

	//Formats the board
	if update.PlayerState.Board == nil {
		board = currentEmbed.Fields[0].Value
	} else {
		board = formatBoard(update.PlayerState.Board)
	}
	playerEmbed.addField("Board", board, true)

	//Formats the stats
	if update.PlayerState.Stats == nil {
		stats = currentEmbed.Fields[1].Value
	} else {
		for stat, value := range update.PlayerState.Stats {
			stats += fmt.Sprintf("%s = %s\n", stat, value)
		}
		playerEmbed.addField("Stats", stats, true)
	}

	if update.Type != "local" || update.Type != "global" {
		//Creates a new embed
		opEmbed := newEmbed()

		//Formats the board
		if update.OpState.Board == nil {
			board = currentEmbed.Fields[0].Value
		} else {
			board = formatBoard(update.OpState.Board)
		}
		opEmbed.addField("Board", board, true)

		//Formats the stats
		if update.OpState.Stats == nil {
			stats = currentEmbed.Fields[1].Value
		} else {
			for stat, value := range update.OpState.Stats {
				stats += fmt.Sprintf("%s = %s\n", stat, value)
			}
		}
		opEmbed.addField("Stats", stats, true)
	}

	//Checks the type of update
	switch update.Type {

	case "local":
		//Adds the selected options field
		for _, field := range currentEmbed.Fields[2:] {
			if field.Name == "Selected Options" {
				playerEmbed.addField("Selected Options", field.Value, false)
				break
			}
		}

		//Adds option field and gameID
		playerEmbed.addField(update.Option.Message, update.Option.Name, false)
		playerEmbed.setFooter(currentEmbed.Footer.Text, "", "")

		//Sends the new message and deletes the old one
		newMessage := playerEmbed.send(info.Name, fmt.Sprintf("%s game against %s", info.Name, opponent.Username), playerChannel.ID)
		err = Session.ChannelMessageDelete(playerChannel.ID, currentGameID)
		if err != nil {
			Log.Error(err.Error())
			return
		}

		//Adds the reactions to the message
		addOption(update.Option, playerChannel.ID, newMessage.ID)

	case "global":
		playerEmbed.setColor(info.Color)
		opEmbed.setColor(info.Color)
		//Adds the selected options field
		for _, field := range currentEmbed.Fields[2:] {
			if field.Name == "Selected Options" {
				playerEmbed.addField("Selected Options", field.Value, false)
				break
			}
		}

		//Sends the new message
		newMessage := playerEmbed.send(info.Name, fmt.Sprintf("%s game against %s", info.Name, player.Username), opChannel.ID)

		//Adds option field and gameID
		playerEmbed.addField(update.Option.Message, update.Option.Name, false)
		playerEmbed.setFooter(currentGameID, "", "")

		//Adds the reactions to the message
		addOption(update.Option, opChannel.ID, newMessage.ID)

	case "playerwin":
		playerEmbed.setColor(Yellow)
		err = Session.ChannelMessageDelete(playerChannel.ID, currentGameID)
		if err != nil {
			Log.Error(err.Error())
		}
		playerEmbed.send("You Won!", fmt.Sprintf("You won your %s game against <@%s>", info.Name, opponentID), playerChannel.ID)
		playerEmbed.setColor(Red)
		playerEmbed.send("You Lost!", fmt.Sprintf("You lost your %s game against <@%s>", info.Name, playerID), opChannel.ID)
		return

	case "opwin":
		playerEmbed.setColor(Red)
		err = Session.ChannelMessageDelete(playerChannel.ID, currentGameID)
		if err != nil {
			Log.Error(err.Error())
		}
		playerEmbed.send("You Won!", fmt.Sprintf("You won your %s game against %s", info.Name, playerID), opponent.Username)
		playerEmbed.setColor(Red)
		playerEmbed.send("You Lost!", fmt.Sprintf("You lost your %s game against %s", info.Name, opponentID), player.Username)
		return

	case "err":
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
				return ""
			}
			LineString += emoji.MessageFormat()
		}
		BoardString += LineString + "\n"
		LineString = ""
	}
	return BoardString
}
