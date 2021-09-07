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
	//args := strings.Split(m.Content, " ")[1:]

	switch command {
	case "ping":
		_, err := s.ChannelMessageSend(m.ChannelID, "pong")
		if err != nil {
			fmt.Println(err)
		}
	// Switch case for handling commands
	case "test":

	case "playgame":
	case "gameinfo":
	}
}
