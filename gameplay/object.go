package gameplay

import (
	"game/util"
	"image"
	"math"
	"math/rand"

	"github.com/hajimehoshi/ebiten/v2"
)

type Object struct {
	X, Y, Speed float64
	Direction   float64
	Active      bool
}

func Location(x, y float64, g *Game) (float64, float64) {
	xIndex := x - math.Mod(float64(x), 1)
	yIndex := y - math.Mod(float64(y), 1)
	return xIndex, yIndex
}

func (o *Object) Collisions(g *Game) ([]*Object, []*BodyPart) {
	objects := make([]*Object, 0)
	bodyParts := make([]*BodyPart, 0)
	if Collision(o, &g.ActiveLevel.Head.Object, g) {
		bodyParts = append(bodyParts, &g.ActiveLevel.Head)
	}
	if Collision(o, &g.ActiveLevel.Torso.Object, g) {
		bodyParts = append(bodyParts, &g.ActiveLevel.Torso)
	}
	if Collision(o, &g.ActiveLevel.ArmLeft.Object, g) {
		bodyParts = append(bodyParts, &g.ActiveLevel.ArmLeft)
	}
	if Collision(o, &g.ActiveLevel.ArmRight.Object, g) {
		bodyParts = append(bodyParts, &g.ActiveLevel.ArmRight)
	}
	if Collision(o, &g.ActiveLevel.LegLeft.Object, g) {
		bodyParts = append(bodyParts, &g.ActiveLevel.LegLeft)
	}
	if Collision(o, &g.ActiveLevel.LegRight.Object, g) {
		bodyParts = append(bodyParts, &g.ActiveLevel.LegRight)
	}
	for _, object := range(g.ActiveLevel.Enemies) {
		xDist := math.Abs(object.X - o.X)
		yDist := math.Abs(object.Y - o.Y)
		if xDist + yDist < 0.25 {
			objects = append(objects, &object.Object)
		}
	}
	return objects, bodyParts
}

func Collision(a, b *Object, g *Game) bool {
	xDist := math.Abs(a.X - b.X)
	yDist := math.Abs(a.Y - b.Y)
	return xDist < 0.25 &&  yDist < 0.25
}

type BodyPartType uint8

const (
	HEAD BodyPartType = iota
	TORSO
	ARM_RIGHT
	ARM_LEFT
	LEG_RIGHT
	LEG_LEFT
)

var BODY_PART_UV_INDICES = map[BodyPartType][4]int{
	HEAD:      {1, 0, 2, 2},
	TORSO:     {1, 2, 2, 2},
	ARM_LEFT:  {0, 0, 1, 2},
	ARM_RIGHT: {3, 0, 1, 2},
	LEG_LEFT:  {0, 2, 1, 2},
	LEG_RIGHT: {3, 2, 1, 2},
}

var BODY_PART_TABLE_OFFSETS = map[BodyPartType][2]int{
	HEAD:      {8, 0},
	TORSO:     {8, 14},
	ARM_LEFT:  {0, 14},
	ARM_RIGHT: {24, 14},
	LEG_LEFT:  {8, 30},
	LEG_RIGHT: {16, 30},
}

type CreatureType uint8

const (
	GREY_BOI CreatureType = iota
	GREY_GURL
	SKINNY_BOI
	SKINNY_GURL
	CHUBBY_BOI
	CHUBBY_GURL
	AVERAGE_BOI
	AVERAGE_GURL
	ATHLETIC_BOI
	ATHLETIC_GURL
	STRONK_BOI
	STRONK_GURL
)

var CREATURE_IMAGE_PATH = map[CreatureType]string{
	GREY_BOI:      "asset/image/body/grey-boi.png",
	GREY_GURL:     "asset/image/body/grey-gurl.png",
	SKINNY_BOI:    "asset/image/body/skinny-boi.png",
	SKINNY_GURL:   "asset/image/body/skinny-gurl.png",
	CHUBBY_BOI:    "asset/image/body/chubby-boi.png",
	CHUBBY_GURL:   "asset/image/body/chubby-gurl.png",
	AVERAGE_BOI:   "asset/image/body/average-boi.png",
	AVERAGE_GURL:  "asset/image/body/average-gurl.png",
	ATHLETIC_BOI:  "asset/image/body/athletic-boi.png",
	ATHLETIC_GURL: "asset/image/body/athletic-gurl.png",
	STRONK_BOI:    "asset/image/body/stronk-boi.png",
	STRONK_GURL:   "asset/image/body/stronk-gurl.png",
}

var CREATURE_IMAGE = map[CreatureType]*ebiten.Image{}

type BodyPart struct {
	Object
	BodyPartType
	CreatureType
	Assembled bool
	Image     *ebiten.Image
}

func (b *BodyPart) Activate(g *Game) {
	b.Active = true
	r := rand.Intn(len(g.ActiveLevel.SpawnerCoordinates))
	coords := g.ActiveLevel.SpawnerCoordinates[r]
	b.X = float64(coords[0]) + 0.5
	b.Y = float64(coords[1]) + 0.5
}	

func (b *BodyPart) Draw(g *Game, screen *ebiten.Image) {
	if !b.Active && !b.Assembled {
		b.DrawIcon(g, screen)
		return
	}
	if b.Assembled {
		b.DrawAssembled(g, screen)
		return
	}
	if b.Active {
		b.DrawActive(g, screen)
	}
}

func (b *BodyPart) DrawActive(g *Game, screen *ebiten.Image) {
	bounds := b.Image.Bounds()

	xOffset := (b.X * float64(g.TileSize)) - float64(bounds.Dx() / 2)
	yOffset := (b.Y * float64(g.TileSize)) - float64(bounds.Dy() / 2)
	op := ebiten.DrawImageOptions{}
	op.GeoM.Translate(float64(xOffset), float64(yOffset))
	screen.DrawImage(b.Image, &op)
}

func (b *BodyPart) DrawAssembled(g *Game, screen *ebiten.Image) {
	xOffset := (g.ActiveLevel.BodyCoordinates[0] * g.TileSize) + (BODY_PART_TABLE_OFFSETS[b.BodyPartType][0])
	yOffset := (g.ActiveLevel.BodyCoordinates[1] * g.TileSize) + (BODY_PART_TABLE_OFFSETS[b.BodyPartType][1])
	op := ebiten.DrawImageOptions{}
	op.GeoM.Translate(float64(xOffset), float64(yOffset))
	screen.DrawImage(b.Image, &op)
}

func (b *BodyPart) DrawIcon(g *Game, screen *ebiten.Image) {
	op := ebiten.DrawImageOptions{}
	x := float64((g.TileSize + 5) * (1 + int(b.BodyPartType)))
	y := float64(g.TileSize * 31)
	op.GeoM.Translate(x, y)
	screen.DrawImage(b.Image, &op)
}

func NewBodyPart(bodyPartType BodyPartType, creatureType CreatureType, g *Game) BodyPart {
	bodyPart := BodyPart{
		BodyPartType: bodyPartType,
		CreatureType: creatureType,
		Assembled:    false,
		Object: Object{
			X:      -1,
			Y:      -1,
			Active: false,
		},
	}
	bodyPart.LoadImage(g)
	return bodyPart
}

func (b *BodyPart) LoadImage(g *Game) {
	filepath := CREATURE_IMAGE_PATH[b.CreatureType]
	creatureImage, ok := CREATURE_IMAGE[b.CreatureType]
	if !ok {
		creatureImage = util.LoadImage(filepath)
	}
	if creatureImage == nil {
		return
	}
	if !ok {
		CREATURE_IMAGE[b.CreatureType] = creatureImage
	}
	uvIndices := BODY_PART_UV_INDICES[b.BodyPartType]
	sliceSize := g.TileSize / 4
	x1 := uvIndices[0] * sliceSize
	y1 := uvIndices[1] * sliceSize
	width := uvIndices[2] * sliceSize
	height := uvIndices[3] * sliceSize
	x2 := x1 + width
	y2 := y1 + height
	rect := image.Rect(x1, y1, x2, y2)
	image := creatureImage.SubImage(rect).(*ebiten.Image)
	b.Image = image
}
