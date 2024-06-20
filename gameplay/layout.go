package gameplay


func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	// g.SetTileSize(outsideWidth, outsideHeight)
	g.ScreenWidth = outsideWidth
	g.ScreenHeight = outsideHeight

	min := g.ScreenWidth
	if g.ScreenHeight < min {
		min = g.ScreenHeight
	}
	g.Scale = float64(min / g.TileSize) / float64(g.TileSize)
	// cameraSize := min * g.TileSize
	// g.Camera = ebiten.NewImage(cameraSize, cameraSize)
	// g.ActiveLevel.Render(g)
	return g.ScreenWidth, g.ScreenHeight
}

// func (g *Game) SetTileSize(outsideWidth, outsideHeight int) {
// 	g.ScreenWidth = outsideWidth
// 	g.ScreenHeight = outsideHeight

// 	min := g.ScreenWidth
// 	if g.ScreenHeight < min {
// 		min = g.ScreenHeight
// 	}

// 	// g.TileSize = min / g.TileImageSize
// 	// g.PlayerController.Player.Layout(g)
// }
