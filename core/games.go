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
	Color        int
	ExampleBoard [][]string
	StartFunc    func(*Instance) GameUpdate
	UpdateFunc   func(*Instance) GameUpdate
}

//Instance
//An instance of a game
type Instance struct {
	ID               string
	Game             Game
	Board            [][]string
	Stats            map[string]string
	Turn             int
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
	Current  Player
	Opponent Player
	Board    [][]string
	Input    Input
}

//Games
//Map games names to their game struct
var Games = make(map[string]Game)

//Instances
//Map of IDs to their game instance
var Instances = make(map[string]*Instance)

//AddGame
//Adds a game to the game map
func AddGame(Name string, Description string, Rules string, Color int, ExampleBoard [][]string, StartFunc func(*Instance) GameUpdate, UpdateFunc func(*Instance) GameUpdate) {
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
func CreateGameUpdate(Player Player, Board [][]string, Input Input) GameUpdate {
	gU := GameUpdate{
		Board: Board,
		Input: Input,
	}
	return gU
}

//CreateInstance
//Creates an instance
func CreateInstance(ID string, Game Game, Board [][]string, stats map[string]string, Turn int, CurrentMessageID string, CurrentInput Input, Players []Player) *Instance {
	newInstance := Instance{
		ID:               ID,
		Game:             Game,
		Board:            Board,
		Stats:            stats,
		Turn:             Turn,
		CurrentMessageID: CurrentMessageID,
		CurrentInput:     CurrentInput,
		Players:          Players,
	}
	Instances[ID] = &newInstance
	return &newInstance
}

//gameUpdate
//Sends the game update
func gameUpdate(instance *Instance, update GameUpdate) {
	//Creates a new embed
	Embed := newEmbed()

	//Formats the player game board and sets color
	Embed.addField("Board", formatBoard(update.Board), true)
	Embed.setColor(instance.Game.Color)

	//Sends the new messages and deletes old one
	err := Session.ChannelMessageDelete(update.Current.ChannelID, instance.CurrentMessageID)
	if err != nil {
		log.Error(err.Error())
		return
	}

	Embed.send(instance.Game.Name, fmt.Sprintf("%s game against %s", instance.Game.Name, update.Opponent.Name), update.Current.ChannelID)
	newMessage := Embed.send(instance.Game.Name, fmt.Sprintf("%s game against %s", instance.Game.Name, update.Current.Name), update.Opponent.ChannelID)

	//Adds input field and gameID
	Embed.addField(update.Input.Message, update.Input.Name, false)
	Embed.setFooter(instance.ID, "", "")

	//Adds the reactions to the message
	addInput(update.Input, update.Opponent.ChannelID, newMessage.ID)
}

//startGame
//starts am instance of a game
func startGame(game Game, update GameUpdate, Current Player, Opponent Player) *Instance {
	//Creates a new embed and instance ID
	Embed := newEmbed()
	ID := uuid.NewString()
	update := game.StartFunc()

	//Formats and sends the message
	Embed.addField("Board", formatBoard(update.Board), true)
	Embed.setColor(game.Color)
	Embed.setFooter(ID, "", "")
	newMessage := Embed.send(game.Name, fmt.Sprintf("%s game against %s", game.Name, update.Current.Name), update.Opponent.ChannelID)

	//Adds input field and gameID
	Embed.addField(update.Input.Message, update.Input.Name, false)

	//Adds the options to the message
	addInput(update.Input, update.Opponent.ChannelID, newMessage.ID)

	//Sets the instance ID

	//Creates and returns an instance
	return instance{
		ID:               ID,
		Game:             game,
		Board:            nil,
		Stats:            nil,
		Turn:             0,
		CurrentMessageID: newMessage.ID,
		CurrentInput:     update.Input,
		Players:          nil,
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
				log.Error(err.Error())
				return ""
			}
			LineString += emoji.MessageFormat()
		}
		BoardString += LineString + "\n"
		LineString = ""
	}
	return BoardString
}
