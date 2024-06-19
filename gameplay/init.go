package gameplay

import (
	"github.com/hajimehoshi/ebiten/v2"
)

func (g *Game) Init() {
	g.InitSettings()
	width, height := ebiten.WindowSize()
	ebiten.SetWindowSize(width, height)
	player := Player{X: 1.5, Y: 1.5, Speed: 0.1}
	g.PlayerController = NewPlayerController(player)
	ebiten.SetWindowTitle("Body Builder")
	ebiten.SetWindowResizingMode(ebiten.WindowResizingModeEnabled)
}
