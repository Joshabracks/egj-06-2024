package gameplay

import (
	"bytes"
	"image/png"
	"image/color"
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

var TILE_MAP = map[[4]uint32]uint8 {
	{0, 0, 0, 0}: TILE_EMPTY,
	{8738, 8224, 13364, 65535}: TILE_WALL,
	{65535, 65535, 65535, 65535}: TILE_BODY,
	{44204, 12850, 12850, 65535}: TILE_SPAWN,
}

var TILE_COLOR_MAP = map[uint8]color.Color {
	TILE_EMPTY: color.White,
	TILE_WALL: color.Black,
	TILE_BODY: color.RGBA{150, 150, 150, 255},
	TILE_SPAWN: color.RGBA{255, 0, 0, 255},
}


type Level struct {
	Filepath string
	Map [][]uint8
	MapImage *ebiten.Image
	width, height int
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
	bounds := image.Bounds()
	l.width = bounds.Max.X
	l.height = bounds.Max.Y

	l.Map = make([][]uint8, l.width)
	for x := 0; x < l.width; x++ {
		l.Map[x] = make([]uint8, l.height)
		for y := 0; y < l.height; y++ {
			pixel := image.At(x, y)
			r, g, b, a := pixel.RGBA()
			tileType, ok := TILE_MAP[[4]uint32{r,g,b,a}]
			if !ok {
				tileType = TILE_EMPTY
			}
			l.Map[x][y] = tileType
		}
	}
	return nil
}

func (l *Level) Render(tileSize int) {
	l.MapImage = ebiten.NewImage(l.width * tileSize, l.height * tileSize)
	for x := 0; x < l.width; x++ {
		for y := 0; y < l.height; y++ {
			tileType := l.Map[x][y]
			vector.DrawFilledRect(l.MapImage, float32(x * tileSize), float32(y * tileSize), float32(tileSize), float32(tileSize), TILE_COLOR_MAP[tileType], false)
		}
	}
}