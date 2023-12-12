package menuMain

import (
	"errors"
	"minesweeper/pkg/config"
	"minesweeper/pkg/game/rendering"
	"minesweeper/pkg/game/scenes"

	"github.com/pkg/browser"
	"github.com/veandco/go-sdl2/sdl"
	"github.com/veandco/go-sdl2/ttf"
)

type MainScene struct {
	widgets          []rendering.Widget
	renderer         *rendering.CustomRenderer
	sceneManager     *scenes.SceneManager
	font             *ttf.Font
	lang             config.MainMenuLang
	selectedWidgetID int
}

func Initialize(sceneManager *scenes.SceneManager, renderer *rendering.CustomRenderer, font *ttf.Font, lang config.MainMenuLang) (*MainScene, error) {
	var err error
	btnCount := 0
	textboxCount := 0
	cfg := sceneManager.GetConfig()
	for _, widget := range widgetsData {
		if widget.wType == buttonWidget {
			btnCount += 1
		} else {
			textboxCount += 1
		}
	}
	widgets := make([]rendering.Widget, len(widgetsData)+len(socialsWidgetData))
	for i, widget := range widgetsData {
		if widget.wType == buttonWidget {
			selectable := widget.action != actionNone
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
			if err := btn.UpdateTexture(renderer.SDLrenderer, font); err != nil {
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
	for i, widget := range socialsWidgetData {
		if widget.wType == buttonWidget {

			btn := rendering.NewButton(
				sdl.Rect{X: 0, Y: 0, W: 10, H: 10},
				true,
				true,
				true,
				"",
				widget.action,
				*widget.textColor,
				nil,
				widget.hoverColor,
			)
			if err := btn.UpdateTexture(renderer.SDLrenderer, font); err != nil {
				return nil, err
			}
			widgets[len(widgetsData)+i] = btn
			if w, ok := widgets[len(widgetsData)+i].(*rendering.Button); ok {
				texture, err := renderer.LoadTexture(cfg.Window.ResourcesPath + widget.text)
				if err != nil {
					return nil, err
				}
				w.SetTexture(texture)
			}
		} else {
			return nil, errors.New("got socialsWidgetData with wtype != buttonWidget")
		}
	}
	s := &MainScene{
		widgets:          widgets,
		font:             font,
		renderer:         renderer,
		sceneManager:     sceneManager,
		selectedWidgetID: -1,
		lang:             lang,
	}
	return s, nil
}

func (s *MainScene) processButtonClick(b *rendering.Button) {
	switch b.ActionId {
	case actionNone:
		return
	case actionOpenSettingsMenu:
		s.sceneManager.SetScene("settings", false)
	case actionOpenNewGame:
		s.sceneManager.SetScene("game", true)
	case actionOpenLastGame:
		s.sceneManager.SetScene("game", false)
	case actionOpenBrowserGithub:
		browser.OpenURL(githubURL)
	case actionOpenBrowserInstagram:
		browser.OpenURL(intagramURL)
	case actionExit:
		s.sceneManager.Quit()
	}

}

func (s *MainScene) ProcessEvent(e sdl.Event) scenes.EventState {
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
				s.sceneManager.Quit()
				return scenes.EventProcessed
			} else if keyCode == sdl.K_f {
				s.renderer.ToggleFullscreen()
			} else if keyCode == sdl.K_TAB {
				if (sdl.GetModState() & sdl.KMOD_SHIFT) != 0 {
					s.selectPreviousWidget()

				} else {
					s.selectNextWidget()
				}
			} else if keyCode == sdl.K_RETURN && s.selectedWidgetID != -1 {
				if btn, ok := s.widgets[s.selectedWidgetID].(*rendering.Button); ok {
					s.processButtonClick(btn)
				}
			}
		}
	}
	return scenes.EventToProcess
}

func (s *MainScene) toggleSelect(w interface{ rendering.Selectable }) {
	w.Selected(!w.IsSelected())
}

func (s *MainScene) selectPreviousWidget() {
	if s.selectedWidgetID == -1 {
		var widgetID int
		for widgetID = len(s.widgets) - 1; widgetID > -1; {
			if newSelected, ok := s.widgets[widgetID].(rendering.Selectable); ok && newSelected.IsSelectable() {
				s.toggleSelect(newSelected)
				s.selectedWidgetID = widgetID
				break
			}
			widgetID += 1
		}

	} else {
		var widgetID int
		for widgetID = s.selectedWidgetID - 1; widgetID != s.selectedWidgetID; {
			if widgetID == -1 {
				widgetID = len(s.widgets) - 1
			}
			if newSelected, ok := s.widgets[widgetID].(rendering.Selectable); ok && newSelected.IsSelectable() {
				if oldSelected, ok := s.widgets[s.selectedWidgetID].(rendering.Selectable); ok {
					s.toggleSelect(oldSelected)
				}
				s.toggleSelect(newSelected)
				s.selectedWidgetID = widgetID
				break
			}
			widgetID -= 1
		}
	}
}

func (s *MainScene) selectNextWidget() {
	if s.selectedWidgetID == -1 {
		var widgetID int
		for widgetID = 0; widgetID < len(s.widgets)-1; {
			if newSelected, ok := s.widgets[widgetID].(rendering.Selectable); ok && newSelected.IsSelectable() {
				s.toggleSelect(newSelected)
				s.selectedWidgetID = widgetID
				break
			}
			widgetID += 1
		}

	} else {
		var widgetID int
		for widgetID = s.selectedWidgetID + 1; widgetID != s.selectedWidgetID; {
			if widgetID == len(s.widgets) {
				widgetID = 0
			}
			if newSelected, ok := s.widgets[widgetID].(rendering.Selectable); ok && newSelected.IsSelectable() {
				if oldSelected, ok := s.widgets[s.selectedWidgetID].(rendering.Selectable); ok {
					s.toggleSelect(oldSelected)
				}
				s.toggleSelect(newSelected)
				s.selectedWidgetID = widgetID
				break
			}
			widgetID += 1
		}
	}
}

func (s *MainScene) ProcessResize(w, h int32) {
	var maxWidth int32 = 0
	var maxHeight int32 = 0

	for _, w := range s.widgets[:len(s.widgets)-len(socialsWidgetData)] {
		if btn, ok := w.(*rendering.Button); ok {
			if btn.TextureRect.W > maxWidth {
				maxWidth = btn.TextureRect.W
			}
			if btn.TextureRect.H > maxHeight {
				maxHeight = btn.TextureRect.H
			}
		}
	}

	padding := maxHeight
	maxHeight += padding
	maxWidth += padding
	var margin int32 = 5
	top := (h - maxHeight*int32(len(s.widgets))) / 2
	for i, widget := range s.widgets[:len(s.widgets)-len(socialsWidgetData)] {
		widget.SetCenter(w/2, top+int32(i)*(maxHeight+margin))
		if btn, ok := widget.(*rendering.Button); ok {
			btn.SetBackgroundSize(maxWidth, maxHeight)
		}
	}

	linksLeft := (w-(maxHeight+margin)*int32(len(socialsWidgetData)))/2 + maxHeight/2

	for i := len(s.widgets) - len(socialsWidgetData); i < len(s.widgets); i++ {
		if btn, ok := s.widgets[i].(*rendering.Button); ok {
			btn.SetCenter(linksLeft, h-margin-maxHeight/2)
			btn.SetBackgroundSize(maxHeight, maxHeight)
			btn.SetTextureSize(maxHeight, maxHeight)
			linksLeft += maxHeight + margin
		}
	}
}

func (s *MainScene) Update(deltaMS uint64) {
	mousePos := sdl.Point{}
	mousePos.X, mousePos.Y, _ = sdl.GetMouseState()
	for _, widget := range s.widgets {
		if btn, ok := widget.(*rendering.Button); ok {
			btn.SetHovered(mousePos.InRect(&btn.Rect))
		}
	}
}
func (s *MainScene) Draw(renderer rendering.CustomRenderer) {
	for _, widget := range s.widgets {
		widget.Draw(&renderer)
	}
}

func (s *MainScene) Exit() {
	// sm.Quit()
}

func (s *MainScene) Enter(reload bool) error {
	if reload {
		err := s.load()
		if err != nil {
			return err
		}
	}
	if gameScene, err := s.sceneManager.GetScene("game"); err == nil {
		continueWidget := s.widgets[3]
		if continueBtn, ok := continueWidget.(*rendering.Button); ok {
			if gameScene.IsLoaded() {
				continueBtn.SetText(continueBtn.Text, s.renderer.SDLrenderer, s.font, rendering.ColorWhite)
				continueBtn.SetHoverColor(&rendering.ColorWhite)
				continueBtn.SetSelectable(true)
				continueBtn.SetAction(actionOpenLastGame)
			} else {
				continueBtn.SetText(continueBtn.Text, s.renderer.SDLrenderer, s.font, rendering.ColorGrey)
				continueBtn.SetSelectable(false)
				continueBtn.SetHoverColor(nil)
				continueBtn.SetAction(actionNone)
			}
		}
	}

	w, h := s.renderer.SDLwindow.GetSize()
	s.ProcessResize(w, h)
	return nil
}

func (s *MainScene) load() error {
	return nil
}

func (s *MainScene) Unload() {
	// for _, b := range s.buttons {
	// 	b.Destroy()
	// }
}

func (s *MainScene) IsLoaded() bool {
	return true
}

func (s *MainScene) NeedsRedraw() bool {
	return true
}
