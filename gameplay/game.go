package gameplay

import (
	"os"
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
	files, err := os.ReadDir("asset/level")
	if err != nil {
		log.Fatal(err)
	}
	file := files[levelCount % len(files)].Name()
	log.Println("Level: ", levelCount, file)
	level := Level{
		Filepath: "asset/level/" + file,
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
