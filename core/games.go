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
	ExampleBoard [][]string
	StartFunc    func(*Instance)
	UpdateFunc   func(*Instance, Output)
}

//Instance
//An instance of a game
type Instance struct {
	ID               string
	Game             Game
	Board            [][]string
	DisplayBoard     [][]string
	CurrentInput     Input
	Stats            map[string][]string
	CurrentMessageID string
	Players          []Player
	Turn             int
}

//Player
//The player of a game
type Player struct {
	ID        string
	Name      string
	ChannelID string
}

//Games
//Map games names to their game struct
var Games = make(map[string]Game)

//Instances
//Map of IDs to their game instance
var Instances = make(map[string]*Instance)

//AddGame
//Adds a game to the game map
func AddGame(Name string, Description string, Rules string, ExampleBoard [][]string, StartFunc func(*Instance), UpdateFunc func(*Instance, Output)) {
	game := Game{
		Name:         Name,
		Description:  Description,
		Rules:        Rules,
		ExampleBoard: ExampleBoard,
		StartFunc:    StartFunc,
		UpdateFunc:   UpdateFunc,
	}
	Games[strings.ToLower(Name)] = game
}

//UpdateGame
//Sends an update of a game instance
func UpdateGame(instance *Instance, Board [][]string, Input Input) {
	Current := instance.Players[instance.Turn]
	Opponent := instance.Players[-instance.Turn-1]

	//Creates a new embed
	Embed := newEmbed()

	//Formats the player game board and sets color
	Embed.addField("Board", formatBoard(Board), true)
	Embed.addField("Input", Input.Message, true)
	Embed.setColor(Blue)

	Embed.send(instance.Game.Name, fmt.Sprintf("%s game against %s", instance.Game.Name, Opponent.Name), Current.ChannelID)
	newMessage := Embed.send(instance.Game.Name, fmt.Sprintf("%s game against %s", instance.Game.Name, Current.Name), Opponent.ChannelID)

	//Adds input field and gameID
	Embed.addField(Input.Message, Input.Name, false)
	Embed.setFooter(instance.ID, "", "")

	//Adds the reactions to the message
	addInput(Input, Opponent.ChannelID, newMessage.ID)

	//Changes the turn
	instance.Turn = -instance.Turn - 1
}

//StartGame
//Starts an instance of a game
func StartGame(instance *Instance, Board [][]string, Input Input) {
	var Current Player
	var Opponent Player

	Embed := newEmbed()

	Current = instance.Players[instance.Turn]
	Opponent = instance.Players[-instance.Turn-1]

	//Formats and sends the message
	Embed.addField("Board", formatBoard(Board), true)
	Embed.addField("Input", Input.Message, true)
	Embed.setColor(Blue)
	Embed.setFooter(instance.ID, "", "")
	newMessage := Embed.send(instance.Game.Name, fmt.Sprintf("%s game against %s", instance.Game.Name, Current.Name), Opponent.ChannelID)

	//Adds input field and gameID
	Embed.addField(Input.Message, Input.Name, false)

	//Adds the options to the message
	addInput(Input, Opponent.ChannelID, newMessage.ID)
}

//EndGame
//Ends a game with a winner and loser
func EndGame(instance *Instance, Winner Player, Looser Player, Board [][]string) {
	Embed := newEmbed()

	//Formats the embed
	Embed.addField("Board", formatBoard(Board), true)

	//Sets color and sends to winner
	Embed.setColor(Yellow)
	Embed.send(instance.Game.Name, fmt.Sprintf("%s game against %s", instance.Game.Name, Looser.Name), Winner.ChannelID)

	//Sets color and sends to looser
	Embed.setColor(Red)
	Embed.send(instance.Game.Name, fmt.Sprintf("%s game against %s", instance.Game.Name, Winner.Name), Looser.ChannelID)

	//Removes instance from instances
	delete(Instances, instance.ID)
}

//EditGame
//Edits a current games message instead of sending a new one
func EditGame(instance *Instance, Board [][]string, Input Input) {
	//Gets players
	Current := instance.Players[instance.Turn]
	Opponent := instance.Players[-instance.Turn-1]

	//Creates a new embed
	Embed := newEmbed()

	//Formats the player game board and sets color
	Embed.addField("Board", formatBoard(Board), true)
	Embed.setColor(Blue)

	//Deletes the old message
	err := Session.ChannelMessageDelete(Current.ChannelID, instance.CurrentMessageID)
	if err != nil {
		log.Error(err.Error())
		return
	}

	//Sends the new message
	newMessage := Embed.send(instance.Game.Name, fmt.Sprintf("%s game against %s", instance.Game.Name, Opponent.Name), Current.ChannelID)

	//Adds input field and gameID
	Embed.addField(Input.Message, Input.Name, false)
	Embed.setFooter(instance.ID, "", "")

	//Adds the reactions to the message
	addInput(Input, Opponent.ChannelID, newMessage.ID)
}
