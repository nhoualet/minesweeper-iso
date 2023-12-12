package menuSettings

import (
	"fmt"
	"log"
	"minesweeper/pkg/config"
	"minesweeper/pkg/game/rendering"
	"minesweeper/pkg/game/scenes"

	"github.com/veandco/go-sdl2/sdl"
	"github.com/veandco/go-sdl2/ttf"
)

type SettingsScene struct {
	widgets      []rendering.Widget
	renderer     *rendering.CustomRenderer
	sceneManager *scenes.SceneManager
	font         *ttf.Font
}

func Initialize(sceneManager *scenes.SceneManager, renderer *rendering.CustomRenderer, font *ttf.Font) (*SettingsScene, error) {
	var err error
	btnCount := 0
	textboxCount := 0
	for _, widget := range widgetsData {
		if widget.wType == buttonWidget {
			btnCount += 1
		} else {
			textboxCount += 1
		}
	}
	widgets := make([]rendering.Widget, len(widgetsData))
	for i, widget := range widgetsData {
		selectable := widget.action != actionNone
		if widget.wType == buttonWidget {
			btn := rendering.NewButton(
				sdl.Rect{X: 0, Y: 0, W: 10, H: 10},
				true,
				true,
				selectable,
				widget.text,
				widget.action,
				*widget.textColor,
				widget.backgroundColor,
				widget.hoverColor,
			)
			if btn.UpdateTexture(renderer.SDLrenderer, font) != nil {
				return nil, err
			}
			widgets[i] = btn
		} else {
			widgets[i], err = rendering.NewTextbox(
				sdl.Rect{X: 0, Y: 0, W: 10, H: 10},
				true,
				true,
				widget.text,
				renderer.SDLrenderer,
				font,
				*widget.textColor,
			)

			if err != nil {
				return nil, err
			}
		}
	}
	s := &SettingsScene{widgets: widgets, font: font, renderer: renderer, sceneManager: sceneManager}
	s.updateText()
	return s, nil
}

func (s *SettingsScene) processButtonClick(b *rendering.Button) {
	switch b.ActionId {
	case actionNone:
		return
	case actionToggleFullscreen:
		s.renderer.ToggleFullscreen()
	case actionToggleBorders:
		s.renderer.ToggleBorders()
	case actionSettingIncreaseColumn:
		cfg := s.sceneManager.GetConfig()
		cfg.Game.GridColumns += 1
		s.sceneManager.SetConfig(cfg)
		s.updateText()

	case actionSettingDecreaseColumn:
		cfg := s.sceneManager.GetConfig()
		cfg.Game.GridColumns -= 1
		if cfg.Game.GridColumns < 1 {
			cfg.Game.GridColumns = 1
		}
		s.sceneManager.SetConfig(cfg)
		s.updateText()
	case actionSettingIncreaseRow:
		cfg := s.sceneManager.GetConfig()
		cfg.Game.GridRows += 1
		s.sceneManager.SetConfig(cfg)
		s.updateText()

	case actionSettingDecreaseRow:
		cfg := s.sceneManager.GetConfig()
		cfg.Game.GridRows -= 1
		if cfg.Game.GridRows < 1 {
			cfg.Game.GridRows = 1
		}
		s.sceneManager.SetConfig(cfg)
		s.updateText()

	case actionSettingIncreaseBombPercent:
		cfg := s.sceneManager.GetConfig()
		cfg.Game.BombPercent += 1
		if cfg.Game.BombPercent > 99 {
			cfg.Game.BombPercent = 99
		}
		s.sceneManager.SetConfig(cfg)
		s.updateText()

	case actionSettingDecreaseBombPercent:
		cfg := s.sceneManager.GetConfig()
		cfg.Game.BombPercent -= 1
		if cfg.Game.BombPercent < 1 {
			cfg.Game.BombPercent = 1
		}
		s.sceneManager.SetConfig(cfg)
		s.updateText()
	case actionSettingIncreaseLives:
		cfg := s.sceneManager.GetConfig()
		if cfg.Game.Lives < 0 {
			cfg.Game.Lives = 1
		} else {
			cfg.Game.Lives += 1
		}
		s.sceneManager.SetConfig(cfg)
		s.updateText()
	case actionSettingDecreaseLives:
		cfg := s.sceneManager.GetConfig()
		cfg.Game.Lives -= 1
		if cfg.Game.Lives < 1 {
			cfg.Game.Lives = -1
		}
		s.sceneManager.SetConfig(cfg)
		s.updateText()

	case actionExit:
		s.Exit()
	}
}

func (s *SettingsScene) ProcessEvent(e sdl.Event) scenes.EventState {
	switch t := e.(type) {
	case *sdl.MouseButtonEvent:
		if t.State == sdl.PRESSED {
			mousePos := sdl.Point{
				X: t.X,
				Y: t.Y,
			}
			for _, w := range s.widgets {
				if btn, ok := w.(*rendering.Button); ok {
					if btn.OnButton(mousePos) {
						s.processButtonClick(btn)
						return scenes.EventProcessed
					}
				}
			}
		}

	case *sdl.KeyboardEvent:
		keyCode := t.Keysym.Sym
		pressed := (t.State == sdl.PRESSED)
		if pressed {
			if keyCode == sdl.K_ESCAPE {
				s.Exit()
				return scenes.EventProcessed
			}
		}
	}
	return scenes.EventToProcess
}

func (s *SettingsScene) ProcessResize(w, h int32) {
	var maxHeight int32 = 0

	for _, widget := range s.widgets {
		if btn, ok := widget.(*rendering.Button); ok {
			if btn.TextureRect.H > maxHeight {
				maxHeight = btn.TextureRect.H
			}
		}
	}

	padding := maxHeight
	maxHeight += padding
	var margin int32 = 5
	y := (h - (maxHeight+margin)*int32(len(s.widgets)-1)) / 2
	for _, widget := range s.widgets[:len(s.widgets)-1] {
		widget.SetCenter(w/2, y)
		if btn, ok := widget.(*rendering.Button); ok {
			btn.SetBackgroundSize(btn.TextureRect.W+padding, maxHeight)
		} else if tbox, ok := widget.(*rendering.Textbox); ok {
			tbox.SetBackgroundSize(tbox.Rect.W+padding, maxHeight)
		}
		y += maxHeight + margin
	}
	if btn, ok := s.widgets[len(s.widgets)-1].(*rendering.Button); ok {
		btn.SetBackgroundSize(btn.TextureRect.W+padding, btn.TextureRect.H+padding)
		btn.SetTopLeft(margin, h-margin-maxHeight)
	}
}

func (s *SettingsScene) Update(deltaMS uint64) {
	mousePos := sdl.Point{}
	mousePos.X, mousePos.Y, _ = sdl.GetMouseState()
	for _, widget := range s.widgets {
		if btn, ok := widget.(*rendering.Button); ok {
			btn.SetHovered(mousePos.InRect(&btn.Rect))
		}
	}
}
func (s *SettingsScene) Draw(renderer rendering.CustomRenderer) {
	for _, widget := range s.widgets {
		widget.Draw(&renderer)
	}
}

func (s *SettingsScene) Exit() {
	s.SaveSettings()
	s.sceneManager.SetLastScene(false)
	fmt.Printf("%+v\n", s.sceneManager.GetConfig())
}

func (s *SettingsScene) Enter(reload bool) error {
	var cfg config.Config
	err := config.LoadConfig(config.ConfigFilePath, &cfg)
	if err == nil {
		s.sceneManager.SetConfig(cfg)
	}
	if reload {
		err := s.load()
		if err != nil {
			return err
		}
	}
	w, h := s.renderer.SDLwindow.GetSize()
	s.ProcessResize(w, h)
	s.updateText()
	return nil
}

func (s *SettingsScene) load() error {
	return nil
}

func (s *SettingsScene) Unload() {
	// for _, b := range s.buttons {
	// 	b.Destroy()
	// }
}

func (s *SettingsScene) SaveSettings() {
	cfg := s.sceneManager.GetConfig()
	cfg.Window.Borderless = !s.renderer.IsBordered()
	cfg.Window.Fullscreen = s.renderer.IsFullscreen()
	if !cfg.Window.Fullscreen {
		cfg.Window.Width, cfg.Window.Height = s.renderer.SDLwindow.GetSize()
	}
	s.sceneManager.SetConfig(cfg)
	err := config.SaveConfig(config.ConfigFilePath, cfg)
	if err != nil {
		log.Printf("SaveSettings error: %s\n", err.Error())
	}
}

func (s *SettingsScene) updateText() {
	cfg := s.sceneManager.GetConfig()
	if tbox, ok := s.widgets[6].(*rendering.Textbox); ok {
		text := fmt.Sprintf("Grid columns: %d", cfg.Game.GridColumns)
		tbox.SetText(text, s.renderer.SDLrenderer, s.font, rendering.ColorWhite)
	}
	if tbox, ok := s.widgets[9].(*rendering.Textbox); ok {
		text := fmt.Sprintf("Grid rows: %d", cfg.Game.GridRows)
		tbox.SetText(text, s.renderer.SDLrenderer, s.font, rendering.ColorWhite)
	}
	if tbox, ok := s.widgets[12].(*rendering.Textbox); ok {
		text := fmt.Sprintf("bombs: %d (%d%% of the tiles)", cfg.Game.GridColumns*cfg.Game.GridRows*uint32(cfg.Game.BombPercent)/100, cfg.Game.BombPercent)
		tbox.SetText(text, s.renderer.SDLrenderer, s.font, rendering.ColorWhite)
	}
	if tbox, ok := s.widgets[15].(*rendering.Textbox); ok {
		var text string
		if cfg.Game.Lives < 1 {
			text = "lives : no limit"

		} else {
			text = fmt.Sprintf("lives : %d", cfg.Game.Lives)
		}
		tbox.SetText(text, s.renderer.SDLrenderer, s.font, rendering.ColorWhite)
	}
}

func (s *SettingsScene) IsLoaded() bool {
	return true
}

func (s *SettingsScene) NeedsRedraw() bool {
	return true
}
