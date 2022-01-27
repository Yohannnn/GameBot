package core

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
	"strings"
)

//Handler for handling commands
func commandHandler(s *discordgo.Session, m *discordgo.MessageCreate) {
	//Defers panic to error handler
	defer handleCommandError(m.GuildID, m.ChannelID, m.Author.ID)
	var newMessage *discordgo.Message

	//Ignore all messages created by the bot itself
	if m.Author.ID == s.State.User.ID {
		return
	}

	//Checks if the message has the command prefix
	if string(m.Content[0]) != "!" {
		return
	}

	//Parses the message for the command and arguments
	command := strings.ToLower(strings.Split(m.Content, " ")[0][1:])
	args := strings.Split(m.Content, " ")[1:]

	//Switch case for handling commands
	switch command {
	case "test":
		embed := newEmbed()
		embed.setFooter("899872740180369408:899526815238996018", "", "")
		embed.send("test", "test", m.ChannelID)
	case "playgame":
		//Pares args for gameInfo
		game := Games[args[0]]

		//Formats embed message
		embed := newEmbed()
		embed.setColor(Green)

		//Sets invite channel and description
		if len(args) >= 2 {
			newMessage = embed.send(fmt.Sprintf("%s invite!", game.Name), fmt.Sprintf(
				"%s invited %s to play %s! React with ✅ to accept.",
				m.Author.Mention(), args[1], game.Name), m.ChannelID)
		} else {
			newMessage = embed.send(fmt.Sprintf("%s Invite!", game.Name),
				fmt.Sprintf("%s invited anyone to play %s! React with ✅ to accept.",
					m.Author.Mention(), game.Name), m.ChannelID)
		}

		//Adds confirmation emoji
		err := s.MessageReactionAdd(m.ChannelID, newMessage.ID, "✅")
		if err != nil {
			log.Error(err.Error())
			return
		}

		//Deletes the command message
		err = s.ChannelMessageDelete(m.ChannelID, m.ID)
		if err != nil {
			log.Error(err.Error())
		}

	case "gameinfo":
		//Creates a new embed
		embed := newEmbed()
		//Checks if game exits
		if _, ok := Games[args[0]]; !ok {
			embed.setColor(14495300)
			embed.send(m.ChannelID, "Error!", "The requested game does not exist.")
		}
		//Parses args for gameInfo
		game := Games[args[0]]
		//Formats and sends embed message
		embed.addField("Rules", game.Rules, true)
		embed.addField("Board", formatBoard(game.ExampleBoard), true)
		embed.setColor(game.Color)
		embed.send(m.ChannelID, game.Name, game.Description)
	}
}
func handleCommandError(gID string, cId string, uId string) {
	if r := recover(); r != nil {
		log.Errorf("Message from %s in %s in %s caused: %s", uId, cId, gID, r)
	}
	return
}
