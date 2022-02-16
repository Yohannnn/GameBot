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

// log
// The logger for the bot
var log = tlog.NewTaggedLogger("BotCore", tlog.NewColor("38;5;111"))

// graceTerm
// Used to check when the bot is terminating
var graceTerm bool

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

// TODO Repair broken inputs on startup
// TODO Bot admins

// Start
// Loads in token and starts the bot
func Start() {
	// Load token
	err := godotenv.Load("./.env")
	if err != nil {
		log.Panicf("Error when starting bot: %s", err.Error())
		return
	}

	// Create a new Discord session
	Session, err = discordgo.New("Bot " + os.Getenv("Token"))
	if err != nil {
		log.Panicf("Error when starting bot: %s", err.Error())
		return
	}

	// Add Handlers
	Session.AddHandler(commandHandler)
	Session.AddHandler(reactionHandler)

	// Start discord session
	err = Session.Open()
	if err != nil {
		log.Panicf("Error when starting bot: %s", err.Error())
		return
	}

	// Checks for broken inputs and repairs them
	log.Info("Checking for and repairing broken inputs")
	for _, inst := range Instances {
		currentMessage, err := Session.ChannelMessage(inst.Players[inst.Turn].ChannelID, inst.CurrentMessageID)
		if err != nil {
			log.Errorf("Error when repairing inputs: %s", err.Error())
			return
		}
		for _, r := range currentMessage.Reactions {
			if !Contains(inst.CurrentInput.Options, r.Emoji.ID) && !Contains(inst.CurrentInput.Options, r.Emoji.Name) && r.Emoji.Name != "✅" {
				err = addInput(inst.CurrentInput, inst.Players[inst.Turn].ChannelID, inst.CurrentMessageID)
				if err != nil {
					log.Errorf("Error when repairing inputs: %s", err.Error())
					return
				}
				log.Infof("Repaired input for %s", inst.ID)
				break
			}
		}
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
	graceTerm = true

	// Make a second sig channel that will respond to user term signal immediately
	sigInstant := make(chan os.Signal, 1)
	signal.Notify(sigInstant, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)

	// Waits for reaction and command handler to finish
	go func() {
		log.Info("Waiting for reaction handler to finish")
		reactionLock.Lock()
		log.Info("Reaction handler stopped gracefully")
		log.Info("Waiting for command handler to finish")
		commandLock.Lock()
		log.Info("Command handler stopped gracefully")
		// Send our own signal to the instant sig channel
		sigInstant <- syscall.SIGTERM
	}()

	// Keep the thread blocked until the above goroutine finishes closing all handlers, or until another TERM is received
	<-sigInstant

	// Close session when finished
	err = Session.Close()
}

func setStatus() {
	var currentCount int
	for !graceTerm {
		if len(Instances) != currentCount {
			err := Session.UpdateGameStatus(0, fmt.Sprintf("%d Games", len(Instances)))
			if err != nil {
				log.Errorf("Error when setting status: %s", err.Error())
				return
			}
			currentCount = len(Instances)
		}
		time.Sleep(time.Second)
	}
}
