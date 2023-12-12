package rendering

import (
	"errors"
	"fmt"

	"github.com/veandco/go-sdl2/img"
	"github.com/veandco/go-sdl2/sdl"
)

type Spritesheet struct {
	spriteRect        sdl.Rect
	currentSpriteRect sdl.Rect
	currentSpriteId   uint32
	columnCount       uint32
	rowCount          uint32
	texture           *sdl.Texture
}

func NewSpritesheet(renderer *CustomRenderer, filePath string, columnCount, rowCount int32) (*Spritesheet, error) {
	surf, err := img.Load(filePath)
	if err != nil {
		return nil, fmt.Errorf("surface creation error: %s", err)
	}
	defer surf.Free()
	texture, err := renderer.SDLrenderer.CreateTextureFromSurface(surf)
	if err != nil {
		return nil, fmt.Errorf("texture creation error: %s", err)
	}
	_, _, w, h, _ := texture.Query()
	spriteRect := sdl.Rect{
		X: 0,
		Y: 0,
		W: w / columnCount,
		H: h / rowCount,
	}
	return &Spritesheet{
		spriteRect:        spriteRect,
		columnCount:       uint32(columnCount),
		rowCount:          uint32(rowCount),
		texture:           texture,
		currentSpriteRect: spriteRect,
		currentSpriteId:   0,
	}, nil
}

func (s *Spritesheet) SelectSprite(id uint32) error {
	if id == s.currentSpriteId {
		return nil
	}
	if id > s.columnCount*s.rowCount {
		return errors.New("invalid id")
	}
	s.currentSpriteId = id
	s.updateCurrentSpriteRect()
	return nil
}

func (s *Spritesheet) updateCurrentSpriteRect() {
	column := s.currentSpriteId % s.columnCount
	row := (s.currentSpriteId - column) / s.columnCount
	s.currentSpriteRect.X = int32(column) * s.spriteRect.W
	s.currentSpriteRect.Y = int32(row) * s.spriteRect.H
}

func (s *Spritesheet) Draw(r *CustomRenderer, dest sdl.Rect) {
	r.SDLrenderer.Copy(s.texture, &s.currentSpriteRect, &dest)
}
