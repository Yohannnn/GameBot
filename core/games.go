package core

// Game
// Struct for a game
type Game struct {
	Function func(GameInfo)
	PreviewBoard [][]string
	Guide string
	Rules string
	Color int32
}

// GameInfo
// Information about a currently active game
type GameInfo struct {
	GameBoard [][]string
	PlayerIDs []string
}

// Games
// Map games names to their game struct
var Games map[string]Game