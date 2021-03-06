package core

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
	"strings"
	"sync"
)

// TODO Switch to modular commands
// TODO Add admin commands

// commandLock
// Used to check if command handler is running
var commandLock = &sync.Mutex{}

// Command
// Struct that contains data for a command
type Command struct {
	Admin    bool
	Name     string
	Function CommandFunc
}

// Ctx
// Context for a command
type Ctx struct {
	MessageID string
	ChannelID string
	Args      []string
	Author    *discordgo.User
}

// CommandFunc
// Function for a command
type CommandFunc func(*Ctx)

// Commands
// Map of triggers to their corresponding command
var Commands = make(map[string]*Command)

// commandHandler
// Handler for handling commands
func commandHandler(s *discordgo.Session, m *discordgo.MessageCreate) {
	// Graceful termination check
	if graceTerm {
		return
	}
	commandLock.Lock()
	defer commandLock.Unlock()

	// Defers panic to error handler
	defer handleCommandError(m.GuildID, m.ChannelID, m.Author.ID)
	var newMessage *discordgo.Message

	// Ignore all messages created by the bot itself
	if m.Author.ID == s.State.User.ID {
		return
	}

	// Checks if the message has the command prefix
	if string(m.Content[0]) != "!" {
		return
	}

	// Parses the message for the command and arguments
	command := strings.ToLower(strings.Split(m.Content, " ")[0][1:])
	args := strings.Split(m.Content, " ")[1:]

	// Switch case for handling commands
	switch command {

	// Gives a list of all games
	case "gamelist":
		var names string
		for _, g := range Games {
			names += g.Name + ", "
		}

		Embed := newEmbed()
		Embed.setColor(Blue)
		_, err := Embed.send("Games", fmt.Sprintf("The currently playable games are: %s\"To see more information about a particular game you can run !GameInfo <GameName>\"", names), m.ChannelID)
		if err != nil {
			log.Error(err.Error())
			return
		}

		// Deletes the command message
		err = s.ChannelMessageDelete(m.ChannelID, m.ID)
		if err != nil {
			log.Error(err.Error())
		}

	// Creates an invitation to a game
	case "gameinvite":
		// Pares args for gameInfo
		game := Games[strings.ToLower(args[0])]

		// Formats Embed message
		Embed := newEmbed()
		Embed.setColor(Green)

		// Sets invite channel and description
		if len(args) >= 2 {
			var err error
			newMessage, err = Embed.send(fmt.Sprintf("%s invite!", game.Name), fmt.Sprintf(
				"%s invited %s to play %s! React with ??? to accept.",
				m.Author.Mention(), args[1], game.Name), m.ChannelID)
			if err != nil {
				log.Error(err.Error())
				return
			}
		} else {
			var err error
			newMessage, err = Embed.send(fmt.Sprintf("%s Invite!", game.Name),
				fmt.Sprintf("%s invited anyone to play %s! React with ??? to accept.",
					m.Author.Mention(), game.Name), m.ChannelID)
			if err != nil {
				log.Error(err.Error())
				return
			}
		}

		// Adds confirmation emoji
		err := s.MessageReactionAdd(m.ChannelID, newMessage.ID, "???")
		if err != nil {
			log.Error(err.Error())
			return
		}

		// Deletes the command message
		err = s.ChannelMessageDelete(m.ChannelID, m.ID)
		if err != nil {
			log.Error(err.Error())
		}

	// Gives the info for a game
	case "gameinfo":
		// Creates a new Embed
		Embed := newEmbed()
		// Checks if game exits
		if _, ok := Games[args[0]]; !ok {
			Embed.setColor(14495300)
			_, err := Embed.send(m.ChannelID, "Error!", "The requested game does not exist.")
			if err != nil {
				log.Error(err.Error())
				return
			}
		}
		// Parses args for gameInfo
		game := Games[args[0]]
		// Formats and sends Embed message
		Embed.addField("Rules", game.Rules, true)
		Embed.addField("Board", formatBoard(game.ExampleBoard), true)
		Embed.setColor(Blue)
		_, err := Embed.send(game.Name, game.Description, m.ChannelID)
		if err != nil {
			log.Error(err.Error())
			return
		}

		// Deletes the command message
		err = s.ChannelMessageDelete(m.ChannelID, m.ID)
		if err != nil {
			log.Error(err.Error())
		}
	}
}

// handleCommandError
// Handles errors for commands
func handleCommandError(gID string, cId string, uId string) {
	if r := recover(); r != nil {
		log.Errorf("Message from %s in %s in %s caused: %s", uId, cId, gID, r)
	}
	return
}
