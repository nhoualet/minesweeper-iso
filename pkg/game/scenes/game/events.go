package game

import (
	"fmt"
	"minesweeper/pkg/game/rendering"
	"minesweeper/pkg/game/scenes"

	"github.com/veandco/go-sdl2/sdl"
)

func (s *GameScene) ProcessEvent(e sdl.Event) scenes.EventState {
	switch t := e.(type) {
	case *sdl.MouseButtonEvent:
		if t.State == sdl.PRESSED {
			if (eventMouseClick(s, sdl.Point{X: t.X, Y: t.Y})) {
				return scenes.EventProcessed
			}
		}

	case *sdl.KeyboardEvent:
		keyCode := t.Keysym.Sym
		pressed := t.State == sdl.PRESSED
		if pressed {
			if keyCode == sdl.K_ESCAPE {
				s.Exit()
				return scenes.EventProcessed
			} else if keyCode == s.keyConfig.KeyUp {
				eventMoveUP(s)
			} else if keyCode == s.keyConfig.KeyDown {
				eventMoveDOWN(s)
			} else if keyCode == s.keyConfig.KeyRight {
				eventMoveRIGHT(s)
			} else if keyCode == s.keyConfig.KeyLeft {
				eventMoveLEFT(s)
			} else if s.state == gameStatePlaying {
				if keyCode == s.keyConfig.KeyOpen && s.state == gameStatePlaying {
					eventOpenTile(s)
				} else if keyCode == s.keyConfig.KeyFlag && s.state == gameStatePlaying {
					eventToggleFlag(s)
				}
			} else {
				if keyCode == s.keyConfig.KeyReplay {
					s.load()
					s.needsRedraw = true
				}
			}
		}
	}
	return scenes.EventToProcess
}

// returns true if the click was on a button
func eventMouseClick(s *GameScene, pos sdl.Point) bool {
	for _, w := range s.widgets {
		if btn, ok := w.(*rendering.Button); ok {
			if btn.OnButton(pos) {
				s.processButtonClick(btn)
				return true
			}
		}
	}
	return false
}

func eventMoveLEFT(s *GameScene) {
	err := s.moveTo(s.player.pos.col, s.player.pos.row+1)
	if err == nil {
		if s.state == gameStatePlaying {
			s.updateStateMessage(fmt.Sprintf("Moved to tile @%d;%d", s.player.pos.col, s.player.pos.row))
		}
		s.needsRedraw = true
	}
}
func eventMoveRIGHT(s *GameScene) {
	err := s.moveTo(s.player.pos.col, s.player.pos.row-1)
	if err == nil {
		if s.state == gameStatePlaying {
			s.updateStateMessage(fmt.Sprintf("Moved to tile @%d;%d", s.player.pos.col, s.player.pos.row))
		}
		s.needsRedraw = true
	}
}
func eventMoveUP(s *GameScene) {
	err := s.moveTo(s.player.pos.col-1, s.player.pos.row)
	if err == nil {
		if s.state == gameStatePlaying {
			s.updateStateMessage(fmt.Sprintf("Moved to tile @%d;%d", s.player.pos.col, s.player.pos.row))
		}
		s.needsRedraw = true
	}
}
func eventMoveDOWN(s *GameScene) {
	err := s.moveTo(s.player.pos.col+1, s.player.pos.row)
	if err == nil {
		if s.state == gameStatePlaying {
			s.updateStateMessage(fmt.Sprintf("Moved to tile @%d;%d", s.player.pos.col, s.player.pos.row))
		}
		s.needsRedraw = true
	}
}

func eventOpenTile(s *GameScene) {
	count := 0
	err := s.openTile(s.player.pos.col, s.player.pos.row, &count)
	if err == nil {
		s.stats.tilesHidden -= count
		if s.grid.tiles[s.player.pos.col][s.player.pos.row].has(tileStateBomb) {
			s.stats.bombsRemaining -= 1
			s.stats.bombsExploded += 1
			if s.stats.totalLives >= 0 {
				s.stats.livesRemaining -= 1
			}
			s.updateStateMessage(fmt.Sprintf("Bomb exploded @%d:%d", s.player.pos.col, s.player.pos.row))
		} else if count == 1 {
			s.updateStateMessage(fmt.Sprintf("opened tile @%d:%d", s.player.pos.col, s.player.pos.row))
		} else {
			s.updateStateMessage(fmt.Sprintf("opened %d tiles from @%d:%d", count, s.player.pos.col, s.player.pos.row))
		}
		s.checkGameState()
		s.needsRedraw = true
	}
}

func eventToggleFlag(s *GameScene) {
	err := s.flagTile(s.player.pos.col, s.player.pos.row)
	if err == nil {
		if s.grid.tiles[s.player.pos.col][s.player.pos.row].has(tileStateFlagged) {
			if s.partyGameConfig.WrongFlagPenalty && !s.grid.tiles[s.player.pos.col][s.player.pos.row].has(tileStateBomb) {
				s.flagTile(s.player.pos.col, s.player.pos.row)
				s.openTile(s.player.pos.col, s.player.pos.row, nil)
				s.stats.livesRemaining -= 1
				s.updateStateMessage(fmt.Sprintf("Wrong flag set on tile @%d:%d", s.player.pos.col, s.player.pos.row))

			} else {
				s.stats.flagsUsed += 1
				s.updateStateMessage(fmt.Sprintf("Set flag @%d:%d", s.player.pos.col, s.player.pos.row))
			}
		} else {
			s.stats.flagsUsed -= 1
			s.updateStateMessage(fmt.Sprintf("Unset flag @%d:%d", s.player.pos.col, s.player.pos.row))
		}
		s.checkGameState()
		s.needsRedraw = true
	}
}
