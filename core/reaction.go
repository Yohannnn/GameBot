package core

import (
	"github.com/bwmarrin/discordgo"
	"github.com/google/uuid"
	"strings"
)

//Input
//An input for a game update
type Input struct {
	Name    string
	Message string
	Options []string
}

//Output
//The output of a games input (the things that were selected and by whom)
type Output struct {
	Name       string
	SelOptions []string
}

//CreateInput
//Creates an option
func CreateInput(name string, message string, options []string) Input {
	return Input{
		Name:    name,
		Message: message,
		Options: options,
	}
}

//addInput
//Adds an Input to a message
func addInput(option Input, channelID string, messageID string) {
	for _, e := range option.Options {
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
	var output Output

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

	//Gets the game
	game := Games[strings.ToLower(strings.Split(m.Embeds[0].Title, " ")[0])]

	//Checks the message is a game invite
	if len(strings.Split(m.Embeds[0].Title, " ")) > 1 {
		//Checks if the reaction is from the opponent
		Opponent, err := getUser(m.Embeds[0].Description[:21])
		if err != nil {
			log.Error(err.Error())
			return
		}
		if r.UserID == Opponent.ID {
			return
		}

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

		OpponentChan, err := Session.UserChannelCreate(Opponent.ID)
		if err != nil {
			log.Error(err.Error())
			return
		}

		//Creates and defines a new instance
		instance := Instance{
			ID:    uuid.NewString(),
			Game:  game,
			Stats: make(map[string][]string),
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
		Instances[instance.ID] = &instance

		//Starts a new game
		game.StartFunc(&instance)

		//Deletes invite
		err = s.ChannelMessageDelete(m.ChannelID, m.ID)
		if err != nil {
			log.Error(err.Error())
		}
		return
	}

	instance := Instances[m.Embeds[0].Footer.Text]
	output.Name = instance.CurrentInput.Name

	//Gets all valid reactions to the message
	for _, e := range m.Reactions {
		if e.Count == 2 && e.Emoji.Name != "✅" {
			output.SelOptions = append(output.SelOptions, e.Emoji.Name)
		}
	}
	game.UpdateFunc(instance, output)
}
