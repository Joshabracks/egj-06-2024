package gameplay

import "github.com/hajimehoshi/ebiten/v2"

type Game struct {
	ActiveLevel Level
	TileSize    int
	PlayerController
	Settings
	Camera *ebiten.Image
	Scale float64
}
