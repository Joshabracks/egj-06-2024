package gameplay

import (
	"github.com/hajimehoshi/ebiten/v2"
)

func (g *Game) Init() {
	g.InitSettings()
	width, height := ebiten.WindowSize()
	ebiten.SetWindowSize(width, height)

	g.PlayerController = NewPlayerController(g)
	ebiten.SetWindowTitle("Body Builder")
	ebiten.SetWindowResizingMode(ebiten.WindowResizingModeEnabled)
}
