package main

import (
	"game/gameplay"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
)

func main() {
	game := gameplay.Game{TileSize: 32, Camera: ebiten.NewImage(1024, 1024)}
	game.Init()
	game.SaveSettings()
	settingsErr := game.SaveSettings()
	game.LoadLevel(0)
	if settingsErr != nil {
		log.Println("[SaveFile]", settingsErr)
	}
	ebiten.MaximizeWindow()
	if err := ebiten.RunGame(&game); err != nil {
		log.Fatal(err)
	}
}
