package gameplay

func (g *Game) Update() error {
	g.ActiveLevel.SetActiveBodyPart(g)
	g.PlayerController.UpdateInput()
	g.PlayerController.UpdatePlayerPosition(g)
	g.PlayerController.Player.CheckCollisions(g)
	g.PlayerController.Player.UpdateCarriedItem(g)
	return nil
}
