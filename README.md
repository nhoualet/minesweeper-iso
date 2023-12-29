# Minesweeper isometric

## How to build

`go build cmd/game/minesweeper.go` (SDL2 must be installed)

## Configs
The configs can be changed in the file data/config.yml (Window size, keymaps, grid size, lives, amount of bombs, ...)

The keymaps names can be any [SDL2 Keycode value](https://wiki.libsdl.org/SDL2/SDL_Keycode)

Some configs can also be change ingame in the Settings.

## Default configs
keymaps:
- go up: ⬆️
- go down: ⬇️
- go left: ⬅️
- go right: ➡️
- toggle flag: F
- open a tile: SPACE
- replay after the game's end: R

window:
- dimension: 1080x720
- bordered

game:
- 10x10 grid
- 10% of the tiles are bombs
- infinite lives
- no penalty when a flag is wrong


## Screenshots
![Main menu screenshot](../gh-pages/images/mainMenu.png)
![Game screenshot](../gh-pages/images/playing.png)
![Game won screenshot](../gh-pages/images/gameWon.png)
