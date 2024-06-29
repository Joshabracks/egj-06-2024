package gameplay

import (
	"fmt"
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

func pauseScreen(g *Game, screen *ebiten.Image) {
	a := make([]string, 0)
	if g.ActiveLevel.Count == 0 {
		a = append(a, []string{
			"You are the Body Builder",
			"",
			"There's a loose body part somewhere on the stage.",
			"Move to the body part to pick it up and move to the grey square",
			"to add it to the body.",
			"If an enemy drains your energy completely, it's game over.",
			"Sprinting draings your energy, but won't kill you",
			"",
			"Controls:",
			"---Move-- WASD",
			"--Sprint- Spacebar",
			"--Pause-- Enter",
			"",
			"",}...)
	}
	a = append(a, []string{
		fmt.Sprintf("Level %d", g.ActiveLevel.Count + 1),
		"",
		"",
		"Press ENTER to continue",
	}...)
	
	screen.Clear()
	for i, line := range a {
		op := &text.DrawOptions{}
		op.ColorScale.ScaleWithColor(color.White)
		op.GeoM.Translate(float64(g.TileSize), float64(g.TileSize*i))
		text.Draw(screen, line, &text.GoTextFace{Source: mplusFaceSource, Size: 24}, op)
	}
}

func (g *Game) Draw(screen *ebiten.Image) {
	if g.ActiveLevel.Pause {
		pauseScreen(g, screen)
		return
	}
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
