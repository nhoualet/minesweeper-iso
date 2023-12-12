package game

import (
	"errors"
	"math/rand"

	"github.com/veandco/go-sdl2/sdl"
)

type tile struct {
	bombAround int
	state      tileState
}

type grid struct {
	columns uint32
	rows    uint32
	tiles   [][]tile
}
type playerPos struct {
	col int32
	row int32
}
type player struct {
	pos playerPos
}

func (t *tile) set(s tileState) {
	t.state = t.state | s
}
func (t *tile) unset(s tileState) {
	t.state = t.state & (^s)
}
func (t *tile) has(s tileState) bool {
	// fmt.Printf("%8b & %8b = %8b\n", t.state, s, t.state&s)
	return (t.state & s) != 0
}

func (g *grid) tile(col, row int32) *tile {
	if row > int32(g.columns) || row > int32(g.rows) {
		return nil
	}
	return &g.tiles[col][row]
}

func (g grid) tilesAround(col, row int32) []sdl.Point {
	result := []sdl.Point{}
	if row > 0 {
		if col > 0 {
			result = append(result, sdl.Point{X: col - 1, Y: row - 1}) // top left
		}
		result = append(result, sdl.Point{X: col, Y: row - 1}) // top
		if col < int32(g.columns-1) {
			result = append(result, sdl.Point{X: col + 1, Y: row - 1}) // top right
		}
	}
	if col < int32(g.columns)-1 {
		result = append(result, sdl.Point{X: col + 1, Y: row}) // right
	}
	if row < int32(g.rows)-1 {
		if col < int32(g.columns)-1 {
			result = append(result, sdl.Point{X: col + 1, Y: row + 1}) // bottom right
		}
		result = append(result, sdl.Point{X: col, Y: row + 1}) // bottom
		if col > 0 {
			result = append(result, sdl.Point{X: col - 1, Y: row + 1}) // bottom left
		}
	}
	if col > 0 {
		result = append(result, sdl.Point{X: col - 1, Y: row}) // left
	}
	return result
}

func (g *grid) placeBomb(col, row int32) error {
	if g.tiles[col][row].has(tileStateBomb) {
		return errors.New("tile already a bomb")
	}
	if g.tiles[col][row].has(tileStateBorder) {
		return errors.New("tile a border")
	}
	g.tiles[col][row].set(tileStateBomb)
	for _, pos := range g.tilesAround(col, row) {
		if !g.tiles[pos.X][pos.Y].has(tileStateBomb | tileStateBorder) {
			g.tiles[pos.X][pos.Y].bombAround += 1
		}
	}
	return nil
}

func newGrid(col, row, bombCount uint32) *grid {
	col += 2
	row += 2
	tiles := make([][]tile, col)
	for i := range tiles {
		tiles[i] = make([]tile, row)
	}
	g := &grid{
		columns: col,
		rows:    row,
		tiles:   tiles,
	}
	for _, col := range g.tiles {
		for _, t := range col {
			t.state = 0
			t.bombAround = 0
		}
	}
	for i := 0; i < int(g.columns); i += 1 {
		g.tiles[i][0].set(tileStateShown | tileStateBorder)
		g.tiles[i][g.rows-1].set(tileStateShown | tileStateBorder)
	}

	for i := 0; i < int(g.rows); i += 1 {
		g.tiles[0][i].set(tileStateShown | tileStateBorder)
		g.tiles[g.columns-1][i].set(tileStateShown | tileStateBorder)
	}
	for bombCount > 0 {
		col := int32(rand.Intn(int(g.columns) - 1))
		row := int32(rand.Intn(int(g.rows) - 1))
		err := g.placeBomb(col, row)
		if err == nil {
			bombCount -= 1
		}
	}
	return g
}
