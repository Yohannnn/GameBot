package core

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
	"strings"
)

// Handler for handling commands
func commandHandler(s *discordgo.Session, m *discordgo.MessageCreate) {
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

	switch command {
	case "ping":
		_, err := s.ChannelMessageSend(m.ChannelID, "pong")
		if err != nil {
			fmt.Println(err)
		}
	// Switch case for handling commands
	case "test":
		embed := NewEmbed()
		embed.SetTitle("Test Title")
		embed.SetDescription("Test description")
		embed.AddField("Test name", "Test value", true)
		embed.SetColor(0)
		_, err := Session.ChannelMessageSendEmbed(m.ChannelID, embed.MessageEmbed)
		if err != nil {
			Log.Error(err.Error())
		}
	case "playgame":
		var invChannelID string

		//Pares args for game
		game := Games[args[0]]
		embed := NewEmbed()
		embed.SetTitle(fmt.Sprintf("%s invite!", game.Name))
		embed.SetColor(7909721)

		//Sets invite channel and description
		if len(args) > 2 {
			embed.SetDescription(fmt.Sprintf("%s invited you to play %s. React with :white_check_mark: to accept.", m.Author.Username, game.Name))
		} else {
			embed.SetDescription(fmt.Sprintf("%s invited anyone %s. React with :white_check_mark: to accept.", m.Author.Username, game.Name))
			invChannelID = m.ChannelID
		}

		_, err := Session.ChannelMessageSendEmbed(invChannelID, embed.MessageEmbed)
		if err != nil {
			Log.Error(err.Error())
		}
	case "gameinfo":
		//Parses args for game
		game := Games[args[0]]

		//Formats embed message
		embed := NewEmbed()
		embed.SetTitle(game.Name)
		embed.SetDescription(game.Description)
		embed.AddField("Rules", game.Rules, true)
		embed.AddField("Board", "", true)
		embed.SetColor(game.Color)

		//Sends embed message
		_, err := Session.ChannelMessageSendEmbed(m.ChannelID, embed.MessageEmbed)
		if err != nil {
			Log.Error(err.Error())
		}
	}
}
