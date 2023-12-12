package rendering

import (
	"errors"
	"fmt"

	"github.com/veandco/go-sdl2/sdl"
	"github.com/veandco/go-sdl2/ttf"
)

type ButtonActionId int

type Button struct {
	Rect            sdl.Rect
	innerRect       sdl.Rect
	texture         *sdl.Texture
	TextureRect     sdl.Rect
	Text            string
	color           sdl.Color
	backgroundColor sdl.Color
	hoverColor      sdl.Color
	hovered         bool
	selected        bool
	selectable      bool
	ActionId        ButtonActionId
	pos             sdl.Point
	centerX         bool
	centerY         bool
	drawBackground  bool
	drawHover       bool
	Shown           bool
	useTextureSize  bool
}

func DummyButton() *Button {
	b := Button{
		Rect:            sdl.Rect{X: 0, Y: 0, W: 10, H: 10},
		Text:            "DummyButton",
		color:           sdl.Color{R: 200, G: 0, B: 0, A: sdl.ALPHA_OPAQUE},
		drawBackground:  true,
		backgroundColor: sdl.Color{R: 200, G: 0, B: 200, A: sdl.ALPHA_OPAQUE},
		hoverColor:      sdl.Color{R: 200, G: 200, B: 0, A: sdl.ALPHA_OPAQUE},
		drawHover:       true,
		selected:        false,
		hovered:         false,
		selectable:      false,
		pos:             sdl.Point{X: 0, Y: 0},
		centerX:         false,
		centerY:         false,
		ActionId:        0,
		Shown:           true,
		useTextureSize:  true,
	}
	return &b
}

func NewButton(rect sdl.Rect, centerX bool, centerY bool, selectable bool, text string, actionId ButtonActionId, color sdl.Color, backgroundColor *sdl.Color, hoverColor *sdl.Color) *Button {
	b := Button{
		Rect:           rect,
		Text:           text,
		color:          color,
		drawBackground: backgroundColor != nil,
		drawHover:      hoverColor != nil,
		selected:       false,
		hovered:        false,
		selectable:     selectable,
		pos:            sdl.Point{X: rect.X, Y: rect.Y},
		centerX:        centerX,
		centerY:        centerY,
		ActionId:       actionId,
		Shown:          true,
		useTextureSize: true,
	}
	if backgroundColor != nil {
		b.backgroundColor = *backgroundColor
	}
	if hoverColor != nil {
		b.hoverColor = *hoverColor
	}
	return &b
}

func (b *Button) SetHovered(state bool) {
	b.hovered = state
}

func (b *Button) Selected(value bool) {
	if b.selectable {
		b.selected = value
	}
}

func (b *Button) SetSelectable(value bool) {
	b.selectable = value
	if !b.selectable {
		b.selected = false
	}
}

func (b *Button) IsSelectable() bool {
	return b.selectable
}

func (b *Button) UpdateTexture(renderer *sdl.Renderer, font *ttf.Font) error {
	if b.Text != "" {
		surface, err := font.RenderUTF8Solid(b.Text, b.color)
		if err != nil {
			return fmt.Errorf("texture surface error : %s", err)
		}
		defer surface.Free()
		newTexture, err := renderer.CreateTextureFromSurface(surface)
		if err != nil {
			return fmt.Errorf("texture texture error : %s", err)
		}
		if b.texture != nil {
			b.texture.Destroy()
		}
		b.texture = newTexture
		b.updatePos()
	} else {
		b.texture = nil
	}
	return nil
}
func (b *Button) SetText(text string, renderer *sdl.Renderer, font *ttf.Font, color sdl.Color) error {
	b.Text = text
	b.color = color
	b.UpdateTexture(renderer, font)
	return nil
}

func (b *Button) SetTexture(texture *sdl.Texture) error {
	if texture == nil {
		return errors.New("setTexture: texture is nil")
	}
	if b.texture != nil {
		b.texture.Destroy()
	}
	b.texture = texture
	return nil
}

func (b *Button) SetCenter(x int32, y int32) {
	b.pos.X = x
	b.pos.Y = y
	b.centerX = true
	b.centerY = true
	b.updatePos()
}

func (b *Button) SetTopLeft(x int32, y int32) {
	b.pos.X = x
	b.pos.Y = y
	b.centerX = false
	b.centerY = false
	b.updatePos()
}

func (b *Button) SetX(x int32, isCenter bool) {
	b.pos.X = x
	b.centerX = isCenter
	b.updatePos()
}

func (b *Button) SetY(y int32, isCenter bool) {
	b.pos.Y = y
	b.centerY = isCenter
	b.updatePos()
}

func (b *Button) Width() int32 {
	return b.Rect.W
}
func (b *Button) Height() int32 {
	return b.Rect.H
}

func (b *Button) IsHovered() bool {
	return b.hovered
}

func (b *Button) IsSelected() bool {
	return b.selected
}

func (b *Button) SetBackgroundSize(w int32, h int32) {
	b.Rect.W = w
	b.Rect.H = h
	b.updatePos()
}

func (b *Button) SetTextureSize(w int32, h int32) {
	b.innerRect.W = w
	b.innerRect.H = h
	b.useTextureSize = false
	b.updatePos()
}
func (b *Button) SetHoverColor(color *sdl.Color) {
	if color == nil {
		b.drawHover = false
	} else {
		b.hoverColor = *color
		b.drawHover = true
	}
}

func (b *Button) SetAction(id ButtonActionId) {
	b.ActionId = id
}
func (b *Button) updatePos() {

	if b.texture != nil && b.useTextureSize {
		_, _, b.TextureRect.W, b.TextureRect.H, _ = b.texture.Query()
		b.innerRect.W = b.TextureRect.W
		b.innerRect.H = b.TextureRect.H
		// b.Rect.W = b.TextureRect.W
		// b.Rect.H = b.TextureRect.H
	}
	if b.centerX {
		b.Rect.X = b.pos.X - b.Rect.W/2
		b.innerRect.X = b.pos.X - b.innerRect.W/2
	} else {
		b.Rect.X = b.pos.X
		b.innerRect.X = b.pos.X + (b.Rect.W-b.innerRect.W)/2
	}
	if b.centerY {
		b.Rect.Y = b.pos.Y - b.Rect.H/2
		b.innerRect.Y = b.pos.Y - b.innerRect.H/2
	} else {
		b.Rect.Y = b.pos.Y
		b.innerRect.Y = b.pos.Y + (b.Rect.H-b.innerRect.H)/2
	}
}

func (b *Button) OnButton(pos sdl.Point) bool {
	if !b.Shown {
		return false
	}
	return pos.InRect(&b.Rect)
}

func (b *Button) Destroy() {
	if b.texture != nil {
		b.texture.Destroy()
	}
}

func (b *Button) Draw(r *CustomRenderer) {
	if !b.Shown {
		return
	}
	if b.drawBackground {
		r.SetDrawColor(b.backgroundColor)
		r.SDLrenderer.FillRect(&b.Rect)
	}
	if b.texture != nil {
		r.SDLrenderer.Copy(b.texture, &sdl.Rect{X: 0, Y: 0, W: b.TextureRect.W, H: b.TextureRect.H}, &b.innerRect)
	}
	if b.drawHover && b.hovered {
		r.SetDrawColor(b.hoverColor)
		r.SDLrenderer.DrawRect(&b.Rect)
	} else if b.selected {
		r.SetDrawColor(b.color)
		r.SDLrenderer.DrawRect(&b.Rect)
	}
}
