package gameplay

import (
	"github.com/hajimehoshi/ebiten/v2"
)

func (g *Game) Draw(screen *ebiten.Image){
	screen.DrawImage(g.ActiveLevel.MapImage, &ebiten.DrawImageOptions{})
	op := &ebiten.DrawImageOptions{}
	g.PlayerController.Player.Render(g)
	op.GeoM.Translate(
		(g.Player.X * float64(g.TileDrawSize)) - float64(g.TileDrawSize / 2), 
		(g.Player.Y * float64(g.TileDrawSize)) - float64(g.TileDrawSize / 2))
	g.ActiveLevel.RenderOverlay(g, screen)
	screen.DrawImage(g.Player.Image, op)
}
