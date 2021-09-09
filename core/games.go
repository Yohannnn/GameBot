package core

import (
	"fmt"
	"strings"
)

//Game
//The function and info of a game
type Game struct {
	StartFunc  func() *GameUpdate
	UpdateFunc func(game *GameUpdate) *GameUpdate
	Info       GameInfo
}

//GameInfo
//Info about a game
type GameInfo struct {
	Name         string
	Description  string
	Rules        string
	Color        int
	ExampleBoard [][]string
}

//GameUpdate
//Info about a game update
type GameUpdate struct {
	GameBoard [][]string
	GameStats map[string]string
	Reactions [][]string
}

//Games
//Map games names to their game struct
var Games = make(map[string]*Game)

//AddGame
//Adds a game to the game map
func AddGame(startFunc func() *GameUpdate, updateFunc func(game *GameUpdate) *GameUpdate, gI *GameInfo) {
	game := Game{
		StartFunc:  startFunc,
		UpdateFunc: updateFunc,
		Info:       *gI,
	}
	Games[strings.ToLower(gI.Name)] = &game
}

//CreateGameInfo
//Creates game info for a command
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
