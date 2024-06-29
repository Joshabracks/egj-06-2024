package gameplay

import (
	"image/color"
	"math"
	"time"

	"github.com/fzipp/astar"
)

type EnemyController struct {
	Path                 astar.Path[Vector]
	Horizontal, Vertical float64
	Character
}

var enemyStartSpeed float64 = 0.05

func NewEnemy(x, y float64, g *Game) *EnemyController {
	character := Character{
		Object: Object{
			X:     x,
			Y:     y,
			Speed: enemyStartSpeed,
		},
		Boost:   false,
		Stamina: 100,
		Color:   color.RGBA{R: 255, G: 0, B: 0, A: 255},
	}
	character.Layout(g)
	return &EnemyController{Character: character}
}

func (ec *EnemyController) Move(g *Game) {
	targetX, targetY := Location(g.PlayerController.X, g.PlayerController.Y, g)
	x, y := Location(ec.X, ec.Y, g)
	clearPath := false

	if targetX == x {
		currentX := int(x)
		incr := 1
		if targetX < x {
			incr = -1
		}
		clearPath = true
		for currentX != int(targetX) {
			if currentX < 0 || currentX >= g.ActiveLevel.width {
				clearPath = false
				break
			}
			row := g.ActiveLevel.Map[currentX]

			cell := row[int(y)]
			if cell == TILE_WALL {
				clearPath = false
				break
			}
			currentX += incr
		}
		// if clearPath {
		// 	ec.Horizontal = float64(incr)
		// }
	} else if targetY == y {
		if targetY == y {
			currentY := int(y)
			incr := 1
			if targetY < y {
				incr = -1
			}
			clearPath = true
			for currentY != int(targetY) {
				row := g.ActiveLevel.Map[int(x)]
				if currentY < 0 || currentY >= g.ActiveLevel.height {
					clearPath = false
					break
				}
				cell := row[int(currentY)]
				if cell == TILE_WALL {
					clearPath = false
					break
				}
				currentY += incr
			}
			// if clearPath {
			// 	ec.Vertical = float64(incr)
			// }
		}
	}
	if clearPath { // enemy chases player
		
		xDiff := math.Abs(ec.X - g.PlayerController.X)
		yDiff := math.Abs(ec.Y - g.PlayerController.Y)
		distance := math.Sqrt(yDiff * yDiff + xDiff * xDiff)
		normX := (ec.X - g.PlayerController.X) / distance
		normY := (ec.Y - g.PlayerController.Y) / distance
		ec.X = ec.X - (normX * ec.Speed * 2)
		ec.Y = ec.Y - (normY * ec.Speed * 2)
		return
	}
	if ec.Horizontal == 0 && ec.Vertical == 0 {
		ec.Horizontal = 1
		return
	}
	nextX := ec.X + ec.Horizontal
	nextY := ec.Y + ec.Vertical
	if nextX < 0 || nextX >= float64(g.ActiveLevel.width) {
		ec.Horizontal *= -1
		return
	}
	if nextY < 0 || nextY >= float64(g.ActiveLevel.height) {
		nextY *= -1
		return
	}
	nextTile := g.ActiveLevel.Map[int(nextX)][int(nextY)]
	if nextTile == TILE_WALL {
		r := time.Now().UnixMilli() % 4
		switch r {
		case 0:
			ec.Horizontal = 1
			ec.Vertical = 0
		case 1:
			ec.Horizontal = -1
			ec.Vertical = 0
		case 2:
			ec.Horizontal = 0
			ec.Vertical = 1
		case 3:
			ec.Horizontal = 0
			ec.Vertical = -1
		}
		return
	}

	speed := ec.Speed
	if ec.Boost && ec.Stamina > 0 {
		speed *= 2
		ec.Stamina -= 0.5
		if ec.Stamina < 0 {
			ec.Stamina = 0
		}
	} else if !ec.Boost && ec.Stamina < 100 {
		ec.Stamina += 0.05
		if ec.Stamina > 100 {
			ec.Stamina = 100
		}
	}
	x1 := ec.X
	y1 := ec.Y
	x2 := x1 + ec.Horizontal
	y2 := y1 + ec.Vertical
	xDiff := x2 - x1
	yDiff := y2 - y1
	ec.Direction = math.Atan2(float64(yDiff), float64(xDiff)) * 180 / math.Pi
	xDist := math.Cos(ec.Direction*math.Pi/180) * float64(speed)
	yDist := math.Sin(ec.Direction*math.Pi/180) * float64(speed)
	locX := x1 + xDist
	locY := y1 + yDist
	if InsideWall(locX, locY, g) {
		if !InsideWall(locX, y1, g) {
			ec.Character.X = locX
			return
		}
		if !InsideWall(x1, locY, g) {
			ec.Character.Y = locY
			return
		}
		return
	}
	ec.Character.X = locX
	ec.Character.Y = locY
}
