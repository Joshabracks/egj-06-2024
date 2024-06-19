package gameplay

import (
	"game/util"
	"math"

	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

type Player struct {
	X, Y,
	Speed float32
	Item uint8
}

func (p *Player) Draw(g *Game, screen *ebiten.Image) {
	x := p.X * float32(g.TileDrawSize)
	y := p.Y * float32(g.TileDrawSize)
	vector.DrawFilledCircle(screen, x, y, float32(g.TileDrawSize / 3), color.Black, true)
	vector.DrawFilledCircle(screen, x, y, float32(g.TileDrawSize / 4), color.RGBA{R: 255, G: 255, B: 0, A: 255}, true)
}

type PlayerController struct {
	Vertical, Horizontal, Scroll float32
	Player
	Left  []ebiten.Key
	Right []ebiten.Key
	Up    []ebiten.Key
	Down  []ebiten.Key
}

func NewPlayerController(p Player) PlayerController {
	return PlayerController{
		Vertical: 0,
		Horizontal: 0,
		Scroll: 0,
		Player: p,
		Left: []ebiten.Key{ebiten.KeyLeft, ebiten.KeyA},
		Right: []ebiten.Key{ebiten.KeyRight, ebiten.KeyD},
		Up: []ebiten.Key{ebiten.KeyUp, ebiten.KeyW},
		Down: []ebiten.Key{ebiten.KeyDown, ebiten.KeyS},
	}
}

func (pc *PlayerController) Truncate() {
	pc.Horizontal = util.Clamp(pc.Horizontal, -1, 1)
	pc.Vertical = util.Clamp(pc.Vertical, -1, 1)
}

func (pc *PlayerController) Reset() {
	pc.Vertical = 0
	pc.Horizontal = 0
	pc.Scroll = 0
}

func IsKeyPressed(keys []ebiten.Key) bool {
	for _, key := range(keys) {
		if ebiten.IsKeyPressed(key) {
			return true
		}
	}
	return false
}

func (pc *PlayerController) UpdateInput() {
	pc.Reset()
	if IsKeyPressed(pc.Up) {
		pc.Vertical--
	}
	if IsKeyPressed(pc.Down) {
		pc.Vertical++
	}
	if IsKeyPressed(pc.Left) {
		pc.Horizontal--
	}
	if IsKeyPressed(pc.Right) {
		pc.Horizontal++
	}
	pc.Truncate()
}

func (pc *PlayerController) UpdatePlayerPosition(game *Game) {
	if pc.Horizontal == 0 && pc.Vertical == 0 {
		return
	}
	x1 := pc.Player.X
	y1 := pc.Player.Y
	x2 := x1 + pc.Horizontal
	y2 := y1 + pc.Vertical
	xDiff := x2 - x1
	yDiff := y2 - y1
	dir := math.Atan2(float64(yDiff), float64(xDiff)) * 180 / math.Pi
	xDist := math.Cos(dir*math.Pi/180) * float64(pc.Speed)
	yDist := math.Sin(dir*math.Pi/180) * float64(pc.Speed)
	locX := x1 + float32(xDist)
	locY := y1 + float32(yDist)
	if InsideWall(locX, locY, game) {
		if !InsideWall(locX, float32(y1), game) {
			pc.Player.X = locX
			return
		}
		if !InsideWall(float32(x1), locY, game) {
			pc.Player.Y = locY
			return
		}
		return
	}
	pc.Player.X = locX
	pc.Player.Y = locY
}

func InsideWall(x, y float32, game *Game) bool {
	xIndex := int(float64(x) - math.Mod(float64(x), 1))
	yIndex := int(float64(y) - math.Mod(float64(y), 1))
	if xIndex < 0 || xIndex >= game.ActiveLevel.width || yIndex < 0 || yIndex >= game.ActiveLevel.height {
		return true
	}
	tileType := game.ActiveLevel.Map[xIndex][yIndex]
	return tileType == TILE_WALL
}
