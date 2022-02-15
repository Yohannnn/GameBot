package core

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
	"github.com/joho/godotenv"
	tlog "github.com/ubergeek77/tinylog"
	"os"
	"os/signal"
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

// TODO Wait for some functions to finish before shutting down

// Start
// Starts the bot
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

	// Sets status
	go setStatus()

	log.Info("Bot is now running")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc

	// Close session when finished
	err = Session.Close()
}

func setStatus() {
	for {
		err := Session.UpdateGameStatus(0, fmt.Sprintf("%d Games", len(Instances)))
		if err != nil {
			log.Error(err.Error())
			return
		}
		time.Sleep(time.Minute)
	}
}
