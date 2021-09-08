package core

import "github.com/bwmarrin/discordgo"

// Embed
// Struct for response
type Embed struct {
	*discordgo.MessageEmbed
}

// NewEmbed
// Returns a new embed object
func NewEmbed() *Embed {
	return &Embed{&discordgo.MessageEmbed{}}
}

//SetTitle
//Sets the title of an embed
func (e *Embed) SetTitle(title string) {
	if len(title) > 256 {
		title = title[:256]
	}
	e.Title = title
}

//SetDescription
//Sets the description of the embed
func (e *Embed) SetDescription(description string) {
	if len(description) > 2048 {
		description = description[:2048]
	}
	e.Description = description
}

// AddField
// Adds a field to a response object
func (e *Embed) AddField(name string, value string, inline bool) {
	// Cuts value short if it's longer than 1024 characters
	if len(value) > 1024 {
		value = value[:1024]
	}

	// Cuts name short if it's longer than 1024 characters
	if len(name) > 1024 {
		value = value[:1024]
	}

	// Adds the field to the embed object
	e.Fields = append(e.Fields, &discordgo.MessageEmbedField{
		Name:   name,
		Value:  value,
		Inline: inline,
	})
}

//SetColor ...
func (e *Embed) SetColor(clr int) {
	e.Color = clr
}
