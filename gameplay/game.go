package gameplay

type Game struct {
	ActiveLevel  Level
	TileDrawSize int
	PlayerController
	Settings
}
