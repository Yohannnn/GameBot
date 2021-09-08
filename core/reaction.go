package core

import "github.com/bwmarrin/discordgo"

//reactionHandler
//Handles reactions for messages the bot has sent
func reactionHandler(s *discordgo.Session, r discordgo.MessageReactionAdd) {
	message, err := s.ChannelMessage(r.ChannelID, r.MessageID)
	if err != nil {
		Log.Error(err.Error())
	}

	// If the message is not sent by the bot ignores the reaction
	if message.Author.ID != s.State.User.ID {
		return
	}

	//
}
