package main

import (
	"log"
	"github.com/hajimehoshi/ebiten/v2"
	"game/gameplay"
)

func main() {
	game := gameplay.Game{TileSize: 32, Camera: ebiten.NewImage(1024, 1024)}
	// w, h := ebiten.WindowSize()
	// game.SetTileSize(w, h)
	game.Init()
	game.SaveSettings()
	settingsErr := game.SaveSettings()
	testLevel := gameplay.Level{
		Filepath: "asset/level/map_01.png",
	}
	
	testLevelErr := testLevel.Load()
	if testLevelErr != nil {
		log.Println("[TEST LEVEL ERROR]", testLevelErr)
	}
	testLevel.LoadParts(&game)
	testLevel.Render(&game)
	game.ActiveLevel = testLevel
	if settingsErr != nil {
		log.Println("[SaveFile]", settingsErr)
	}
	if err := ebiten.RunGame(&game); err != nil {
		log.Fatal(err)
	}
}