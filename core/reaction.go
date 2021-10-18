package core

import (
	"github.com/bwmarrin/discordgo"
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
	Name     string
	Type     string
	Message  string
	Rollback bool
	Cord     struct {
		x     [2]int
		y     [2]int
		value [2]int
	}
	Select struct {
		Options []string
		value   string
	}
}

//addAReactionOption to a message
func addOption(option Option, channelID string, messageID string) {
	if option.Rollback {
		err := Session.MessageReactionAdd(channelID, messageID, "❌")
		if err != nil {
			Log.Error(err.Error())
			return
		}
	}

	switch option.Type {
	case "select":
		for _, e := range option.Select.Options {
			err := Session.MessageReactionAdd(channelID, messageID, e)
			if err != nil {
				Log.Error(err.Error())
				return
			}
		}
	case "cord":
		xArray := IntArray(option.Cord.x[0], option.Cord.x[1])
		for _, x := range xArray {
			err := Session.MessageReactionAdd(channelID, messageID, NumCord[x])
			if err != nil {
				Log.Error(err.Error())
			}
		}
	}

	err := Session.MessageReactionAdd(channelID, messageID, "✅")
	if err != nil {
		Log.Error(err.Error())
		return
	}
}

//reactionHandler
//Handles reactions for messages the bot has sent
func reactionHandler(s *discordgo.Session, r *discordgo.MessageReactionAdd) {
	//Gets the message that the reaction was put on
	m, err := s.ChannelMessage(r.ChannelID, r.MessageID)
	if err != nil {
		Log.Error(err.Error())
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
		//game := Games[strings.Split(m.Embeds[0].Title, " ")[0]]
		//startUpdate := game.StartFunc
		//sendGameUpdate(game.Info, startUpdate(), m.Embeds[0].Description)
	}

	//Checks if the reaction was an option given by the bot
	for i, e := range m.Reactions {
		if e.Emoji.ID == r.Emoji.ID && e.Me {
			break
		} else if i == len(m.Reactions) {
			return
		}
	}
}
