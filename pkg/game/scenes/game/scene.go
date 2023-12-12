package game

import (
	"fmt"
	"math/rand"
	"minesweeper/pkg/config"
	"minesweeper/pkg/game/rendering"
	"minesweeper/pkg/game/scenes"

	"github.com/veandco/go-sdl2/sdl"
	"github.com/veandco/go-sdl2/ttf"
)

type GameScene struct {
	widgets         [2]rendering.Widget
	statsMessage    [5]*rendering.Textbox
	bigMessage      *rendering.Textbox
	bigMessageRect  sdl.Rect
	statsRect       sdl.Rect
	stats           gameStats
	renderer        *rendering.CustomRenderer
	sceneManager    *scenes.SceneManager
	font            *ttf.Font
	spreadsheet     *rendering.Spritesheet
	grid            *grid
	tileSize        sdl.Rect
	player          player
	isLoaded        bool
	needsRedraw     bool
	state           gameState
	partyGameConfig config.GameConfig
	keyConfig       config.ControlCodes
}

func Initialize(sceneManager *scenes.SceneManager, renderer *rendering.CustomRenderer, font *ttf.Font) (*GameScene, error) {
	var err error
	var widgets [2]rendering.Widget
	cfg := sceneManager.GetConfig()

	widgets[0] = rendering.NewButton(
		sdl.Rect{X: 0, Y: 0, W: 10, H: 10},
		true,
		true,
		false,
		textButtonExit,
		actionExit,
		colorWhite,
		&colorBack,
		&colorWhite,
	)

	if err := widgets[0].(*rendering.Button).UpdateTexture(renderer.SDLrenderer, font); err != nil {
		return nil, err
	}
	widgets[1] = rendering.NewButton(
		sdl.Rect{X: 0, Y: 0, W: 10, H: 10},
		true,
		true,
		false,
		textButtonSettings,
		actionOpenSettingsMenu,
		colorWhite,
		&colorBack,
		&colorWhite,
	)

	if err := widgets[1].(*rendering.Button).UpdateTexture(renderer.SDLrenderer, font); err != nil {
		return nil, err
	}

	spritesheet, err := rendering.NewSpritesheet(renderer, cfg.Window.ResourcesPath+spritesheetFile, spritesheetColumns, spritesheetRows)
	if err != nil {
		return nil, err
	}
	var statsMessage [5]*rendering.Textbox
	for i := range statsMessage {
		statsMessage[i], err = rendering.NewTextbox(sdl.Rect{X: 0, Y: 0, W: 0, H: 0}, true, true, "x", renderer.SDLrenderer, font, rendering.ColorWhite)
		if err != nil {
			return nil, err
		}
	}
	bigMessage, err := rendering.NewTextbox(sdl.Rect{X: 0, Y: 0, W: 0, H: 0}, true, true, "", renderer.SDLrenderer, font, rendering.ColorWhite)
	if err != nil {
		return nil, err
	}
	s := &GameScene{
		widgets:         widgets,
		statsMessage:    statsMessage,
		statsRect:       sdl.Rect{X: 0, Y: 0, W: 0, H: 0},
		font:            font,
		renderer:        renderer,
		sceneManager:    sceneManager,
		spreadsheet:     spritesheet,
		grid:            nil,
		tileSize:        sdl.Rect{X: 0, Y: 0, W: 0, H: 0},
		player:          player{pos: playerPos{col: 0, row: 0}},
		isLoaded:        false,
		needsRedraw:     true,
		stats:           gameStats{},
		bigMessage:      bigMessage,
		state:           gameStatePlaying,
		partyGameConfig: cfg.Game,
		keyConfig:       cfg.Controls.Codes,
	}
	return s, nil
}

func (s *GameScene) processButtonClick(b *rendering.Button) {
	switch b.ActionId {
	case actionNone:
		return
	case actionOpenSettingsMenu:
		s.sceneManager.SetScene("settings", false)
	case actionExit:
		s.Exit()
	}

}

func (s *GameScene) replaceStateMessage() {
	var maxWidth int32 = 0
	var maxHeight int32 = 0
	var lineCount int32
	for _, tbox := range s.statsMessage {
		if tbox.Rect.W > maxWidth {
			maxWidth = tbox.Rect.W
		}
		if tbox.Rect.H > maxHeight {
			maxHeight = tbox.Rect.H
		}
		if tbox.Text != "" {
			lineCount += 1
		}
	}
	w, h := s.renderer.SDLwindow.GetSize()
	s.statsRect.X = w - maxWidth - 10
	s.statsRect.Y = h / 5
	s.statsRect.W = maxWidth + 10*2
	s.statsRect.H = (maxHeight+10)*lineCount + 10
	top := s.statsRect.Y + 10
	left := s.statsRect.X + 10
	for _, tbox := range s.statsMessage {
		if tbox.Text != "" {
			tbox.SetBackgroundSize(maxWidth, tbox.Rect.H)
			tbox.SetTopLeft(left, top)
			top += maxHeight + 10
		}
	}
}

func (s *GameScene) updateBigMessage(msg string) {
	s.bigMessage.SetText(msg, s.renderer.SDLrenderer, s.font, rendering.ColorWhite)
	s.replaceBigMessage()
}
func (s *GameScene) replaceBigMessage() {
	w, h := s.renderer.SDLwindow.GetSize()
	s.bigMessage.SetCenter(w/2, h/10+s.bigMessage.Rect.H)
	s.bigMessageRect = sdl.Rect{X: s.bigMessage.Rect.X - 10, Y: s.bigMessage.Rect.Y - 10, W: s.bigMessage.Rect.W + 20, H: s.bigMessage.Rect.H + 20}
}

func (s *GameScene) updateStateMessage(msg string) {
	var msgs [5]string
	if s.state == gameStatePlaying {
		var livesMsg string
		if s.stats.totalLives < 0 {
			livesMsg = "lives left: ∞"
		} else {
			livesMsg = fmt.Sprintf("lives left: %d/%d", s.stats.livesRemaining, s.stats.totalLives)
		}
		msgs = [...]string{
			msg,
			fmt.Sprintf("Tiles hidden: %d/%d", s.stats.tilesHidden, s.stats.totalTiles),
			fmt.Sprintf("flags used: %d", s.stats.flagsUsed),
			fmt.Sprintf("Bombs remaining: %d", s.stats.bombsRemaining),
			livesMsg,
		}
	} else {
		var livesMsg string
		if s.stats.totalLives < 0 {
			livesMsg = "unlimited lives"
		} else if s.stats.livesRemaining == 0 {
			livesMsg = fmt.Sprintf("%d/%d lives left", s.stats.livesRemaining, s.stats.totalLives)
		}
		msgs = [...]string{
			msg,
			"",
			fmt.Sprintf("%d tiles", s.stats.totalTiles),
			fmt.Sprintf("%d/%d bombs exploded", s.stats.bombsExploded, s.stats.totalBombs),
			livesMsg,
		}
	}
	for i := range msgs {
		s.statsMessage[i].SetText(msgs[i], s.renderer.SDLrenderer, s.font, rendering.ColorWhite)
	}
	s.replaceStateMessage()
}

func (s *GameScene) discoverGrid() {
	for col := range s.grid.tiles {
		for row := range s.grid.tiles[col] {
			if !s.grid.tiles[col][row].has(tileStateShown) {
				s.grid.tiles[col][row].set(tileStateShown)
			}
		}
	}
	s.stats.tilesHidden = 0
	s.stats.bombsRemaining = 0
}

func (s *GameScene) checkGameState() {
	if s.stats.totalLives >= 0 && s.stats.livesRemaining <= 0 {
		s.updateBigMessage("Game lost (no lives left) press [R] to replay")
		s.updateStateMessage("")
		for col := range s.grid.tiles {
			for row := range s.grid.tiles[col] {
				if !s.grid.tiles[col][row].has(tileStateShown) {
					s.grid.tiles[col][row].set(tileStateShown)
				}
			}
		}
		s.stats.tilesHidden = 0
		s.state = gameStateLost
		s.updateStateMessage("")
		return
	}
	if uint32(s.stats.flagsUsed) > s.stats.bombsRemaining {
		s.updateBigMessage("At least one flag isn't right")
		return
	}
	if uint32(s.stats.flagsUsed) == s.stats.bombsRemaining {
		won := true
		for col := range s.grid.tiles {
			for row := range s.grid.tiles[col] {
				if !s.grid.tiles[col][row].has(tileStateShown) && s.grid.tiles[col][row].has(tileStateFlagged) != s.grid.tiles[col][row].has(tileStateBomb) {
					won = false
					break
				}
			}
		}
		if won {
			s.discoverGrid()
			s.state = gameStateWon
			s.updateBigMessage("Game won, press [R] to replay")
			s.updateStateMessage("")

			s.updateStateMessage("")
		} else {
			s.updateBigMessage("At least one flag isn't right")
		}
		return
	}
	if uint32(s.stats.tilesHidden) == s.stats.bombsRemaining {
		s.discoverGrid()
		s.stats.tilesHidden = 0
		s.stats.bombsRemaining = 0
		s.state = gameStateWon
		s.updateBigMessage("Game won, press [R] to replay")
		s.updateStateMessage("")
		return
	}
	s.updateBigMessage("")
}

func (s *GameScene) ProcessResize(w, h int32) {
	s.replaceStateMessage()
	s.replaceBigMessage()
	if btn, ok := s.widgets[0].(*rendering.Button); ok {
		btn.SetBackgroundSize(btn.TextureRect.W+10, btn.TextureRect.H+10)
		btn.SetTopLeft(10, h-btn.Rect.H-10)
	}
	if btn, ok := s.widgets[1].(*rendering.Button); ok {
		btn.SetBackgroundSize(btn.TextureRect.W+10, btn.TextureRect.H+10)
		btn.SetTopLeft(w-btn.Rect.W-10, h-btn.Rect.H-10)
	}
	minTileSize := sdl.Rect{X: 0, Y: 0, W: w / 11, H: h / 11}
	if minTileSize.W < minTileSize.H {
		minTileSize.H = minTileSize.W
	} else {
		minTileSize.W = minTileSize.H
	}
	s.tileSize.W = w / (int32(s.grid.columns) + 1)
	s.tileSize.H = h / (int32(s.grid.rows) + 1)
	if s.tileSize.W < s.tileSize.H {
		s.tileSize.H = s.tileSize.W
	} else {
		s.tileSize.W = s.tileSize.H
	}
	if s.tileSize.W < minTileSize.W {
		s.tileSize.W = minTileSize.W
		s.tileSize.H = minTileSize.H
	}
	p := cartesianToIso(sdl.Point{X: int32(s.grid.columns), Y: int32(s.grid.rows)}, s.tileSize)
	s.tileSize.X = (w - p.X) / 2
	s.tileSize.Y = (h - p.Y) / 2
	s.needsRedraw = true
}

func (s *GameScene) Update(deltaMS uint64) {
	mousePos := sdl.Point{}
	mousePos.X, mousePos.Y, _ = sdl.GetMouseState()
	var oldState bool
	for _, widget := range s.widgets {
		if btn, ok := widget.(*rendering.Button); ok {
			oldState = btn.IsHovered()
			btn.SetHovered(mousePos.InRect(&btn.Rect))
			if btn.IsHovered() != oldState {
				s.needsRedraw = true
			}
		}
	}
}

func (s *GameScene) tileValueToId(t tile) tileSpriteID {
	if t.has(tileStateBorder) {
		return tileSpriteBorder
	}
	if t.has(tileStateFlagged) {
		return tileSpriteFlag
	}
	if !t.has(tileStateShown) {
		return tileSpriteHidden
	}
	if t.has(tileStateBomb) {
		if t.has(tileStateExploded) {
			return tileSpriteBombExploded
		} else {
			return tileSpriteBomb
		}
	}
	switch t.bombAround {
	case 0:
		return tileNoSprite
	case 1:
		return tileSprite1
	case 2:
		return tileSprite2
	case 3:
		return tileSprite3
	case 4:
		return tileSprite4
	case 5:
		return tileSprite5
	case 6:
		return tileSprite6
	case 7:
		return tileSprite7
	case 8:
		return tileSprite8
	}
	return tileSprite0 + tileSpriteID(t.bombAround)
}

func (s *GameScene) drawTileGround(renderer rendering.CustomRenderer, cstart, cstop, rstart, rstop int32) {
	w, h := renderer.SDLwindow.GetSize()
	rect := s.tileSize
	dp := cartesianToIso(sdl.Point{X: s.player.pos.col, Y: s.player.pos.row}, s.tileSize)
	s.spreadsheet.SelectSprite(uint32(tileSpriteEmpty))
	var pos sdl.Point
	for r := rstart; r < rstop; r += 1 {
		for c := cstart; c < cstop; c += 1 {
			pos = cartesianToIso(sdl.Point{X: c, Y: r}, rect)
			rect.X = pos.X - dp.X + w/2
			rect.Y = pos.Y - dp.Y + h/2
			var spriteID uint32
			if s.grid.tiles[c][r].has(tileStateBorder) {
				spriteID = uint32(tileSpriteBorder)
			} else {
				spriteID = uint32(tileSpriteEmpty)
				if c == s.player.pos.col && r == s.player.pos.row {
					spriteID += uint32(spritesheetColumns)
				}
			}
			s.spreadsheet.SelectSprite(spriteID)
			s.spreadsheet.Draw(&renderer, rect)
		}
	}
}

func (s *GameScene) drawTileContent(renderer rendering.CustomRenderer, cstart, cstop, rstart, rstop int32) {
	w, h := renderer.SDLwindow.GetSize()
	rect := s.tileSize
	dp := cartesianToIso(sdl.Point{X: s.player.pos.col, Y: s.player.pos.row}, s.tileSize)
	var pos sdl.Point
	for r := rstart; r < rstop; r += 1 {
		for c := cstart; c < cstop; c += 1 {
			spriteID := s.tileValueToId(s.grid.tiles[c][r])
			if spriteID != tileNoSprite {
				id := uint32(spriteID)
				if c == s.player.pos.col && r == s.player.pos.row {
					id += uint32(spritesheetColumns)
				}
				s.spreadsheet.SelectSprite(id)
				pos = cartesianToIso(sdl.Point{X: c, Y: r}, rect)
				rect.X = pos.X - dp.X + w/2
				rect.Y = pos.Y - dp.Y + h/2
				s.spreadsheet.Draw(&renderer, rect)
			}
		}
	}
}

func (s *GameScene) Draw(renderer rendering.CustomRenderer) {
	if !s.needsRedraw {
		return
	}
	rstart := s.player.pos.row - viewRange
	rstop := s.player.pos.row + viewRange
	if rstart < 0 {
		rstart = 0
	}
	if rstop > int32(s.grid.rows) {
		rstop = int32(s.grid.rows)
	}

	cstart := s.player.pos.col - viewRange
	cstop := s.player.pos.col + viewRange
	if cstart < 0 {
		cstart = 0
	}
	if cstop > int32(s.grid.columns) {
		cstop = int32(s.grid.columns)
	}
	s.drawTileGround(renderer, cstart, cstop, rstart, rstop)
	s.drawTileContent(renderer, cstart, cstop, rstart, rstop)
	renderer.SDLrenderer.SetDrawColor(0, 0, 0, sdl.ALPHA_OPAQUE)
	renderer.SDLrenderer.FillRect(&s.statsRect)
	if s.bigMessage.Text != "" {
		renderer.SDLrenderer.FillRect(&s.bigMessageRect)
		s.bigMessage.Draw(&renderer)
	}

	for _, tbox := range s.statsMessage {
		tbox.Draw(&renderer)
	}
	for _, widget := range s.widgets {
		widget.Draw(&renderer)
	}
	s.needsRedraw = false
}

func (s *GameScene) Exit() {
	// s.grid = nil
	s.sceneManager.SetSceneDefault(false)
}

func (s *GameScene) Enter(reload bool) error {
	if reload {
		var cfg config.Config
		if config.LoadConfig(config.ConfigFilePath, &cfg) == nil {
			s.sceneManager.SetConfig(cfg)
		}
		err := s.load()
		if err != nil {
			return err
		}
	}

	w, h := s.renderer.SDLwindow.GetSize()
	s.ProcessResize(w, h)
	return nil
}

func (s *GameScene) load() error {
	cfg := s.sceneManager.GetConfig()
	bombCount := cfg.Game.GridColumns * cfg.Game.GridRows * uint32(cfg.Game.BombPercent) / 100
	s.grid = newGrid(cfg.Game.GridColumns, cfg.Game.GridRows, bombCount)
	s.partyGameConfig = cfg.Game
	tileCount := int(cfg.Game.GridColumns * cfg.Game.GridRows)
	s.stats = gameStats{
		tilesHidden:    tileCount,
		totalTiles:     tileCount,
		flagsUsed:      0,
		livesRemaining: cfg.Game.Lives,
		totalLives:     cfg.Game.Lives,
		totalBombs:     bombCount,
		bombsRemaining: bombCount,
		bombsExploded:  0,
	}
	playerPlaced := false
	for !playerPlaced {
		col := int32(rand.Intn(int(s.grid.columns)))
		row := int32(rand.Intn(int(s.grid.rows)))
		var firstOpenCount int
		if !s.grid.tiles[col][row].has(tileStateBomb|tileStateBorder) && s.grid.tiles[col][row].bombAround == 0 && s.openTile(col, row, &firstOpenCount) == nil {
			s.player.pos.col = col
			s.player.pos.row = row
			s.stats.tilesHidden -= firstOpenCount
			playerPlaced = true
		}
	}
	fmt.Printf("(load) New game config: %+v\n", cfg.Game)
	s.updateStateMessage("")
	s.updateBigMessage("")
	s.isLoaded = true
	s.needsRedraw = true
	s.state = gameStatePlaying
	return nil
}

func (s *GameScene) Unload() {
	// for _, b := range s.buttons {
	// 	b.Destroy()
	// }
}

func (s *GameScene) IsLoaded() bool {
	return s.isLoaded
}

func (s *GameScene) NeedsRedraw() bool {
	return s.needsRedraw
}
