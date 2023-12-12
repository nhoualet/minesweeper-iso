package app

import (
	"errors"
	"fmt"
	"minesweeper/pkg/config"
	"minesweeper/pkg/game/rendering"
	"minesweeper/pkg/game/scenes"
	"minesweeper/pkg/game/scenes/game"
	"minesweeper/pkg/game/scenes/menuMain"
	"minesweeper/pkg/game/scenes/menuSettings"
	"strings"

	"github.com/veandco/go-sdl2/img"
	"github.com/veandco/go-sdl2/sdl"
	"github.com/veandco/go-sdl2/ttf"
)

type Program struct {
	fps          int32
	timePerFrame uint64
	config       config.Config
	language     config.LangConfig
	renderer     *rendering.CustomRenderer
	icon         *sdl.Surface
	font         *ttf.Font
	isRunning    bool
	sceneManager *scenes.SceneManager
}

/*
Generates a new program with the window/renderer and all the scenes
*/
func NewProgram() (*Program, error) {
	var cfg config.Config
	err := config.LoadConfig(config.ConfigFilePath, &cfg)
	if err == nil {
		err = cfg.Check()
		if err != nil {
			fmt.Printf("Invalid config: %s (fallback to defaults)\n", err)
		}
	}
	if err != nil {
		cfg = config.DefaultConfig
		err = cfg.Check()
		if err != nil {
			return nil, fmt.Errorf("couldn't load default config: %s", err)
		}
	}
	if !strings.HasSuffix(cfg.Window.ResourcesPath, "/") {
		cfg.Window.ResourcesPath += "/"
	}
	fmt.Printf("%+v\n", cfg)
	if cfg.Window.FPS < 1 {
		return nil, errors.New("program load: fps < 0")
	}
	windowFlags := sdl.WINDOW_SHOWN | sdl.WINDOW_RESIZABLE
	if cfg.Window.Borderless {
		windowFlags = windowFlags | sdl.WINDOW_BORDERLESS
	}
	if cfg.Window.Fullscreen {
		windowFlags = windowFlags | sdl.WINDOW_FULLSCREEN_DESKTOP
	}

	var langCfg config.LangConfig
	err = config.LoadConfig("data/lang.yml", &cfg)
	if err == nil {
		err = langCfg.Check()
		if err != nil {
			fmt.Printf("Invalid config: %s (fallback to defaults)\n", err)
		}
	}
	if err != nil {
		langCfg = config.DefaultLang
		/*
			err = config.CheckConfig(&langCfg)
			 if err != nil {
				return nil, fmt.Errorf("couldn't load default config: %s", err)
			}
		*/
	}
	fmt.Printf("%+v\n", langCfg)

	customRenderer, err := rendering.CreateCustomRenderer(
		"13Noodles' Isometric minesweeper",
		sdl.WINDOWPOS_UNDEFINED, sdl.WINDOWPOS_UNDEFINED,
		int32(cfg.Window.Width), int32(cfg.Window.Height),
		uint32(windowFlags),
		-1,
		sdl.RENDERER_ACCELERATED|sdl.RENDERER_PRESENTVSYNC,
	)
	if err != nil {
		return nil, fmt.Errorf("customRenderer init: %s", err)
	}

	icon, err := img.Load(cfg.Window.ResourcesPath + cfg.Window.IconFile)
	if err != nil {
		return nil, fmt.Errorf("icon load: %s", err)
	}
	customRenderer.SDLwindow.SetIcon(icon)
	ttf.Init()
	font, err := ttf.OpenFont(cfg.Window.ResourcesPath+cfg.Window.FontFile, 12)
	if err != nil {
		return nil, fmt.Errorf("font load: %s", err)
	}

	program := &Program{
		fps:          cfg.Window.FPS,
		timePerFrame: uint64(1000 / cfg.Window.FPS),
		config:       cfg,
		language:     langCfg,
		renderer:     customRenderer,
		icon:         icon,
		font:         font,
		isRunning:    true,
	}
	sceneManager, err := scenes.NewSceneManager(&program.isRunning, *program.renderer, cfg)
	if err != nil {
		return nil, fmt.Errorf("sceneManager load: %s", err)
	}
	default_scene, err := menuMain.Initialize(sceneManager, customRenderer, font, langCfg.MainMenu)
	if err != nil {
		return nil, fmt.Errorf("default scene load: %s", err)
	}
	settingsScene, err := menuSettings.Initialize(sceneManager, customRenderer, font)
	if err != nil {
		return nil, fmt.Errorf("scene load: %s", err)
	}
	gameScene, err := game.Initialize(sceneManager, customRenderer, font)
	if err != nil {
		return nil, fmt.Errorf("scene load: %s", err)
	}

	sceneManager.AddDefaultScene(default_scene, "main")
	sceneManager.AddScene(settingsScene, "settings")
	sceneManager.AddScene(gameScene, "game")
	sceneManager.SetScene("main", true)
	program.sceneManager = sceneManager
	return program, nil
}

// Cleanup the Program and call SceneManager's destroy
func (p *Program) Destroy() {
	if p.icon != nil {
		p.icon.Free()
	}
	if p.font != nil {
		p.font.Close()
	}
	if p.renderer != nil {
		p.renderer.Destroy()
	}
}

// Starts the program's loop, calling the current Scene's ProcessEvent(), Update() and Draw() functions
func (p *Program) Start() {
	var startTime uint64
	var deltaMS uint64 = 0
	endTime := sdl.GetTicks64() - p.timePerFrame

	p.isRunning = true
	for p.isRunning {
		p.processEvents()
		startTime = sdl.GetTicks64()
		deltaMS = startTime - endTime
		if deltaMS > p.timePerFrame {
			p.update(deltaMS)
			p.draw()
			endTime = startTime
		} else {
			sdl.Delay(uint32(p.timePerFrame - deltaMS))
		}
	}
}

func (p *Program) processEvents() {
	for event := sdl.PollEvent(); event != nil; event = sdl.PollEvent() {
		switch event.(type) {
		case *sdl.QuitEvent:
			p.isRunning = false
			continue
		}
		p.sceneManager.ProcessEvent(event)
	}
}

func (p *Program) update(deltaMS uint64) {
	p.sceneManager.Update(deltaMS)
}

func (p *Program) draw() {
	p.sceneManager.Draw(*p.renderer)
}
