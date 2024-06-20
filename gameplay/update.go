package gameplay

func (g *Game) Update() error {
	g.PlayerController.UpdateInput()
	g.PlayerController.UpdatePlayerPosition(g)
	return nil
}
