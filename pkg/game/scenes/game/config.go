package game

import (
	"minesweeper/pkg/game/rendering"

	"github.com/veandco/go-sdl2/sdl"
)

const (
	spritesheetFile          = "spritesheet_iso.png"
	spritesheetColumns int32 = 15
	spritesheetRows    int32 = 2
)

const (
	gridColumns = 10
	gridRows    = 10
	bombPercent = 10
	viewRange   = 22
)

var (
	colorWhite    = sdl.Color{R: 255, G: 255, B: 255, A: sdl.ALPHA_OPAQUE}
	colorDarkGrey = sdl.Color{R: 20, G: 20, B: 20, A: sdl.ALPHA_OPAQUE}
	colorBack     = sdl.Color{R: 0, G: 0, B: 0, A: sdl.ALPHA_OPAQUE}
)

const (
	actionNone rendering.ButtonActionId = iota
	actionOpenSettingsMenu
	actionExit
)

type tileSpriteID uint32

const (
	tileNoSprite           tileSpriteID = 666
	tileSpriteBorder       tileSpriteID = 0
	tileSpriteEmpty        tileSpriteID = 1
	tileHover              tileSpriteID = 2
	tileSpriteHidden       tileSpriteID = 3
	tileSpriteFlag         tileSpriteID = 4
	tileSpriteBomb         tileSpriteID = 5
	tileSpriteBombExploded tileSpriteID = 6
	tileSprite0            tileSpriteID = tileSpriteEmpty
	tileSprite1            tileSpriteID = 7
	tileSprite2            tileSpriteID = 8
	tileSprite3            tileSpriteID = 9
	tileSprite4            tileSpriteID = 10
	tileSprite5            tileSpriteID = 11
	tileSprite6            tileSpriteID = 12
	tileSprite7            tileSpriteID = 13
	tileSprite8            tileSpriteID = 14
)
const (
	textButtonExit     = "Main menu"
	textButtonSettings = "Settings"
)

type tileState byte

const (
	tileStateShown tileState = 1 << iota
	tileStateFlagged
	tileStateBomb
	tileStateExploded
	tileStateBorder
)

type gameStats struct {
	tilesHidden    int
	totalTiles     int
	flagsUsed      int
	livesRemaining int
	totalLives     int
	totalBombs     uint32
	bombsRemaining uint32
	bombsExploded  uint32
}

type gameState byte

const (
	gameStatePlaying gameState = 1 << iota
	gameStateWon
	gameStateLost
)
