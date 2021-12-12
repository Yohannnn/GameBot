package core

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
	"strings"
)

//NumCord
//A map of ints to their emoji (except 0 cause its sus)
var NumCord = map[int]string{
	1:  "1️⃣",
	2:  "2️⃣",
	3:  "3️⃣",
	4:  "4️⃣",
	5:  "5️⃣",
	6:  "6️⃣",
	7:  "7️⃣",
	8:  "8️⃣",
	9:  "9️⃣",
	10: "🔟",
}

//LetCord
//A map of a letters place in the alphabet to its emoji
var LetCord = map[int]string{
	1:  "🇦",
	2:  "🇧",
	3:  "🇨",
	4:  "🇩",
	5:  "🇪",
	6:  "🇫",
	7:  "🇬",
	8:  "🇭",
	9:  "🇮",
	10: "🇯",
	11: "🇰",
	12: "🇱",
	13: "🇲",
	14: "🇳",
	15: "🇴",
	16: "🇵",
	17: "🇶",
	18: "🇷",
	19: "🇸",
	20: "🇹",
	21: "🇺",
	23: "🇻",
	24: "🇼",
	25: "🇽",
	26: "🇾",
	27: "🇿",
}

//Option
//An input for a game update
type Option struct {
	Name      string
	Message   string
	Rollback  bool
	Reactions []string
}

//CreateOption
//Creates an option
func CreateOption(name string, message string, rollback bool, reactions []string) Option {
	return Option{
		Name:      name,
		Message:   message,
		Rollback:  rollback,
		Reactions: reactions,
	}
}

//addOption
//Adds A Reaction Option to a message
func addOption(option Option, channelID string, messageID string) {
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

	//Checks if the reaction was an option given by the bot
	for i, e := range m.Reactions {
		if e.Emoji.ID == r.Emoji.ID && e.Me {
			break
		} else if i == len(m.Reactions) {
			return
		}
	}

	//Checks if the emoji is the confirmation or rollback
	if r.Emoji.Name != "✅" && r.Emoji.Name != "❌" {
		return
	}

	//Checks if the message is a game invite
	if m.Embeds[0].Title[7:] == "Invite!" && r.Emoji.Name == "✅" {
		invitee := strings.Split(m.Embeds[0].Description, " ")[2]
		if invitee == "anyone" || invitee == r.UserID {
			game := Games[strings.Split(m.Embeds[0].Title, " ")[0]]
			OpID := strings.Split(m.Embeds[0].Description, " ")[0]
			startGame(game.StartFunc, *game.Info, r.UserID, OpID)
		}
		return
	}

	//Detects if it is rollback or confirmation
	if r.Emoji.Name == "✅"{
		var reactions []string
		//Gets the options that the player selected
		for _, e := range m.Reactions {
			if e.Me && e.Emoji.Name != "✅"{
				reactions = append(reactions, fmt.Sprintf("%s:%s", e.Emoji.Name, e.Emoji.ID))
			}
		}

		//Gets the option name
		optionName := m.Embeds[0].Fields[len(m.Embeds[0].Fields)-1].Value

		//Gets the board
		for m.Embeds[0].Fields[0]{
			
		}

		}
	}
}