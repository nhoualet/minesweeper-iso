package rendering

import (
	"fmt"

	"github.com/veandco/go-sdl2/sdl"
	"github.com/veandco/go-sdl2/ttf"
)

type Textbox struct {
	texture *sdl.Texture
	Rect    sdl.Rect
	Text    string
	color   sdl.Color
	pos     sdl.Point
	centerX bool
	centerY bool
}

func NewTextbox(rect sdl.Rect, centerX bool, centerY bool, text string, renderer *sdl.Renderer, font *ttf.Font, color sdl.Color) (*Textbox, error) {
	tbox := &Textbox{}
	if text != "" {
		surface, err := font.RenderUTF8Solid(text, color)
		if err != nil {
			return nil, fmt.Errorf("texture surface error : %s", err)
		}
		defer surface.Free()
		tbox.texture, err = renderer.CreateTextureFromSurface(surface)
		if err != nil {
			return nil, fmt.Errorf("texture texture error : %s", err)
		}
		_, _, tbox.Rect.W, tbox.Rect.H, _ = tbox.texture.Query()
	} else {
		tbox.texture = nil
	}
	tbox.Text = text
	tbox.color = color
	tbox.pos = sdl.Point{X: rect.X, Y: rect.Y}
	tbox.centerX = centerX
	tbox.centerY = centerY
	tbox.updatePos()
	return tbox, nil
}

func (tbox *Textbox) SetText(text string, renderer *sdl.Renderer, font *ttf.Font, color sdl.Color) error {
	if text == tbox.Text {
		return nil
	}
	if text != "" {
		// surface, err := font.RenderUTF8BlendedWrapped(text, color, 10)
		surface, err := font.RenderUTF8Solid(text, color)
		if err != nil {
			return fmt.Errorf("texture surface error : %s", err)
		}
		defer surface.Free()
		newTexture, err := renderer.CreateTextureFromSurface(surface)
		if err != nil {
			return fmt.Errorf("texture texture error : %s", err)
		}
		if tbox.texture != nil {
			tbox.texture.Destroy()
		}
		tbox.texture = newTexture
		tbox.updatePos()
	} else if tbox.texture != nil {
		tbox.texture.Destroy()
		tbox.texture = nil
	}
	tbox.Text = text
	return nil
}

func (tbox *Textbox) SetCenter(x int32, y int32) {
	tbox.pos.X = x
	tbox.pos.Y = y
	tbox.centerX = true
	tbox.centerY = true
	tbox.updatePos()
}

func (tbox *Textbox) SetTopLeft(x int32, y int32) {
	tbox.pos.X = x
	tbox.pos.Y = y
	tbox.centerX = false
	tbox.centerY = false
	tbox.updatePos()
}

func (tbox *Textbox) SetX(x int32, isCenter bool) {
	tbox.pos.X = x
	tbox.centerX = isCenter
	tbox.updatePos()
}

func (tbox *Textbox) SetY(y int32, isCenter bool) {
	tbox.pos.Y = y
	tbox.centerY = isCenter
	tbox.updatePos()
}

func (tbox *Textbox) Width() int32 {
	return tbox.Rect.W
}
func (tbox *Textbox) Height() int32 {
	return tbox.Rect.H
}

func (tbox *Textbox) SetBackgroundSize(w int32, h int32) {
	tbox.Rect.W = w
	tbox.Rect.H = h
	tbox.updatePos()
}

func (tbox *Textbox) updatePos() {
	if tbox.texture != nil {
		_, _, tbox.Rect.W, tbox.Rect.H, _ = tbox.texture.Query()

		if tbox.centerX {
			tbox.Rect.X = tbox.pos.X - tbox.Rect.W/2
		} else {
			tbox.Rect.X = tbox.pos.X
		}
		if tbox.centerY {
			tbox.Rect.Y = tbox.pos.Y - tbox.Rect.H/2
		} else {
			tbox.Rect.Y = tbox.pos.Y
		}
	}
}

func (tbox *Textbox) Destroy() {
	if tbox.texture != nil {
		tbox.texture.Destroy()
	}
}

func (tbox *Textbox) Draw(r *CustomRenderer) {
	if tbox.texture != nil {
		r.SDLrenderer.Copy(tbox.texture, &sdl.Rect{X: 0, Y: 0, W: tbox.Rect.W, H: tbox.Rect.H}, &tbox.Rect)
	}
}
