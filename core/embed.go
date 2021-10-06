package core

import "github.com/bwmarrin/discordgo"

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

//sendEmbed
//Sends the embed message
func (e *embed) send(ChannelID string, title string, description string) *discordgo.Message {
	if len(title) > 256 {
		title = title[:256]
	}
	e.Title = title
	if len(description) > 2048 {
		description = description[:2048]
	}
	e.Description = description
	m, err := Session.ChannelMessageSendEmbed(ChannelID, e.MessageEmbed)
	if err != nil {
		Log.Error(err.Error())
	}
	return m
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
