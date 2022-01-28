package core

import (
	"fmt"
	"github.com/google/uuid"
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
	UpdateFunc   func(*Instance)
}

//Instance
//An instance of a game
type Instance struct {
	ID               string
	Game             Game
	Board            [][]string
	Stats            map[string]string
	CurrentMessageID string
	Players          []Player
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
func AddGame(Name string, Description string, Rules string, ExampleBoard [][]string, StartFunc func(*Instance), UpdateFunc func(*Instance)) {
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
func UpdateGame(instance *Instance, Board [][]string, Input Input, Current Player, Opponent Player) {
	//Creates a new embed
	Embed := newEmbed()

	//Formats the player game board and sets color
	Embed.addField("Board", formatBoard(Board), true)
	Embed.setColor(Blue)

	//Sends the new messages and deletes old one
	err := Session.ChannelMessageDelete(Current.ChannelID, instance.CurrentMessageID)
	if err != nil {
		log.Error(err.Error())
		return
	}

	Embed.send(instance.Game.Name, fmt.Sprintf("%s game against %s", instance.Game.Name, Opponent.Name), Current.ChannelID)
	newMessage := Embed.send(instance.Game.Name, fmt.Sprintf("%s game against %s", instance.Game.Name, Current.Name), Opponent.ChannelID)

	//Adds input field and gameID
	Embed.addField(Input.Message, Input.Name, false)
	Embed.setFooter(instance.ID, "", "")

	//Adds the reactions to the message
	addInput(Input, Opponent.ChannelID, newMessage.ID)
}

//StartGame
//Starts an instance of a game
func StartGame(game Game, Current Player, Opponent Player, Board [][]string, Input Input) {
	Embed := newEmbed()
	instance := Instance{ID: uuid.NewString(), Players: []Player{Current, Opponent}}
	Instances[instance.ID] = &instance

	//Formats and sends the message
	Embed.addField("Board", formatBoard(Board), true)
	Embed.setColor(Blue)
	Embed.setFooter(instance.ID, "", "")
	newMessage := Embed.send(game.Name, fmt.Sprintf("%s game against %s", game.Name, Opponent.Name), Current.ChannelID)

	//Adds input field and gameID
	Embed.addField(Input.Message, Input.Name, false)

	//Adds the options to the message
	addInput(Input, Opponent.ChannelID, newMessage.ID)
}
