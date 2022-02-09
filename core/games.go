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

//JSONInstance
//Data for an instance that needs to be written to a json
type JSONInstance struct {
	ID               string              `json:"id"`
	GameName         string              `json:"game_name"`
	Board            [][]string          `json:"board"`
	DisplayBoard     [][]string          `json:"display_board"`
	CurrentInput     Input               `json:"current_input"`
	Stats            map[string][]string `json:"stats"`
	CurrentMessageID string              `json:"current_message_id"`
	Players          []Player            `json:"players"`
	Turn             int                 `json:"turn"`
}

//Player
//The player of a game
type Player struct {
	ID        string
	Name      string
	ChannelID string
}

//TODO create separate player struct that deals that keeps tracks of wins and is stored in a JSON

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
func UpdateGame(instance *Instance, Input Input) {
	Current := instance.Players[instance.Turn]
	Opponent := instance.Players[-(instance.Turn - 1)]

	//Creates new embeds
	Embed := newEmbed()
	SentEmb := newEmbed()

	//Formats and sends message
	Embed.addField("Board", formatBoard(instance.DisplayBoard), true)
	Embed.addField("Input", Input.Message, true)
	Embed.setColor(Blue)
	Embed.setFooter(instance.ID, "", "")
	newMessage := Embed.send(instance.Game.Name, fmt.Sprintf("%s game against %s", instance.Game.Name, Current.Name), Opponent.ChannelID)

	//Formats sent embed
	SentEmb.setColor(Green)
	SentEmb.setFooter(instance.ID, "", "")

	//Edits old message to sent
	SentEmb.edit("Sent!", fmt.Sprintf("Sent %s update to %s", instance.Game.Name, Opponent.Name), Current.ChannelID, instance.CurrentMessageID)

	//Adds the reactions to the message
	addInput(Input, Opponent.ChannelID, newMessage.ID)

	//Changes the turn
	instance.Turn = -(instance.Turn - 1)

	//Sets current input
	instance.CurrentInput = Input

	//Sets current message ID
	instance.CurrentMessageID = newMessage.ID

	//Save instances to JSON
	err := saveInstances()
	if err != nil {
		log.Error(err.Error())
		return
	}
}

//StartGame
//Starts an instance of a game
func StartGame(instance *Instance, Input Input) {
	var Current Player
	var Opponent Player

	Embed := newEmbed()

	Current = instance.Players[0]
	Opponent = instance.Players[1]

	//Formats and sends the message
	Embed.addField("Board", formatBoard(instance.DisplayBoard), true)
	Embed.addField("Input", Input.Message, true)
	Embed.setColor(Blue)
	Embed.setFooter(instance.ID, "", "")
	newMessage := Embed.send(instance.Game.Name, fmt.Sprintf("%s game against %s", instance.Game.Name, Opponent.Name), Current.ChannelID)

	//Adds the options to the message
	addInput(Input, Current.ChannelID, newMessage.ID)

	//Sets current input
	instance.CurrentInput = Input

	//Sets current message ID
	instance.CurrentMessageID = newMessage.ID

	//Save instances to JSON
	err := saveInstances()
	if err != nil {
		log.Error(err.Error())
		return
	}
}

//EndGame
//Ends a game with a winner and loser
func EndGame(instance *Instance, Winner Player, Looser Player) {
	Embed := newEmbed()

	//Formats the embed
	Embed.addField("Board", formatBoard(instance.DisplayBoard), true)

	//Sets color and sends to winner
	Embed.setColor(Yellow)
	Embed.send(instance.Game.Name, fmt.Sprintf("You won your %s game against %s", instance.Game.Name, Looser.Name), Winner.ChannelID)

	//Sets color and sends to looser
	Embed.setColor(Red)
	Embed.send(instance.Game.Name, fmt.Sprintf("You lost your %s game against %s", instance.Game.Name, Winner.Name), Looser.ChannelID)

	//Removes instance from instances
	delete(Instances, instance.ID)

	//Save instances to JSON
	err := saveInstances()
	if err != nil {
		log.Error(err.Error())
		return
	}
}

//EditGame
//Edits a current games message instead of sending a new one
func EditGame(instance *Instance, Input Input) {
	//Gets players
	Current := instance.Players[instance.Turn]
	Opponent := instance.Players[-(instance.Turn - 1)]

	//Creates a new embed
	Embed := newEmbed()

	//Formats the player game board and sets color
	Embed.addField("Board", formatBoard(instance.DisplayBoard), true)
	Embed.addField("Input", Input.Message, true)
	Embed.setColor(Blue)

	//Sends the new message
	newMessage := Embed.send(instance.Game.Name, fmt.Sprintf("%s game against %s", instance.Game.Name, Opponent.Name), Current.ChannelID)

	//Deletes the old message
	err := Session.ChannelMessageDelete(Current.ChannelID, instance.CurrentMessageID)
	if err != nil {
		log.Error(err.Error())
		return
	}

	//Adds input field and gameID
	Embed.addField(Input.Message, Input.Name, false)
	Embed.setFooter(instance.ID, "", "")

	//Adds the reactions to the message
	addInput(Input, Current.ChannelID, newMessage.ID)

	//Sets current input
	instance.CurrentInput = Input

	//Sets current message ID
	instance.CurrentMessageID = newMessage.ID

	//Save instances to JSON
	err = saveInstances()
	if err != nil {
		log.Error(err.Error())
		return
	}
}
