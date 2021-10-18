package core

import (
	"github.com/bwmarrin/discordgo"
	"reflect"
)

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
func (e *embed) send(title string, description string, channelID string) *discordgo.Message {
	if len(title) > 256 {
		title = title[:256]
	}
	e.Title = title

	if len(description) > 256 {
		description = description[:256]
	}
	e.Description = description
	m, err := Session.ChannelMessageSendEmbed(channelID, e.MessageEmbed)
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

//setFooter
//Sets the footer of an embed
func (e *embed) setFooter(text string, iconURL string, proxyIconURL string) {
	e.Footer = &discordgo.MessageEmbedFooter{
		Text:         text,
		IconURL:      iconURL,
		ProxyIconURL: proxyIconURL,
	}
}

//setColor
//Sets the color of an embed message
func (e *embed) setColor(clr int) {
	e.Color = clr
}

//Contains
//Checks if an array contains an element
func Contains(arrayType interface{}, item interface{}) bool {
	arr := reflect.ValueOf(arrayType)

	if arr.Kind() != reflect.Array {
		panic("Invalid data-type")
	}

	for i := 0; i < arr.Len(); i++ {
		if arr.Index(i).Interface() == item {
			return true
		}
	}

	return false
}

//IntArray
//Returns an array that starts at one value and end at another
func IntArray(x, y int) []int {
	var a []int
	z := x
	for z <= y {
		a = append(a, z)
		z++
	}
	return a
}
