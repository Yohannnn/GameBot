package core

import (
	"fmt"
	"strings"
)

//Game
//The update function, start function, and game information for a game
type Game struct {
	UpdateFunc func(*GameUpdate)
	StartFunc  func() GameState
	Info       GameInfo
}

//GameState
//The state of a game
type GameState struct {
	GameBoard [][]string
	GameStats map[string]string
}

//GameInfo
//Information about a game
type GameInfo struct {
	Name         string
	Description  string
	Rules        string
	Color        int
	ExampleBoard [][]string
}

//GameUpdate
//Information about a game update
type GameUpdate struct {
	PlayerState  GameState
	OpponentSate GameState
	OptionType   string
	Option       []string
}

//Games
//Map games names to their game struct
var Games = make(map[string]*Game)

//AddGame
//Adds a game to the game map
func AddGame(updateFunc func(game *GameUpdate), gI *GameInfo) {
	game := Game{
		UpdateFunc: updateFunc,
		Info:       *gI,
	}
	Games[strings.ToLower(gI.Name)] = &game
}

//CreateGameInfo
//Creates game info for a game
func CreateGameInfo(name string, description string, rules string, color int, exampleBoard [][]string) *GameInfo {
	gI := &GameInfo{
		Name:         name,
		Description:  description,
		Rules:        rules,
		ExampleBoard: exampleBoard,
		Color:        color,
	}
	return gI
}

//formatBoard
//Formats a game board into a string
func formatBoard(board [][]string) string {
	var BoardString string
	var LineString string
	for _, l := range board {
		for _, e := range l {
			LineString += fmt.Sprintf(":%s:", e)
		}
		BoardString += LineString + "\n"
		LineString = ""
	}
	return BoardString
}
