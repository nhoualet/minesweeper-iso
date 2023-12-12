package rendering

import "github.com/veandco/go-sdl2/sdl"

var (
	ColorBlack    = sdl.Color{R: 0, G: 0, B: 0, A: sdl.ALPHA_OPAQUE}
	ColorYellow   = sdl.Color{R: 255, G: 255, B: 0, A: sdl.ALPHA_OPAQUE}
	ColorPink     = sdl.Color{R: 255, G: 0, B: 255, A: sdl.ALPHA_OPAQUE}
	ColorWhite    = sdl.Color{R: 255, G: 255, B: 255, A: sdl.ALPHA_OPAQUE}
	ColorGrey     = sdl.Color{R: 120, G: 120, B: 120, A: sdl.ALPHA_OPAQUE}
	ColorDarkGrey = sdl.Color{R: 20, G: 20, B: 20, A: sdl.ALPHA_OPAQUE}
	ColorRed      = sdl.Color{R: 255, G: 0, B: 0, A: sdl.ALPHA_OPAQUE}
	ColorGreen    = sdl.Color{R: 0, G: 255, B: 0, A: sdl.ALPHA_OPAQUE}
	ColorBlue     = sdl.Color{R: 0, G: 0, B: 255, A: sdl.ALPHA_OPAQUE}
)
