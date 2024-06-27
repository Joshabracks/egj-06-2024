package gameplay

import (
	"game/util"
	"math"

	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

type Character struct {
	CarriedBodyPart *BodyPart
	Boost 			bool
	Stamina         float32
	Image           *ebiten.Image
	Color			color.RGBA
	Object
}

func (c *Character) Render(g *Game) {
	midpoint := float32(g.TileSize) / 2
	x := midpoint
	y := midpoint
	c.Image.Clear()
	vector.DrawFilledCircle(c.Image, x, y, float32(g.TileSize/3), color.Black, true)
	vector.DrawFilledCircle(c.Image, x, y, float32(g.TileSize/4), c.Color, true)
}

func (c *Character) Layout(g *Game) {
	c.Image = ebiten.NewImage(g.TileSize, g.TileSize)
}

func (c *Character) UpdateCarriedItem(g *Game) {
	if c.CarriedBodyPart == nil {
		return
	}
	bodyCoords := g.ActiveLevel.BodyCoordinates
	xIndex, yIndex := Location(c.X, c.Y, g)
	if int(xIndex) == bodyCoords[0] && int(yIndex) == bodyCoords[1] {
		c.CarriedBodyPart.Assembled = true
		c.CarriedBodyPart = nil
		return
	}

	c.CarriedBodyPart.X = c.X
	c.CarriedBodyPart.Y = c.Y
}

func (c *Character) CheckCollisions(g *Game) {
	_, bpCollisions := c.Collisions(g)
	for _, bp := range bpCollisions {
		if !bp.Active || bp.Assembled || c.CarriedBodyPart == bp {
			continue
		}
		c.CarriedBodyPart = bp
	}
}

type PlayerController struct {
	Vertical, Horizontal, Scroll float64
	Character
	KeyBindings KeyBindingMap
}

type KeyBindingMap map[string][]ebiten.Key

func NewPlayerController(g *Game) PlayerController {
	player := Character{
		Object: Object{
			X:     1.5,
			Y:     1.5,
			Speed: 0.1,
		},
		Color: 	color.RGBA{R: 255, G: 255, B: 0, A: 255},
		Boost:   false,
		Stamina: 100,
	}
	player.Layout(g)
	return PlayerController{
		Vertical:   0,
		Horizontal: 0,
		Scroll:     0,
		Character:  player,
		KeyBindings: KeyBindingMap{
			"left":  []ebiten.Key{ebiten.KeyLeft, ebiten.KeyA},
			"right": []ebiten.Key{ebiten.KeyRight, ebiten.KeyD},
			"up":    []ebiten.Key{ebiten.KeyUp, ebiten.KeyW},
			"down":  []ebiten.Key{ebiten.KeyDown, ebiten.KeyS},
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
	for _, binding := range bindings {
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
	x1 := pc.Character.X
	y1 := pc.Character.Y
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
			pc.Character.X = locX
			return
		}
		if !InsideWall(x1, locY, game) {
			pc.Character.Y = locY
			return
		}
		return
	}
	pc.Character.X = locX
	pc.Character.Y = locY
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
