package gameplay

import (
	"bytes"
	"image/color"
	"image/png"
	"os"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

const (
	TILE_EMPTY uint8 = iota
	TILE_WALL
	TILE_BODY
	TILE_SPAWN
)

var TILE_MAP = map[[4]uint32]uint8{
	{0, 0, 0, 0}:                 TILE_EMPTY,
	{8738, 8224, 13364, 65535}:   TILE_WALL,
	{65535, 65535, 65535, 65535}: TILE_BODY,
	{44204, 12850, 12850, 65535}: TILE_SPAWN,
}

var TILE_COLOR_MAP = map[uint8]color.Color{
	TILE_EMPTY: color.White,
	TILE_WALL:  color.Black,
	TILE_BODY:  color.RGBA{150, 150, 150, 255},
	TILE_SPAWN: color.RGBA{255, 0, 0, 255},
}

type Level struct {
	Filepath                                          string
	Map                                               [][]uint8
	BodyCoordinates                                   [2]int
	SpawnerCoordinates                                [][2]int
	MapImage                                          *ebiten.Image
	width, height                                     int
	Head, Torso, ArmRight, ArmLeft, LegRight, LegLeft BodyPart
	Objects []Object
}

func (l *Level) Load() error {
	data, err := os.ReadFile(l.Filepath)
	if err != nil {
		return err
	}
	buffer := bytes.NewBuffer(data)
	image, err := png.Decode(buffer)
	if err != nil {
		return err
	}
	hasBody := false
	bounds := image.Bounds()
	l.width = bounds.Max.X
	l.height = bounds.Max.Y
	l.SpawnerCoordinates = make([][2]int, 0)
	l.Map = make([][]uint8, l.width)
	for x := 0; x < l.width; x++ {
		l.Map[x] = make([]uint8, l.height)
		for y := 0; y < l.height; y++ {
			pixel := image.At(x, y)
			r, g, b, a := pixel.RGBA()
			tileType, ok := TILE_MAP[[4]uint32{r, g, b, a}]
			if !ok || (tileType == TILE_BODY && hasBody) {
				tileType = TILE_EMPTY
			}
			if tileType == TILE_BODY && !hasBody {
				l.BodyCoordinates = [2]int{x, y}
				hasBody = true
			}
			if tileType == TILE_SPAWN {
				l.SpawnerCoordinates = append(l.SpawnerCoordinates, [2]int{x, y})
			}
			l.Map[x][y] = tileType
		}
	}
	return nil
}

func (b *BodyPart) SetActive(g *Game) bool {
	if !b.Active {
		b.Activate(g)
	}
	return b.Assembled
}

func (l *Level) SetActiveBodyPart(g *Game) {
	if !l.Head.SetActive(g) {
		return
	}
	if !l.Torso.SetActive(g) {
		return
	}
	if !l.ArmLeft.SetActive(g) {
		return
	}
	if !l.ArmRight.SetActive(g) {
		return
	}
	if !l.LegLeft.SetActive(g) {
		return
	}
	l.LegRight.SetActive(g)

}

func (l *Level) LoadParts(g *Game) {
	l.Head = NewBodyPart(HEAD, GREY_BOI, g)
	l.Torso = NewBodyPart(TORSO, GREY_BOI, g)
	l.ArmLeft = NewBodyPart(ARM_LEFT, GREY_BOI, g)
	l.ArmRight = NewBodyPart(ARM_RIGHT, GREY_BOI, g)
	l.LegLeft = NewBodyPart(LEG_LEFT, GREY_BOI, g)
	l.LegRight = NewBodyPart(LEG_RIGHT, GREY_BOI, g)
}

func (l *Level) Render(g *Game) {
	l.MapImage = ebiten.NewImage(l.width*g.TileSize, l.height*g.TileSize)
	for x := 0; x < l.width; x++ {
		for y := 0; y < l.height; y++ {
			tileType := l.Map[x][y]
			vector.DrawFilledRect(l.MapImage, float32(x*g.TileSize), float32(y*g.TileSize), float32(g.TileSize), float32(g.TileSize), TILE_COLOR_MAP[tileType], false)
		}
	}
}

func (l *Level) RenderOverlay(g *Game, screen *ebiten.Image) {
	stamina := g.Player.Stamina
	if stamina > 0 {
		stamina = stamina / 100
	}
	vector.DrawFilledRect(screen, float32(g.TileSize), float32(g.TileSize/4), float32(g.TileSize*8)*stamina, float32(g.TileSize/2), color.White, false)

	l.Torso.Draw(g, screen)
	l.ArmLeft.Draw(g, screen)
	l.ArmRight.Draw(g, screen)
	l.LegLeft.Draw(g, screen)
	l.LegRight.Draw(g, screen)
	l.Head.Draw(g, screen)
}
