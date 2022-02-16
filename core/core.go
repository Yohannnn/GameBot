package core

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
	"github.com/joho/godotenv"
	tlog "github.com/ubergeek77/tinylog"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"
)

// Session
// The session of the bot
var Session *discordgo.Session

// TODO Rework logging
// log
// The logger for the bot
var log = tlog.NewTaggedLogger("BotCore", tlog.NewColor("38;5;111"))

// graceTerm
// Used to check when the bot is terminating
var graceTerm = true

// botAdmins
// A list of user IDs that are designated as "Bot Administrators"
var botAdmins = make(map[string]bool)

// IsAdmin
// Allow commands to check if a user is an admin or not
func IsAdmin(userId string) bool {
	return botAdmins[userId]
}

// IsCommand
// Check if a given string is a command registered to the core bot
func IsCommand(trigger string) bool {
	if _, ok := Commands[strings.ToLower(trigger)]; ok {
		return true
	}
	return false
}

// TODO Wait for some game functions to finish before terminating
// TODO Repair broken inputs on startup
// TODO Bot admins

// Start
// Loads in token and starts the bot
func Start() {
	// Load token
	err := godotenv.Load("./.env")
	if err != nil {
		log.Error(err.Error())
	}

	// Create a new Discord session
	Session, err = discordgo.New("Bot " + os.Getenv("Token"))
	if err != nil {
		log.Error(err.Error())
	}

	// Add Handlers
	Session.AddHandler(commandHandler)
	Session.AddHandler(reactionHandler)

	// Start discord session
	err = Session.Open()
	if err != nil {
		log.Error(err.Error())
		return
	}

	// Starts go routine for status
	go setStatus()

	log.Info("Bot is now running")

	// -- GRACEFUL TERMINATION -- //

	// Set up a sigterm channel, so we can detect when the application receives a TERM signal
	sigChannel := make(chan os.Signal, 1)
	signal.Notify(sigChannel, syscall.SIGINT, syscall.SIGTERM, syscall.SIGKILL, os.Interrupt, os.Kill)

	// Keep this thread blocked forever, until a TERM signal is received
	<-sigChannel

	log.Info("Received TERM signal, terminating gracefully.")

	// Waits 10 seconds
	time.Sleep(time.Second * 10)

	// Close session when finished
	err = Session.Close()
}

func setStatus() {
	var currentCount int
	for {
		if len(Instances) != currentCount {
			err := Session.UpdateGameStatus(0, fmt.Sprintf("%d Games", len(Instances)))
			if err != nil {
				log.Error(err.Error())
				return
			}
			currentCount = len(Instances)
		}
		time.Sleep(time.Second)
	}
}
