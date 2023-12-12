package scenes

import (
	"minesweeper/pkg/game/rendering"

	"github.com/veandco/go-sdl2/sdl"
)

type MenuLine struct {
	widgets []rendering.Widget
}

type Menu struct {
	widgetLines []MenuLine
}

func NewMenuLine(w []rendering.Widget) MenuLine {
	return MenuLine{
		widgets: w,
	}
}

func (l *MenuLine) Height() int32 {
	var maxHeight int32 = 0
	for _, w := range l.widgets {
		if w.Height() > maxHeight {
			maxHeight = w.Height()
		}
	}
	return maxHeight
}

func (l *MenuLine) Width() int32 {
	var width int32 = 0
	for _, w := range l.widgets {
		width += w.Width()
	}
	return width
}

func NewMenu() Menu {
	return Menu{
		make([]MenuLine, 0),
	}
}

func (m *Menu) AppendLine(l MenuLine) {
	m.widgetLines = append(m.widgetLines, l)
}

func (m *Menu) InsertLine(lineID int, l MenuLine) {
	m.widgetLines = append(m.widgetLines[:lineID+1], m.widgetLines[lineID:]...)
	m.widgetLines[lineID] = l
}

func (m *Menu) Width() int32 {
	var maxWidth int32 = 0
	for _, l := range m.widgetLines {
		if l.Width() > maxWidth {
			maxWidth = l.Width()
		}
	}
	return maxWidth
}

func (m *Menu) Height() int32 {
	var height int32 = 0
	for _, l := range m.widgetLines {
		height += l.Height()
	}
	return height
}

func (m *Menu) Resize(rect sdl.Rect) {
	w := m.Width()
	h := m.Height()
	pos := rect
	if w < rect.W {
		pos.X = rect.X - (rect.W-w)/2
	}
	if h < rect.H {
		pos.Y = rect.Y - (rect.H-h)/2
	}
	currentY := pos.Y
	currentX := pos.X
	for _, l := range m.widgetLines {
		for _, w := range l.widgets {
			w.SetTopLeft(currentX, currentY)
			currentX += w.Width()
		}
		currentY += l.Height()
		currentX = pos.X
	}
}

func (l *MenuLine) Draw(r *rendering.CustomRenderer) {
	for _, w := range l.widgets {
		w.Draw(r)
	}
}
func (m *Menu) Draw(r *rendering.CustomRenderer) {
	for _, l := range m.widgetLines {
		l.Draw(r)
	}
}
