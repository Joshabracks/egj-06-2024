package gameplay

import (
	"bytes"
	"image/color"
	"image/png"
	"log"
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

func (l *Level) LoadParts(g *Game) {
	l.Head = BodyPart{BodyPartType: HEAD, CreatureType: GREY_BOI}
	l.Torso = BodyPart{BodyPartType: TORSO, CreatureType: GREY_BOI}
	l.ArmLeft = BodyPart{BodyPartType: ARM_LEFT, CreatureType: GREY_BOI}
	l.ArmRight = BodyPart{BodyPartType: ARM_RIGHT, CreatureType: GREY_BOI}
	l.LegLeft = BodyPart{BodyPartType: LEG_LEFT, CreatureType: GREY_BOI}
	l.LegRight = BodyPart{BodyPartType: LEG_RIGHT, CreatureType: GREY_BOI}

	l.Head.LoadImage(g)
	l.Torso.LoadImage(g)
	l.ArmLeft.LoadImage(g)
	l.ArmRight.LoadImage(g)
	l.LegLeft.LoadImage(g)
	l.LegRight.LoadImage(g)
}

func (l *Level) Render(g *Game) {
	l.MapImage = ebiten.NewImage(l.width*g.TileDrawSize, l.height*g.TileDrawSize)
	for x := 0; x < l.width; x++ {
		for y := 0; y < l.height; y++ {
			tileType := l.Map[x][y]
			vector.DrawFilledRect(l.MapImage, float32(x*g.TileDrawSize), float32(y*g.TileDrawSize), float32(g.TileDrawSize), float32(g.TileDrawSize), TILE_COLOR_MAP[tileType], false)
		}
	}
	op := ebiten.DrawImageOptions{}
	tileWarp := float64(g.TileDrawSize) / float64(g.TileImageSize)
	op.GeoM.Scale(tileWarp, tileWarp)
	op.GeoM.Translate(float64(g.TileDrawSize), float64(g.TileDrawSize) * 31)
	l.MapImage.DrawImage(l.Head.Image, &op)
	op.GeoM.Translate(float64(g.TileDrawSize), 0)
	l.MapImage.DrawImage(l.Torso.Image, &op)
	op.GeoM.Translate(float64(g.TileDrawSize), 0)
	l.MapImage.DrawImage(l.ArmLeft.Image, &op)
	op.GeoM.Translate(float64(g.TileDrawSize), 0)
	l.MapImage.DrawImage(l.ArmRight.Image, &op)
	op.GeoM.Translate(float64(g.TileDrawSize), 0)
	l.MapImage.DrawImage(l.LegLeft.Image, &op)
	op.GeoM.Translate(float64(g.TileDrawSize), 0)
	l.MapImage.DrawImage(l.LegRight.Image, &op)
}

func (l *Level) RenderOverlay(g *Game, screen *ebiten.Image) {
	stamina := g.Player.Stamina
	if stamina > 0 {
		stamina = stamina / 100
	}
	vector.DrawFilledRect(screen, float32(g.TileDrawSize), float32(g.TileDrawSize/4), float32(g.TileDrawSize*8)*stamina, float32(g.TileDrawSize/2), color.White, false)
}
