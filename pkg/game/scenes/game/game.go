package game

import (
	"errors"

	"github.com/veandco/go-sdl2/sdl"
)

func isoToCartesian(isometricPos sdl.Point, tileSize sdl.Rect) sdl.Point {
	return sdl.Point{
		X: (isometricPos.Y/(tileSize.H/2) + isometricPos.X/(tileSize.W/2)) / 2,
		Y: (isometricPos.Y/(tileSize.H/2) - isometricPos.X/(tileSize.W/2)) / 2,
	}
}

func cartesianToIso(cartesianPos sdl.Point, tileSize sdl.Rect) sdl.Point {
	return sdl.Point{
		X: (cartesianPos.X - cartesianPos.Y) * (tileSize.W / 2),
		Y: (cartesianPos.X + cartesianPos.Y) * (tileSize.H / 4),
	}
}

func (s *GameScene) flagTile(col, row int32) error {
	if s.grid.tiles[col][row].has(tileStateShown) {
		return errors.New("can't flag: tile already opened")
	}
	if s.grid.tiles[col][row].has(tileStateFlagged) {
		s.grid.tiles[col][row].unset(tileStateFlagged)
	} else {
		s.grid.tiles[col][row].set(tileStateFlagged)
	}
	return nil
}

func (s *GameScene) openTile(col, row int32, count *int) error {
	tile := &s.grid.tiles[col][row]
	if tile.has(tileStateFlagged) {
		return errors.New("can't open: tile flagged")
	}
	if tile.has(tileStateShown) {
		return errors.New("can't open: tile already opened")
	}
	tile.set(tileStateShown)
	if tile.has(tileStateBomb) {
		tile.set(tileStateExploded)
	}
	if count != nil {
		*count += 1
	}
	if tile.bombAround == 0 && !tile.has(tileStateBomb) {
		for _, pos := range s.grid.tilesAround(col, row) {
			s.openTile(pos.X, pos.Y, count)
		}
	}
	return nil
}

func (s *GameScene) moveTo(col, row int32) error {
	if col < 0 || col > int32(s.grid.columns)-1 || row < 0 && row > int32(s.grid.rows)-1 {
		return errors.New("can't move: invalid position")
	}
	if s.grid.tiles[col][row].has(tileStateBorder) {
		return errors.New("can't move: border reached")
	}
	s.player.pos.col = col
	s.player.pos.row = row
	return nil
}
