package gameplay

import (
	"github.com/hajimehoshi/ebiten/v2"
)

func (g *Game) Init() {
	textInit()
	g.InitSettings()
	width, height := ebiten.WindowSize()
	ebiten.SetWindowSize(width, height)

	ebiten.SetWindowTitle("Body Builder")
	ebiten.SetWindowResizingMode(ebiten.WindowResizingModeEnabled)
}
