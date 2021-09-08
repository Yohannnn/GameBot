package core

// Game
// Struct for a game
type Game struct {
	Name         string
	Function     func(GameInfo)
	PreviewBoard [][]string
	Description  string
	Rules        string
	Color        int
}

//GameInfo
//Contains information about a game
type GameInfo struct {
	GameBoard [][]string
	GameStats map[string]string
}

// Games
// Map games names to their game struct
var Games = make(map[string]Game)

// AddGame
// Adds a game to the game map
func AddGame(name string, game Game) {
	Games[name] = game
}
