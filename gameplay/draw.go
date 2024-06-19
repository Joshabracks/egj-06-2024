package gameplay

import (
	"github.com/hajimehoshi/ebiten/v2"
)

func (g *Game) Draw(screen *ebiten.Image){
	screen.DrawImage(g.ActiveLevel.MapImage, &ebiten.DrawImageOptions{})
	g.PlayerController.Player.Draw(g, screen)
}
