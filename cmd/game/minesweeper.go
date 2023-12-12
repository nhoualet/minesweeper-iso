package main

import (
	"log"
	"minesweeper/pkg/game/app"

	"github.com/veandco/go-sdl2/sdl"
)

func main() {
	if err := sdl.Init(sdl.INIT_EVERYTHING); err != nil {
		log.Fatal(err)
	}
	defer sdl.Quit()

	// rand.Seed(time.Now().UnixNano())
	g, err := app.NewProgram()
	if err != nil {
		log.Fatal(err)
	}
	defer g.Destroy()

	log.Println("starting")
	g.Start()
}
