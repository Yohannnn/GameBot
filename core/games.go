package core

// Game
// Struct for a game
type Game struct {
	Name          string
	Function      func(update GameUpdate)
	StartingBoard [][]string
	Description   string
	Rules         string
	Color         int
}

//GameUpdate
//Contains information a game update
type GameUpdate struct {
	GameBoard [][]string
	GameStats map[string]string
	Reactions [][]string
}

// Games
// Map games names to their game struct
var Games = make(map[string]Game)

// AddGame
// Adds a game to the game map
func AddGame(name string, game Game) {
	Games[name] = game
}
