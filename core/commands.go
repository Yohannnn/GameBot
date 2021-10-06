package core

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
	"strings"
)

//Handler for handling commands
func commandHandler(s *discordgo.Session, m *discordgo.MessageCreate) {
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
	case "playgame":
		//Pares args for gameInfo
		game := Games[args[0]].Info
		//Formats embed message
		embed := newEmbed()
		embed.setColor(7909721)
		//Sets invite channel and description
		if len(args) > 2 {
			newMessage = embed.send(m.ChannelID, fmt.Sprintf("%s invite!", game.Name), fmt.Sprintf(
				"%s invited you to play %s! React with ✅ to accept.",
				m.Author.Mention(), game.Name))
		} else {
			newMessage = embed.send(m.ChannelID, fmt.Sprintf("%s Invite!", game.Name),
				fmt.Sprintf("%s invited anyone to play %s! React with ✅ to accept.",
					m.Author.Mention(), game.Name))
		}
		//Adds confirmation emoji
		err := Session.MessageReactionAdd(m.ChannelID, newMessage.ID, "✅")
		if err != nil {
			Log.Error(err.Error())
			return
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
		gameInfo := Games[args[0]].Info
		//Formats and sends embed message
		embed.addField("Rules", gameInfo.Rules, true)
		embed.addField("Board", formatBoard(gameInfo.ExampleBoard), true)
		embed.setColor(gameInfo.Color)
		embed.send(m.ChannelID, gameInfo.Name, gameInfo.Description)
	}
}
