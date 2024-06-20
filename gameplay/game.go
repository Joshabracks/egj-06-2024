package gameplay

type Game struct {
	ActiveLevel  Level
	TileDrawSize, TileImageSize int
	PlayerController
	Settings
}
