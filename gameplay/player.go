package gameplay

import (
	"game/util"
	"math"

	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

type Player struct {
	Item uint8
	Boost bool
	Stamina float32
	Image *ebiten.Image
	Object
}

func (p *Player) Render(g *Game) {
	midpoint := float32(g.TileSize) / 2
	x := midpoint
	y := midpoint
	p.Image.Clear()
	vector.DrawFilledCircle(p.Image, x, y, float32(g.TileSize / 3), color.Black, true)
	vector.DrawFilledCircle(p.Image, x, y, float32(g.TileSize / 4), color.RGBA{R: 255, G: 255, B: 0, A: 255}, true)
}

func (p *Player) Layout(g *Game) {
	p.Image = ebiten.NewImage(g.TileSize, g.TileSize)
}

type PlayerController struct {
	Vertical, Horizontal, Scroll float64
	Player
	KeyBindings KeyBindingMap
}

type KeyBindingMap map[string][]ebiten.Key

func NewPlayerController(g *Game) PlayerController {
	player := Player{
		Object: Object{
			X:     1.5,
			Y:     1.5,
			Speed: 0.1,
		},
		Boost:   false,
		Stamina: 100,
	}
	player.Layout(g)
	return PlayerController{
		Vertical: 0,
		Horizontal: 0,
		Scroll: 0,
		Player: player,
		KeyBindings: KeyBindingMap{
			"left": []ebiten.Key{ebiten.KeyLeft, ebiten.KeyA},
			"right": []ebiten.Key{ebiten.KeyRight, ebiten.KeyD},
			"up": []ebiten.Key{ebiten.KeyUp, ebiten.KeyW},
			"down": []ebiten.Key{ebiten.KeyDown, ebiten.KeyS},
			"boost": []ebiten.Key{ebiten.KeySpace, ebiten.KeyEnter},
		},
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

func (pc *PlayerController) IsKeyPressed(key string) bool {
	bindings, ok := pc.KeyBindings[key]
	if !ok {
		return false
	}
	for _, binding := range(bindings) {
		if ebiten.IsKeyPressed(binding) {
			return true
		}
	}
	return false
}

func (pc *PlayerController) UpdateInput() {
	pc.Reset()
	if pc.IsKeyPressed("up") {
		pc.Vertical--
	}
	if pc.IsKeyPressed("down") {
		pc.Vertical++
	}
	if pc.IsKeyPressed("left") {
		pc.Horizontal--
	}
	if pc.IsKeyPressed("right") {
		pc.Horizontal++
	}
	if pc.IsKeyPressed("boost") {
		pc.Boost = true
	} else {
		pc.Boost = false
	}
	pc.Truncate()
}

func (pc *PlayerController) UpdatePlayerPosition(game *Game) {
	speed := pc.Speed
	if pc.Boost && pc.Stamina > 0 {
		speed *= 2
		pc.Stamina -= 0.5
		if pc.Stamina < 0 {
			pc.Stamina = 0
		}
	} else if !pc.Boost && pc.Stamina < 100 {
		pc.Stamina += 0.05
		if pc.Stamina > 100 {
			pc.Stamina = 100
		}
	}
	if pc.Horizontal == 0 && pc.Vertical == 0 {
		return
	}
	x1 := pc.Player.X
	y1 := pc.Player.Y
	x2 := x1 + pc.Horizontal
	y2 := y1 + pc.Vertical
	xDiff := x2 - x1
	yDiff := y2 - y1
	pc.Direction = math.Atan2(float64(yDiff), float64(xDiff)) * 180 / math.Pi
	xDist := math.Cos(pc.Direction*math.Pi/180) * float64(speed)
	yDist := math.Sin(pc.Direction*math.Pi/180) * float64(speed)
	locX := x1 + xDist
	locY := y1 + yDist
	if InsideWall(locX, locY, game) {
		if !InsideWall(locX, y1, game) {
			pc.Player.X = locX
			return
		}
		if !InsideWall(x1, locY, game) {
			pc.Player.Y = locY
			return
		}
		return
	}
	pc.Player.X = locX
	pc.Player.Y = locY
}

func InsideWall(x, y float64, game *Game) bool {
	xIndex := int(x - math.Mod(float64(x), 1))
	yIndex := int(y - math.Mod(float64(y), 1))
	if xIndex < 0 || xIndex >= game.ActiveLevel.width || yIndex < 0 || yIndex >= game.ActiveLevel.height {
		return true
	}
	tileType := game.ActiveLevel.Map[xIndex][yIndex]
	return tileType == TILE_WALL
}
