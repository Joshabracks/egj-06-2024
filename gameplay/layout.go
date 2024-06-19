package gameplay

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	g.SetTileSize(outsideWidth, outsideHeight)
	g.ActiveLevel.Render(g.TileDrawSize)
	return g.ScreenWidth, g.ScreenHeight
}

func (g *Game) SetTileSize(outsideWidth, outsideHeight int) {
	g.ScreenWidth = outsideWidth
	g.ScreenHeight = outsideHeight

	min := g.ScreenWidth
	if g.ScreenHeight < min {
		min = g.ScreenHeight
	}

	g.TileDrawSize = min / 32
}