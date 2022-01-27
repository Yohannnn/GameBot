package core

import (
	"fmt"
	"strings"
)

//Game
//Information about a game
type Game struct {
	Name         string
	Description  string
	Rules        string
	Color        int
	ExampleBoard [][]string
	StartFunc    func() (*GameInstance, GameUpdate)
	UpdateFunc   func(*GameInstance) GameUpdate
}

//GameInstance
//An instance of a game
type GameInstance struct {
	ID               string
	Game             Game
	OneBoard         [][]string
	TwoBoard         [][]string
	Stats            map[string]string
	Options          []string
	OneTurn          bool
	CurrentMessageID string
	OneName          string
	TwoName          string
	OneChannelID     string
	TwoChannelID     string
}

//GameUpdate
//An update to a game
type GameUpdate struct {
	Type     string
	OneBoard [][]string
	TwoBoard [][]string
	Input    Input
}

//Games
//Map games names to their game struct
var Games = make(map[string]Game)

//Instances
//Map of IDs to their game instance
var Instances = make(map[string]*GameInstance)

//AddGame
//Adds a game to the game map
func AddGame(Name string, Description string, Rules string, Color int, ExampleBoard [][]string, StartFunc func() (*GameInstance, GameUpdate), UpdateFunc func(*GameInstance) GameUpdate) {
	game := Game{
		Name:         Name,
		Description:  Description,
		Rules:        Rules,
		Color:        Color,
		ExampleBoard: ExampleBoard,
		StartFunc:    StartFunc,
		UpdateFunc:   UpdateFunc,
	}
	Games[strings.ToLower(Name)] = game
}

//CreateGameUpdate
//Creates a game update to be sent to a player
func CreateGameUpdate(Type string, POneBoard [][]string, PTwoBoard [][]string, Input Input) GameUpdate {
	//Checks type
	if !Contains([]string{"local", "global", "error", "playerwin", "opponentwin"}, Type) {
		Log.Error("Invalid type for game update")
		return GameUpdate{}
	}

	gU := GameUpdate{
		Type:     Type,
		OneBoard: POneBoard,
		TwoBoard: PTwoBoard,
		Input:    Input,
	}

	return gU
}

//gameUpdate
//Sends the game update
func gameUpdate(instance *GameInstance, update GameUpdate) {
	var PEmbed *embed
	var OpEmbed *embed

	//Creates a new embed
	playerEmbed := newEmbed()

	//Formats the player game board
	PBoard := formatBoard(update.PBoard)
	PEmbed.addField("Board", PBoard, true)

	if update.Type != "local" && update.Type != "global" {
		//Creates a new embed
		OpEmbed = newEmbed()

		//Formats the board
		OpBoard := formatBoard(update.OpBoard)
		OpEmbed.addField("Board", OpBoard, true)
	}

	//Checks the type of update
	switch update.Type {

	case "local":
		//Adds option field and gameID
		playerEmbed.addField(update.Input.Message, update.Input.Name, false)
		playerEmbed.setFooter(instance.ID, "", "")

		//Sends the new message and deletes the old one
		newMessage := playerEmbed.send(instance.Game.Name, fmt.Sprintf("%s game against %s", instance.Game.Name, instance.PName), instance.PChannelID)
		err := Session.ChannelMessageDelete(instance.PChannelID, instance.CurrentMessageID)
		if err != nil {
			Log.Error(err.Error())
			return
		}

		//Adds the reactions to the message
		addOption(update.Input, instance.PChannelID, newMessage.ID)

	case "global":
		playerEmbed.setColor(instance.Game.Color)
		OpEmbed.setColor(instance.Game.Color)

		//Sends the new message
		newMessage := playerEmbed.send(instance.Game.Name, fmt.Sprintf("%s game against %s", instance.Game.Name, instance.PName), instance.OpChannelID)

		//Adds input field and gameID
		playerEmbed.addField(update.Input.Message, update.Input.Name, false)
		playerEmbed.setFooter(instance.CurrentMessageID, "", "")

		//Adds the reactions to the message
		addOption(update.Input, instance.OpChannelID, newMessage.ID)

	case "playerwin":
		playerEmbed.setColor(Yellow)
		err := Session.ChannelMessageDelete(instance.PChannelID, instance.CurrentMessageID)
		if err != nil {
			Log.Error(err.Error())
		}
		playerEmbed.send("You Won!", fmt.Sprintf("You won your %s game against <@%s>", instance.Game.Name, opponentID), instance.PChannelID)
		playerEmbed.setColor(Red)
		playerEmbed.send("You Lost!", fmt.Sprintf("You lost your %s game against <@%s>", instance.Game.Name, playerID), instance.OpChannelID)

	case "opwin":
		playerEmbed.setColor(Red)
		err = Session.ChannelMessageDelete(instance.PChannelID, currentGameID)
		if err != nil {
			Log.Error(err.Error())
		}
		playerEmbed.send("You Won!", fmt.Sprintf("You won your %s game against %s", info.Name, playerID), opponent.Username)
		playerEmbed.setColor(Red)
		playerEmbed.send("You Lost!", fmt.Sprintf("You lost your %s game against %s", info.Name, opponentID), player.Username)

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
