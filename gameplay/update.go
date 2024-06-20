package gameplay

func (g *Game) Update() error {
	g.ActiveLevel.SetActiveBodyPart(g)
	g.PlayerController.UpdateInput()
	g.PlayerController.UpdatePlayerPosition(g)
	return nil
}
