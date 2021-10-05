package core

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
	"github.com/joho/godotenv"
	tlog "github.com/ubergeek77/tinylog"
	"os"
	"os/signal"
	"syscall"
)

//Session
//The session of the bot
var Session *discordgo.Session

//Log
//The logger for the bot
var Log = tlog.NewTaggedLogger("BotCore", tlog.NewColor("38;5;111"))

func Start() error {
	//Load token
	err := godotenv.Load("./.env")
	if err != nil {
		fmt.Println(err)
	}

	//Create a new Discord session
	Session, err = discordgo.New("Bot " + os.Getenv("Token"))
	if err != nil {
		return err
	}

	//Add Handlers
	Session.AddHandler(commandHandler)
	Session.AddHandler(reactionHandler)

	//Star discord session
	err = Session.Open()
	if err != nil {
		return err
	}

	Log.Info("Bot is now running.  Press CTRL-C to exit.")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc

	//Close session when finished
	err = Session.Close()
	return err
}
