package core

import "github.com/bwmarrin/discordgo"

//Colors
const (
	Red    = 14495300
	Brown  = 12675407
	Orange = 16027660
	Yellow = 16632664
	Green  = 7909721
	Blue   = 5614830
	Purple = 11177686
	White  = 15132648
)

//embed
//Struct for response
type embed struct {
	*discordgo.MessageEmbed
}

//newEmbed
//Returns a new embed object
func newEmbed() *embed {
	return &embed{&discordgo.MessageEmbed{}}
}

//send
//Sends the embed message
func (e *embed) send(ChannelID string) *discordgo.Message {
	m, err := Session.ChannelMessageSendEmbed(ChannelID, e.MessageEmbed)
	if err != nil {
		Log.Error(err.Error())
	}
	return m
}

//setTitle
//Sets the title of the embed
func (e *embed) setTitle(title string) {
	if len(title) > 256 {
		title = title[:256]
	}
	e.Title = title
}

//setDescription
//Sets the title of the embed
func (e *embed) setDescription(title string) {
	if len(title) > 256 {
		title = title[:256]
	}
	e.Title = title
}

//addField
//Adds a field to a response object
func (e *embed) addField(name string, value string, inline bool) {
	// Cuts value short if it's longer than 1024 characters
	if len(value) > 1024 {
		value = value[:1024]
	}

	//Cuts name short if it's longer than 1024 characters
	if len(name) > 1024 {
		value = value[:1024]
	}

	//Adds the field to the embed object
	e.Fields = append(e.Fields, &discordgo.MessageEmbedField{
		Name:   name,
		Value:  value,
		Inline: inline,
	})
}

//setColor
//Sets the color of an embed message
func (e *embed) setColor(clr int) {
	e.Color = clr
}
