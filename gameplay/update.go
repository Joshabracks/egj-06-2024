package gameplay

func (g *Game) Update() error {
	g.ActiveLevel.SetActiveBodyPart(g)
	g.PlayerController.UpdateInput()
	g.PlayerController.UpdatePlayerPosition(g)
	g.PlayerController.Character.CheckCollisions(g)
	g.PlayerController.Character.UpdateCarriedItem(g)
	for _, enemy := range(g.ActiveLevel.Enemies) {
		enemy.Move(g)
	}
	return nil
}
