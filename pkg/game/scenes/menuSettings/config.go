package menuSettings

import (
	"minesweeper/pkg/game/rendering"

	"github.com/veandco/go-sdl2/sdl"
)

const (
	actionNone rendering.ButtonActionId = iota
	actionToggleFullscreen
	actionToggleBorders
	actionExit
	actionSettingIncreaseColumn
	actionSettingDecreaseColumn
	actionSettingIncreaseRow
	actionSettingDecreaseRow

	actionSettingIncreaseBombPercent
	actionSettingDecreaseBombPercent

	actionSettingIncreaseLives
	actionSettingDecreaseLives
	actionSettingToggleInfiniteLives
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

type widgetID int

// const (
// 	widgetSettingsText = iota
// 	widgetFullscreenToggleBtn
// 	widgetBorderToggleBtn
// 	widgetEmpty
// 	widgetGridSettingsText
// 	widgetGridColumnsText
// 	widgetGridColumnsIncreaseBtn
// 	widgetGridColumnsDecreaseBtn

// 	widgetGridRowsText
// 	widgetGridRowsIncreaseBtn
// 	widgetGridRowsDecreaseBtn

// 	widgetBombText
// 	widgetBombIncreaseBtn
// 	widgetBombDecreaseBtn

// 	widgetLivesText
// 	widgetLivesIncreateBtn
// 	widgetLivesDecreaseBtn
// 	widgetLivesInfiniteBtn
// )

var widgetsData = [...]widgetLoadingData{
	{textboxWidget, "Window settings", actionNone, &rendering.ColorWhite, &rendering.ColorBlack, &rendering.ColorWhite},
	{buttonWidget, "Toggle fullscreen", actionToggleFullscreen, &rendering.ColorWhite, &rendering.ColorBlack, &rendering.ColorWhite},
	{buttonWidget, "Toggle borders", actionToggleBorders, &rendering.ColorWhite, &rendering.ColorBlack, &rendering.ColorWhite},
	{textboxWidget, "", actionNone, &rendering.ColorWhite, &rendering.ColorBlack, &rendering.ColorWhite},
	{textboxWidget, "---", actionNone, &rendering.ColorDarkGrey, &rendering.ColorBlack, &rendering.ColorWhite},
	{textboxWidget, "Grid settings (changes for the next game)", actionNone, &rendering.ColorWhite, &rendering.ColorBlack, &rendering.ColorWhite},
	{textboxWidget, "Grid columns : {X}", actionNone, &rendering.ColorWhite, &rendering.ColorBlack, &rendering.ColorWhite},
	{buttonWidget, "+", actionSettingIncreaseColumn, &rendering.ColorWhite, &rendering.ColorBlack, &rendering.ColorWhite},
	{buttonWidget, "-", actionSettingDecreaseColumn, &rendering.ColorWhite, &rendering.ColorBlack, &rendering.ColorWhite},
	{textboxWidget, "Grid rows : {X}", actionNone, &rendering.ColorWhite, &rendering.ColorBlack, &rendering.ColorWhite},
	{buttonWidget, "+", actionSettingIncreaseRow, &rendering.ColorWhite, &rendering.ColorBlack, &rendering.ColorWhite},
	{buttonWidget, "-", actionSettingDecreaseRow, &rendering.ColorWhite, &rendering.ColorBlack, &rendering.ColorWhite},
	{textboxWidget, "% of Bomb : {X}%", actionNone, &rendering.ColorWhite, &rendering.ColorBlack, &rendering.ColorWhite},
	{buttonWidget, "+", actionSettingIncreaseBombPercent, &rendering.ColorWhite, &rendering.ColorBlack, &rendering.ColorWhite},
	{buttonWidget, "-", actionSettingDecreaseBombPercent, &rendering.ColorWhite, &rendering.ColorBlack, &rendering.ColorWhite},
	{textboxWidget, "lives : {X}", actionNone, &rendering.ColorWhite, &rendering.ColorBlack, &rendering.ColorWhite},
	{buttonWidget, "+", actionSettingIncreaseLives, &rendering.ColorWhite, &rendering.ColorBlack, &rendering.ColorWhite},
	{buttonWidget, "-", actionSettingDecreaseLives, &rendering.ColorWhite, &rendering.ColorBlack, &rendering.ColorWhite},
	{buttonWidget, "∞", actionNone, &rendering.ColorWhite, &rendering.ColorBlack, &rendering.ColorWhite},

	{buttonWidget, "Go back", actionExit, &rendering.ColorWhite, &rendering.ColorBlack, &rendering.ColorWhite},
}
