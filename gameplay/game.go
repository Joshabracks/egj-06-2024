package gameplay

import (
	"fmt"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
)

type Game struct {
	ActiveLevel Level
	TileSize    int
	PlayerController
	Settings
	Camera *ebiten.Image
	Scale  float64
}

func (g *Game) LoadLevel(levelCount int) {
	g.PlayerController = NewPlayerController(g)
	filepath := fmt.Sprintf("asset/level/map_%d.png", 1)
	log.Println("Level: ", levelCount)
	level := Level{
		Filepath: filepath,
	}

	levelErr := level.Load(levelCount)
	if levelErr != nil {
		log.Println(levelErr)
	}
	level.LoadParts(g)
	level.PopulateEnemies(g, levelCount + 2)
	level.InitGraph()
	level.Render(g)
	g.ActiveLevel = level
}
