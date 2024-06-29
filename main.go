package main

import (
	"game/gameplay"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
)

func main() {
	game := gameplay.Game{TileSize: 32, Camera: ebiten.NewImage(1024, 1024)}
	// w, h := ebiten.WindowSize()
	// game.SetTileSize(w, h)
	game.Init()
	game.SaveSettings()
	settingsErr := game.SaveSettings()
	game.LoadLevel(1)
	if settingsErr != nil {
		log.Println("[SaveFile]", settingsErr)
	}
	if err := ebiten.RunGame(&game); err != nil {
		log.Fatal(err)
	}
}
