package rendering

type Drawable interface {
	Draw(r *CustomRenderer)
}

type Widget interface {
	Drawable
	SetCenter(x, y int32)
	SetTopLeft(x, y int32)
	SetX(x int32, isCenter bool)
	SetY(y int32, isCenter bool)
	Width() int32
	Height() int32
}

type Selectable interface {
	Widget
	Selected(value bool)
	IsSelected() bool
	SetSelectable(value bool)
	IsSelectable() bool
}
