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
	StartFunc    func() (*Instance, GameUpdate)
	UpdateFunc   func(*Instance) GameUpdate
}

//Instance
//An instance of a game
type Instance struct {
	ID               string
	Game             Game
	Board            [][]string
	Stats            map[string]string
	Options          []string
	Turn             bool
	CurrentMessageID string
	CurrentInput     Input
	Players          []Player
}

//Player
//The player of a game
type Player struct {
	ID        string
	Name      string
	ChannelID string
}

//GameUpdate
//An update to a game
type GameUpdate struct {
	Type  string
	Board [][]string
	Input Input
}

//Games
//Map games names to their game struct
var Games = make(map[string]Game)

//Instances
//Map of IDs to their game instance
var Instances = make(map[string]*Instance)

//AddGame
//Adds a game to the game map
func AddGame(Name string, Description string, Rules string, Color int, ExampleBoard [][]string, StartFunc func() (*Instance, GameUpdate), UpdateFunc func(*Instance) GameUpdate) {
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
func CreateGameUpdate(Type string, Board [][]string, Input Input) GameUpdate {
	//Checks type
	if !Contains([]string{"local", "global", "error", "playerwin", "opponentwin"}, Type) {
		Log.Error("Invalid type for game update")
		return GameUpdate{}
	}

	gU := GameUpdate{
		Type:  Type,
		Board: Board,
		Input: Input,
	}

	return gU
}

//gameUpdate
//Sends the game update
func gameUpdate(instance *Instance, update GameUpdate) {
	var Current Player
	var Opponent Player

	if instance.Turn {
		Current = instance.Players[0]
		Opponent = instance.Players[1]
	} else {
		Current = instance.Players[1]
		Opponent = instance.Players[0]
	}

	//Creates a new embed
	Embed := newEmbed()

	//Formats the player game board and sets color
	Board := formatBoard(update.Board)
	Embed.addField("Board", Board, true)
	Embed.setColor(instance.Game.Color)

	//Checks the type of update
	switch update.Type {

	case "local":
		//Adds option field and gameID
		Embed.addField(update.Input.Message, update.Input.Name, false)
		Embed.setFooter(instance.ID, "", "")

		//Sends the new message and deletes the old one
		newMessage := Embed.send(instance.Game.Name, fmt.Sprintf("%s game against %s", instance.Game.Name, Opponent.Name), Current.ChannelID)
		err := Session.ChannelMessageDelete(Current.ChannelID, instance.CurrentMessageID)
		if err != nil {
			Log.Error(err.Error())
			return
		}

		//Adds the reactions to the message
		addInput(update.Input, Current.ChannelID, newMessage.ID)

	case "global":
		//Sends the new message
		newMessage := Embed.send(instance.Game.Name, fmt.Sprintf("%s game against %s", instance.Game.Name, Current.Name), Opponent.ChannelID)

		//Adds input field and gameID
		Embed.addField(update.Input.Message, update.Input.Name, false)
		Embed.setFooter(instance.CurrentMessageID, "", "")

		//Adds the reactions to the message
		addInput(update.Input, Opponent.ChannelID, newMessage.ID)
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
