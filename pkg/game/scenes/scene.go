package scenes

import (
	"minesweeper/pkg/game/rendering"

	"github.com/veandco/go-sdl2/sdl"
)

type EventState int

const (
	EventProcessed = iota
	EventToProcess
)

type Scene interface {
	ProcessEvent(e sdl.Event) EventState
	Update(deltaMS uint64)
	Draw(renderer rendering.CustomRenderer)
	Exit()
	Unload()
	// load() error
	Enter(reload bool) error
	ProcessResize(w, h int32)
	IsLoaded() bool
	NeedsRedraw() bool
}
