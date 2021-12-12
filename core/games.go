package core

import (
	"fmt"
	"strings"
)

//Game
//The update function, start function, and game information for a game
type Game struct {
	UpdateFunc func(GameInput) (GameUpdate, string)
	StartFunc  func() GameStart
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
	OptionName  string
	Reactions   []string
	PlayerState GameState
	OpSate      GameState
}

//GameUpdate
//An update to a game
type GameUpdate struct {
	Type        string
	ErrMessage string
	PlayerState GameState
	OpState     GameState
	Option      Option
}

//GameStart
//The start data for a game
type GameStart struct {
	Type        string
	State GameState
	Option Option
}

//Games
//Map games names to their game struct
var Games = make(map[string]Game)

//AddGame
//Adds a game to the game map
func AddGame(updateFunc func(GameInput) (GameUpdate, string), startFunc func() GameStart, gI *GameInfo) {
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

//CreateGameSate
//Creates a game state
func CreateGameSate(board [][]string, stats map[string]string) GameState{
	gS := GameState{
		Board: board,
		Stats: stats,
	}
	return gS
}

//createGameInput
//Creates the input to a game func (only used by framework)
func createGameInput(OptionName string, Reactions []string, PlayerState GameState, OpState GameState) GameInput {
	gI := GameInput{
		OptionName:  OptionName,
		Reactions:   Reactions,
		PlayerState: PlayerState,
		OpSate:      OpState,
	}
	return gI
}

//CreateGameUpdateLocal
//Creates a local game update to be sent to a player
func CreateGameUpdateLocal(PlayerState GameState, Option Option) GameUpdate {
	gU := GameUpdate{
		Type:        "local",
		PlayerState: PlayerState,
		Option:      Option,
	}
	return gU
}

//CreateGameUpdateGlobal
//Creates a global game update to be sent to the opponent and the player
func CreateGameUpdateGlobal(PlayerState GameState, OpState GameState ,Option Option) GameUpdate {
	gU := GameUpdate{
		Type:        "local",
		PlayerState: PlayerState,
		OpState: OpState,
		Option:      Option,
	}
	return gU
}

//CreateGameUpdateWin
//Creates a win game update to be sent to the opponent and the player
func CreateGameUpdateWin(PlayerState GameState, OpState GameState ,Option Option, Winner bool) GameUpdate {
	gU := GameUpdate{
		PlayerState: PlayerState,
		OpState: OpState,
		Option:      Option,
	}
	if Winner{
		gU.Type = "playerwin"
	}else{
		gU.Type = "opwin"
	}

	return gU
}

//CreateGameUpdateErr
//Creates an err game update to be sent to the player
func CreateGameUpdateErr(Message string, Option Option) GameUpdate {
	gU := GameUpdate{
		ErrMessage:  Message,
		Option: Option,
	}

	return gU
}

//sendGameUpdate
//Sends the game update
func sendGameUpdate(info *GameInfo, update GameUpdate, playerID string, opID string, currentGameID string) {
	var stats string
	var board string
	var opEmbed embed
	//Gets the user struct of each player
	player, err := Session.User(playerID)
	if err != nil {
		log.Error(err.Error())
	}
	opponent, err := Session.User(opID)
	if err != nil {
		log.Error(err.Error())
	}

	//Gets the dm channel for each player
	playerChannel, err := Session.UserChannelCreate(player.ID)
	if err != nil {
		log.Error(err.Error())
	}
	opChannel, err := Session.UserChannelCreate(opponent.ID)
	if err != nil {
		log.Error(err.Error())
	}

	//Gets the current game message
	message, err := Session.ChannelMessage(playerChannel.ID, currentGameID)
	if err != nil {
		log.Error(err.Error())
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

	//Formats opponent state
	if update.Type != "local" && update.Type != "global" {
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
		playerEmbed.setColor(info.Color)
		//Adds the selected options field
		for _, field := range currentEmbed.Fields[2:] {
			if field.Name == "Selected Options" {
				playerEmbed.addField("Selected Options", field.Value, false)
				break
			}
		}

		//Adds option field
		gameID := strings.Split(currentEmbed.Footer.Text, ":")[0]
		playerEmbed.addField(update.Option.Message, update.Option.Name, false)

		//Sets the footer
		if update.Option.Rollback{
			playerEmbed.setFooter(fmt.Sprintf("%s:%s", gameID, message.ID), "", "")
		}else{
			playerEmbed.setFooter(gameID, "", "")

			//Traces back old local updates and deletes them
			tracedEmbed := currentEmbed
			tracedMessageID := message.ID
			for len(strings.Split(tracedEmbed.Footer.Text, ":")) > 1{
				err := Session.ChannelMessageDelete(playerChannel.ID, tracedMessageID)
				if err != nil{
					log.Error(err.Error())
					return
				}
				tracedMessageID = strings.Split(tracedEmbed.Footer.Text, ":")[1]
				tracedMessage, err := Session.ChannelMessage(playerChannel.ID, tracedMessageID)
				tracedEmbed = tracedMessage.Embeds[0]
			}
		}

		//Sends the new message and deletes the old one
		newMessage := playerEmbed.send(info.Name,
			fmt.Sprintf("%s game against %s", info.Name, opponent.Username), playerChannel.ID)
		err = Session.ChannelMessageDelete(playerChannel.ID, currentGameID)
		if err != nil {
			log.Error(err.Error())
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

		//Edits the players embed
		_, err = Session.ChannelMessageEditEmbed(playerChannel.ID, currentGameID, playerEmbed.MessageEmbed)
		if err != nil{
			log.Error(err.Error())
			return
		}

		//Sends the new message
		newMessage := opEmbed.send(info.Name, fmt.Sprintf("%s game against %s", info.Name, player.Username), opChannel.ID)

		//Adds option field and gameID
		opEmbed.addField(update.Option.Message, update.Option.Name, false)
		opEmbed.setFooter(currentGameID, "", "")

		//Adds the reactions to the message
		addOption(update.Option, opChannel.ID, newMessage.ID)
		return

	case "playerwin":
		//Sets colors
		playerEmbed.setColor(Yellow)
		opEmbed.setColor(Red)

		//Deletes old message
		err = Session.ChannelMessageDelete(playerChannel.ID, currentGameID)
		if err != nil {
			log.Error(err.Error())
		}

		//Sends new messages
		playerEmbed.send("You Won!",
			fmt.Sprintf("You won your %s game against %s", info.Name, opponent.Username), playerChannel.ID)
		opEmbed.send("You Lost!",
			fmt.Sprintf("You lost your %s game against %s", info.Name, player.Username), opChannel.ID)
		return

	case "opwin":
		//Sets the colors
		playerEmbed.setColor(Red)
		opEmbed.setColor(Yellow)

		//Deletes the old message
		err = Session.ChannelMessageDelete(playerChannel.ID, currentGameID)
		if err != nil {
			log.Error(err.Error())
		}

		//Sends the new messages
		opEmbed.send("You Won!",
			fmt.Sprintf("You won your %s game against %s", info.Name, player.Username), opChannel.ID)
		playerEmbed.send("You Lost!",
			fmt.Sprintf("You lost your %s game against %s", info.Name, opponent.Username), playerChannel.ID)
		return

	case "err":
		//Sets the color
		playerEmbed.setColor(info.Color)

		//Adds the error field
		playerEmbed.addField("There was an error", update.ErrMessage, false)

		//Adds the selected options field
		for _, field := range currentEmbed.Fields[2:] {
			if field.Name == "Selected Options" {
				playerEmbed.addField("Selected Options", field.Value, false)
				break
			}
		}

		//Deletes the old message
		err = Session.ChannelMessageDelete(playerChannel.ID, currentGameID)
		if err != nil {
			log.Error(err.Error())
		}

		//Sends the new message
		newMessage := playerEmbed.send(info.Name,
			fmt.Sprintf("%s game against %s", info.Name, opponent.Username), playerChannel.ID)

		//Adds the option to the message
		addOption(update.Option, playerChannel.ID, newMessage.ID)
		return
	}
}

//startGame
//Starts a game
func startGame(startFunc func() GameStart, info GameInfo, playerID string, opID string){
	//Runs the games start function
	gameStart := startFunc()

	//Gets the dm channel for the player
	playerChannel, err := Session.UserChannelCreate(playerID)
	if err != nil {
		log.Error(err.Error())
	}

	//Gets the user struct of the opponent
	op, err := Session.User(opID)
	if err != nil{

	}

	//Creates a new embed and sets the color
	embed := newEmbed()
	embed.setColor(info.Color)

	//Formats the board
	board := formatBoard(gameStart.State.Board)
	embed.addField("Board", board, true)

	//Formats the stats
	stats := formatStats(gameStart.State.Stats)
		embed.addField("Stats", stats, true)

	//Sends the embed
	embed.send(info.Name, fmt.Sprintf("%s game against %s", info.Name, op.Username), playerChannel.ID)
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

//formatStats
//Formats game stats into a string
func formatStats(stats map[string]string) string {
	var statsString string
	for stat, value := range stats {
		statsString += fmt.Sprintf("%s = %s\n", stat, value)
	}
	return statsString
}
