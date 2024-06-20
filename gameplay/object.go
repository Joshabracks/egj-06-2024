package gameplay

import (
	"game/util"
	"image"

	"github.com/hajimehoshi/ebiten/v2"
)

type Object struct {
	X, Y, Speed float64
	Direction   float64
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

var CREATURE_IMAGE_PATH = map[CreatureType]string {
	GREY_BOI: "asset/image/body/grey-boi.png",
	GREY_GURL: "asset/image/body/grey-gurl.png",
	SKINNY_BOI: "asset/image/body/skinny-boi.png",
	SKINNY_GURL: "asset/image/body/skinny-gurl.png",
	CHUBBY_BOI: "asset/image/body/chubby-boi.png",
	CHUBBY_GURL: "asset/image/body/chubby-gurl.png",
	AVERAGE_BOI: "asset/image/body/average-boi.png",
	AVERAGE_GURL: "asset/image/body/average-gurl.png",
	ATHLETIC_BOI: "asset/image/body/athletic-boi.png",
	ATHLETIC_GURL: "asset/image/body/athletic-gurl.png",
	STRONK_BOI: "asset/image/body/stronk-boi.png",
	STRONK_GURL: "asset/image/body/stronk-gurl.png",
}

var CREATURE_IMAGE = map[CreatureType]*ebiten.Image {}

type BodyPart struct {
	Object
	BodyPartType
	CreatureType
	Image *ebiten.Image
}

var BODY_PART_UV_INDICES = map[BodyPartType][4]int{
	HEAD:      {1, 0, 2, 2},
	TORSO:     {1, 2, 2, 2},
	ARM_LEFT:  {0, 0, 1, 2},
	ARM_RIGHT: {3, 0, 1, 2},
	LEG_LEFT:  {0, 2, 1, 2},
	LEG_RIGHT: {3, 2, 1, 2},
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
	sliceSize := g.TileImageSize / 4
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
