package gameplay

import (
	"github.com/hajimehoshi/ebiten/v2"
)

func (g *Game) Draw(screen *ebiten.Image) {
	g.Camera.DrawImage(g.ActiveLevel.MapImage, &ebiten.DrawImageOptions{})
	op := &ebiten.DrawImageOptions{}
	g.PlayerController.Player.Render(g)
	op.GeoM.Translate(
		(g.Player.X*float64(g.TileSize))-float64(g.TileSize/2),
		(g.Player.Y*float64(g.TileSize))-float64(g.TileSize/2))
	g.Camera.DrawImage(g.Player.Image, op)
	op = &ebiten.DrawImageOptions{}
	g.ActiveLevel.RenderOverlay(g, g.Camera)
	op.GeoM.Scale(g.Scale, g.Scale)
	screen.DrawImage(g.Camera, op)
}
