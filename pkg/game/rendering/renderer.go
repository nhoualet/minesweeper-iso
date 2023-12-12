package rendering

import (
	"fmt"

	"github.com/veandco/go-sdl2/img"
	"github.com/veandco/go-sdl2/sdl"
	"github.com/veandco/go-sdl2/ttf"
)

type CustomRenderer struct {
	SDLwindow    *sdl.Window
	SDLrenderer  *sdl.Renderer
	isFullscreen bool
	isBordered   bool
}

func (r *CustomRenderer) IsFullscreen() bool {
	return r.isFullscreen
}

func (r *CustomRenderer) IsBordered() bool {
	return r.isBordered
}
func CreateCustomRenderer(title string, x int32, y int32, w int32, h int32, windowflags uint32, rendererIndex int, rendererFlags int) (*CustomRenderer, error) {
	window, err := sdl.CreateWindow(title, x, y, w, h, windowflags)
	if err != nil {
		return nil, fmt.Errorf("window init %s", err)
	}

	renderer, err := sdl.CreateRenderer(window, rendererIndex, uint32(rendererFlags))
	if err != nil {
		window.Destroy()
		return nil, fmt.Errorf("render init: %s", err)
	}
	return &CustomRenderer{
		SDLwindow:    window,
		SDLrenderer:  renderer,
		isFullscreen: (windowflags&sdl.WINDOW_FULLSCREEN != 0 || windowflags&sdl.WINDOW_FULLSCREEN_DESKTOP != 0),
		isBordered:   (windowflags&sdl.WINDOW_BORDERLESS == 0),
	}, nil
}

func (r *CustomRenderer) Destroy() {
	if r.SDLrenderer != nil {
		r.SDLrenderer.Destroy()
	}
	if r.SDLwindow != nil {
		r.SDLwindow.Destroy()
	}
}

func (r *CustomRenderer) SetDrawColor(color sdl.Color) {
	r.SDLrenderer.SetDrawColor(color.R, color.G, color.B, color.A)
}

func (r *CustomRenderer) NewTextTexture(text string, font *ttf.Font, color sdl.Color) (*sdl.Texture, error) {
	if text != "" {
		surface, err := font.RenderUTF8Blended(text, color)
		if err != nil {
			return nil, fmt.Errorf("texture surface error : %s", err)
		}
		defer surface.Free()
		newTexture, err := r.SDLrenderer.CreateTextureFromSurface(surface)
		if err != nil {
			return nil, fmt.Errorf("texture texture error : %s", err)
		}
		return newTexture, nil
	}
	return nil, nil
}

func (r *CustomRenderer) LoadTexture(path string) (*sdl.Texture, error) {
	surface, err := img.Load(path)
	if err != nil {
		return nil, fmt.Errorf("texture surface error : %s", err)
	}
	defer surface.Free()
	newTexture, err := r.SDLrenderer.CreateTextureFromSurface(surface)
	if err != nil {
		return nil, fmt.Errorf("texture texture error : %s", err)
	}
	return newTexture, nil
}

func (r *CustomRenderer) ToggleFullscreen() {
	r.isFullscreen = !r.isFullscreen
	if r.isFullscreen {
		r.SDLwindow.SetFullscreen(sdl.WINDOW_FULLSCREEN_DESKTOP)
	} else {
		r.SDLwindow.SetFullscreen(0)
		r.SDLwindow.SetBordered(r.isBordered)
	}
}

func (r *CustomRenderer) ToggleBorders() {
	r.isBordered = !r.isBordered
	if !r.isFullscreen {
		r.SDLwindow.SetBordered(r.isBordered)
	}
}
