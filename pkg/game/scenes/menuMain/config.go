package menuMain

import (
	"minesweeper/pkg/game/rendering"

	"github.com/veandco/go-sdl2/sdl"
)

var (
	backgroundColor = sdl.Color{R: 0, G: 0, B: 0, A: sdl.ALPHA_OPAQUE}
	hoverColor      = sdl.Color{R: 255, G: 255, B: 0, A: sdl.ALPHA_OPAQUE}
	secondaryColor  = sdl.Color{R: 255, G: 255, B: 255, A: sdl.ALPHA_OPAQUE}
	// primaryColor    = sdl.Color{R: 255, G: 50, B: 255, A: sdl.ALPHA_OPAQUE}
	tertiaryColor = sdl.Color{R: 120, G: 120, B: 120, A: sdl.ALPHA_OPAQUE}
	// wipColor        = sdl.Color{R: 255, G: 0, B: 0, A: 100}
)

const (
	actionNone rendering.ButtonActionId = iota
	actionOpenSettingsMenu
	actionOpenNewGame
	actionOpenLastGame
	actionOpenBrowserGithub
	actionOpenBrowserInstagram
	actionExit

	githubLogoFile    = "github.png"
	instagramLogoFile = "instagram.png"

	githubURL   = "https://github.com/nhoualet"
	intagramURL = "https://www.instagram.com/13noodles_"
)

type widgetType int

const (
	buttonWidget = iota
	textboxWidget
)

type widgetLoadingData struct {
	wType           widgetType
	text            string
	action          rendering.ButtonActionId
	textColor       *sdl.Color
	backgroundColor *sdl.Color
	hoverColor      *sdl.Color
}

const (
	textTitle              = "Isometric minesweeper"
	textButtonNewGame      = "New Game"
	textButtonContinueGame = "Continue"
	textButtonExit         = "Exit"
	textButtonSettings     = "Settings"
)

var widgetsData = [...]widgetLoadingData{
	{textboxWidget, textTitle, actionNone, &secondaryColor, nil, nil},
	{textboxWidget, "", actionNone, &secondaryColor, nil, nil},
	{buttonWidget, textButtonNewGame, actionOpenNewGame, &secondaryColor, &backgroundColor, &tertiaryColor},
	{buttonWidget, textButtonContinueGame, actionOpenLastGame, &secondaryColor, &backgroundColor, &tertiaryColor},
	{buttonWidget, textButtonSettings, actionOpenSettingsMenu, &secondaryColor, &backgroundColor, &hoverColor},
	{buttonWidget, textButtonExit, actionExit, &secondaryColor, &backgroundColor, &hoverColor},
}
var socialsWidgetData = [...]widgetLoadingData{
	{buttonWidget, githubLogoFile, actionOpenBrowserGithub, &secondaryColor, nil, nil},
	{buttonWidget, instagramLogoFile, actionOpenBrowserInstagram, &secondaryColor, nil, nil},
}
