package core

import "github.com/bwmarrin/discordgo"

//reactionHandler
//Handles reactions for messages the bot has sent
func reactionHandler(s *discordgo.Session, r discordgo.MessageReactionAdd) {
	//Reactions added by the bot
	var options []string

	m, err := s.ChannelMessage(r.ChannelID, r.MessageID)
	if err != nil {
		Log.Error(err.Error())
	}

	//Ignore messages that are not sent by the bot
	if m.Author.ID != s.State.User.ID {
		return
	}

	//Parses the message for all options given by the bot
	for _, e := range m.Reactions {
		if e.Me {
			options = append(options, e.Emoji.ID)
		}
	}
}
