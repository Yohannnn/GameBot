package core

import (
	"errors"
	"github.com/bwmarrin/discordgo"
	"regexp"
	"strings"
)

// Colors
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

// Squares
const (
	RedSqr    = "🟥"
	BrownSqr  = "🟫"
	OrangeSqr = "🟧"
	YellowSqr = "🟨"
	GreenSqr  = "🟩"
	BlueSqr   = "🟦"
	PurpleSqr = "🟪"
	WhiteSqr  = "⬜"
	BlackSqr  = "⬛"
)

// NumCord
// A map of numbers to their emoji (except 0 cause its sus)
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

// LetCord
// A map of a letters place in the alphabet to its emoji
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

// embed
// Struct for response
type embed struct {
	*discordgo.MessageEmbed
}

// newEmbed
// Returns a new embed object
func newEmbed() *embed {
	return &embed{&discordgo.MessageEmbed{}}
}

// send
// Sends the embed message
func (e *embed) send(title string, description string, channelID string) (*discordgo.Message, error) {
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
		return nil, err
	}
	return m, nil
}

// edit
// Edits an embed message
func (e *embed) edit(title string, description string, channelID string, messageID string) *discordgo.Message {
	if len(title) > 256 {
		title = title[:256]
	}
	e.Title = title

	if len(description) > 256 {
		description = description[:256]
	}
	e.Description = description
	m, err := Session.ChannelMessageEditEmbed(channelID, messageID, e.MessageEmbed)
	if err != nil {
		log.Error(err.Error())
	}
	return m
}

// addField
// Adds a field to a response object
func (e *embed) addField(name string, value string, inline bool) {
	//  Cuts value short if it's longer than 1024 characters
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

// setFooter
// Sets the footer of an embed
func (e *embed) setFooter(text string, iconURL string, proxyIconURL string) {
	e.Footer = &discordgo.MessageEmbedFooter{
		Text:         text,
		IconURL:      iconURL,
		ProxyIconURL: proxyIconURL,
	}
}

// setColor
// Sets the color of an embed message
func (e *embed) setColor(clr int) {
	e.Color = clr
}

// Contains
// Checks if an array contains an element
func Contains(s []string, item string) bool {
	for i := 0; i < len(s); i++ {
		if s[i] == item {
			return true
		}
	}
	return false
}

// formatBoard
// Formats a game board into a string
func formatBoard(board [][]string) string {
	var BoardString string
	var LineString string

	for _, l := range board {
		for _, e := range l {
			if strings.Contains(e, ":") {
				emoji, err := Session.State.Emoji("806048328973549578", e)
				if err != nil {
					log.Errorf("Error when formatting board: %s", err.Error())
					return ""
				}
				LineString += emoji.MessageFormat() + "‌"
			} else {
				LineString += e + "‌"
			}
		}
		BoardString += LineString + "\n"
		LineString = ""
	}
	return BoardString
}

// RemoveItems
// Remove items from a slice by value
func RemoveItems(slice []string, deleteables []string) []string {
	var newSlice []string
	for _, elem := range slice {
		if !Contains(deleteables, elem) {
			newSlice = append(newSlice, elem)
		}
	}
	return newSlice
}

// ensureNumbers
// Given a string, ensure it contains only numbers
// This is useful for stripping letters and formatting characters from user/role pings
func ensureNumbers(in string) string {
	reg, err := regexp.Compile("[^0-9]+")
	if err != nil {
		log.Errorf("Error when ensuring numbers: %s", err)
		return ""
	}

	return reg.ReplaceAllString(in, "")
}

// cleanId
// Given a string, attempt to remove all numbers from it
// Additionally, ensure it is at least 17 characters in length
// This is a way of "cleaning" a Discord ping into a valid snowflake string
func cleanId(in string) string {
	out := ensureNumbers(in)

	//  Discord IDs must be, at minimum, 17 characters long
	if len(out) < 17 {
		return ""
	}

	return out
}

// getUser
// Given a user ID, get that user's object (global to Discord, not in a guild)
func getUser(userId string) (*discordgo.User, error) {
	cleanedId := cleanId(userId)
	if cleanedId == "" {
		return nil, errors.New("provided ID is invalid")
	}

	return Session.User(cleanedId)
}

// MakeCordOption
// Takes 2 integers (max value for each cord) and returns an array of corresponding emojis
func MakeCordOption(Xmax int, Ymax int) []string {
	var Options []string

	for i := 1; i <= Xmax; i++ {
		Options = append(Options, NumCord[i])
	}

	for i := 1; i <= Ymax; i++ {
		Options = append(Options, LetCord[i])
	}

	return Options
}
