package gameplay

import (
	"bytes"
	"fmt"
	"image/color"
	"image/png"
	"log"
	"math"
	"os"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/examples/resources/fonts"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
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
	Count int
	Filepath                                          string
	Map                                               [][]uint8
	BodyCoordinates                                   [2]int
	SpawnerCoordinates                                [][2]int
	MapImage                                          *ebiten.Image
	width, height                                     int
	Head, Torso, ArmRight, ArmLeft, LegRight, LegLeft BodyPart
	// Objects                                           []Object
	Enemies []*EnemyController
	graph[Vector]
}

func (l *Level) PopulateEnemies(g *Game, count int) {
	for _, coords := range l.SpawnerCoordinates {
		if len(l.Enemies) > count {
			break
		}
		enemy := NewEnemy(float64(coords[0])+0.5, float64(coords[1])+0.5, g)
		l.Enemies = append(l.Enemies, enemy)
	}
}

// func (l *Level) ActiveBodyPartLocation(g *Game) (float64, float64) {
// 	if l.Head.Active && !l.Head.Assembled {
// 		return Location(l.Head.X, l.Head.Y, g)
// 	}
// 	if l.Torso.Active && !l.Torso.Assembled {
// 		return Location(l.Torso.X, l.Torso.Y, g)
// 	}
// 	if l.ArmRight.Active && !l.ArmRight.Assembled {
// 		return Location(l.ArmRight.X, l.ArmRight.Y, g)
// 	}
// 	if l.ArmLeft.Active && !l.ArmLeft.Assembled {
// 		return Location(l.ArmLeft.X, l.ArmLeft.Y, g)
// 	}
// 	if l.LegRight.Active && !l.LegRight.Assembled {
// 		return Location(l.LegRight.X, l.LegRight.Y, g)
// 	}
// 	if l.LegLeft.Active && !l.LegLeft.Assembled {
// 		return Location(l.LegLeft.X, l.LegLeft.Y, g)
// 	}
// 	return -1, -1
// }

func (l *Level) Load(levelCount int) error {
	l.Count = levelCount
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
	if !l.LegRight.SetActive(g) {
		return
	}
	g.LoadLevel(g.ActiveLevel.Count + 1)
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

var mplusFaceSource *text.GoTextFaceSource
func textInit() {
	s, err := text.NewGoTextFaceSource(bytes.NewReader(fonts.MPlus1pRegular_ttf))
	if err != nil {
		log.Fatal(err)
	}
	mplusFaceSource = s
}

func (l *Level) RenderOverlay(g *Game, screen *ebiten.Image) {
	stamina := g.Character.Stamina
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
	op := &text.DrawOptions{}
	op.ColorScale.ScaleWithColor(color.White)
	op.GeoM.Translate(float64(g.TileSize) * 26, 0)
	textOverlay := fmt.Sprint("Level: ", g.ActiveLevel.Count + 1)
	text.Draw(screen, textOverlay, &text.GoTextFace{Source: mplusFaceSource, Size: 24}, op)
}

func nodeDist(p, q Vector) float64 {
	dx := p.X - q.X
	dy := p.Y - q.Y
	return math.Sqrt(dx*dx + dy*dy)
}

type graph[Node comparable] map[Node][]Node
type Vector struct {
	X, Y float64
}

func NewVector(x, y float64) Vector {
	return Vector{X: x, Y: y}
}
func newGraph[Node comparable]() graph[Node] {
	return make(graph[Node])
}

func (g graph[Node]) link(a, b Node) graph[Node] {
	g[a] = append(g[a], b)
	g[b] = append(g[b], a)
	return g
}

func (g graph[Node]) Neighbours(n Node) []Node {
	return g[n]
}

func (l *Level) InitGraph() {
	l.graph = newGraph[Vector]()
	for x, row := range l.Map {
		for y, cell := range row {
			if cell != TILE_WALL {
				vector := Vector{X: float64(x), Y: float64(y)}
				l.graph[vector] = []Vector{
					NewVector(vector.X - 1, vector.Y - 1), NewVector(vector.X, vector.Y - 1), NewVector(vector.X + 1, vector.Y - 1),
					NewVector(vector.X - 1, vector.Y), NewVector(vector.X + 1, vector.Y),
					NewVector(vector.X - 1, vector.Y + 1), NewVector(vector.X, vector.Y + 1), NewVector(vector.X+ 1, vector.Y + 1),
				}
			}
		}
	}
}
