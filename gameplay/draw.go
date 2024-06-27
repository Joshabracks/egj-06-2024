package gameplay

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

func (g *Game) Draw(screen *ebiten.Image) {
	g.Camera.DrawImage(g.ActiveLevel.MapImage, &ebiten.DrawImageOptions{})
	op := &ebiten.DrawImageOptions{}
	g.PlayerController.Render(g)
	op.GeoM.Translate(
		(g.PlayerController.X*float64(g.TileSize))-float64(g.TileSize/2),
		(g.PlayerController.Y*float64(g.TileSize))-float64(g.TileSize/2))
	g.Camera.DrawImage(g.Character.Image, op)
	for _, enemy := range(g.ActiveLevel.Enemies) {
		enemy.Render(g)
		op = &ebiten.DrawImageOptions{}
		op.GeoM.Translate(
			(enemy.X*float64(g.TileSize))-float64(g.TileSize/2),
			(enemy.Y*float64(g.TileSize))-float64(g.TileSize/2))
		g.Camera.DrawImage(enemy.Image, op)
		for _, n := range(enemy.Path) {
			vector.DrawFilledCircle(screen, float32(n.X * float64(g.TileSize)), float32(n.Y * float64(g.TileSize)), float32(g.TileSize/4), color.Black, true)
		}
	}
	op = &ebiten.DrawImageOptions{}
	g.ActiveLevel.RenderOverlay(g, g.Camera)
	op.GeoM.Scale(g.Scale, g.Scale)
	screen.DrawImage(g.Camera, op)
}
