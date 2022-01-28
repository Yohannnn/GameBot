package core

import (
	"github.com/bwmarrin/discordgo"
	"github.com/google/uuid"
	"strings"
)

//Input
//An input for a game update
type Input struct {
	Name      string
	Message   string
	Rollback  bool
	Reactions []string
}

//CreateInput
//Creates an option
func CreateInput(name string, message string, rollback bool, reactions []string) Input {
	return Input{
		Name:      name,
		Message:   message,
		Rollback:  rollback,
		Reactions: reactions,
	}
}

//addInput
//Adds an Input to a message
func addInput(option Input, channelID string, messageID string) {
	if option.Rollback {
		err := Session.MessageReactionAdd(channelID, messageID, "❌")
		if err != nil {
			log.Error(err.Error())
			return
		}
	}

	for _, e := range option.Reactions {
		err := Session.MessageReactionAdd(channelID, messageID, e)
		if err != nil {
			log.Error(err.Error())
			return
		}
	}

	err := Session.MessageReactionAdd(channelID, messageID, "✅")
	if err != nil {
		log.Error(err.Error())
		return
	}
}

//reactionHandler
//Handles reactions for messages the bot has sent
func reactionHandler(s *discordgo.Session, r *discordgo.MessageReactionAdd) {

	//Ignores reactions added by the bot
	if r.UserID == s.State.User.ID {
		return
	}

	//Gets the message that the reaction was put on
	m, err := s.ChannelMessage(r.ChannelID, r.MessageID)
	if err != nil {
		log.Error(err.Error())
	}

	//Ignore messages that are not sent by the bot
	if m.Author.ID != s.State.User.ID {
		return
	}

	//Checks if the emoji is the confirmation emoji
	if r.Emoji.Name != "✅" {
		return
	}

	//Checks the message is a game invite
	if m.Embeds[0].Title[7:] == "Invite!" {
		//Gets the game
		game := Games[strings.Split(m.Embeds[0].Title, " ")[0]]

		//Get the user struct and dm channel for each player
		Current, err := Session.User(r.UserID)
		if err != nil {
			log.Error(err.Error())
			return
		}
		CurrentChan, err := Session.UserChannelCreate(Current.ID)
		if err != nil {
			log.Error(err.Error())
			return
		}
		Opponent := m.Mentions[0]
		OpponentChan, err := Session.UserChannelCreate(Current.ID)
		if err != nil {
			log.Error(err.Error())
			return
		}

		//Creates and defines a new instance
		instance := Instance{
			ID:    uuid.NewString(),
			Game:  game,
			Board: nil,
			Stats: nil,
			Players: []Player{
				{
					ID:        Current.ID,
					Name:      Current.Username,
					ChannelID: CurrentChan.ID,
				}, {
					ID:        Opponent.ID,
					Name:      Opponent.Username,
					ChannelID: OpponentChan.ID,
				}},
		}

		//Starts a new game
		game.StartFunc(&instance)
		return
	}
}
