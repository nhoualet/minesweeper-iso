package config

import (
	"fmt"

	"github.com/veandco/go-sdl2/sdl"
)

const ConfigFilePath = "data/config.yml"

type WindowConfig struct {
	FPS           int32  `yaml:"fps"`
	Width         int32  `yaml:"width"`
	Height        int32  `yaml:"height"`
	Borderless    bool   `yaml:"borderless"`
	Fullscreen    bool   `yaml:"fullscreen"`
	IconFile      string `yaml:"icon_path"`
	FontFile      string `yaml:"font_path"`
	ResourcesPath string `yaml:"resources_path"`
}

type GameConfig struct {
	GridColumns      uint32 `yaml:"grid-column"`
	GridRows         uint32 `yaml:"grid-row"`
	BombPercent      int    `yaml:"bomb-percent"`
	Lives            int    `yaml:"lives"`
	WrongFlagPenalty bool   `yaml:"wrong_flag_penalty"`
}

type ControlNames struct {
	KeyUp     string `yaml:"up"`
	KeyDown   string `yaml:"down"`
	KeyLeft   string `yaml:"left"`
	KeyRight  string `yaml:"right"`
	KeyFlag   string `yaml:"flag"`
	KeyOpen   string `yaml:"open"`
	KeyReplay string `yaml:"replay"`
}

type ControlCodes struct {
	KeyUp     sdl.Keycode
	KeyDown   sdl.Keycode
	KeyLeft   sdl.Keycode
	KeyRight  sdl.Keycode
	KeyFlag   sdl.Keycode
	KeyOpen   sdl.Keycode
	KeyReplay sdl.Keycode
}

type GameControls struct {
	Names ControlNames `yaml:"keys"`
	Codes ControlCodes `yaml:"-"`
}

type Config struct {
	Window   WindowConfig `yaml:"window"`
	Game     GameConfig   `yaml:"game"`
	Controls GameControls `yaml:"controls"`
}

func (c *Config) Check() error {
	if c.Window.FPS < 1 {
		return fmt.Errorf("invalid FPS: got %d (expected FPS>0)", c.Window.FPS)
	}
	c.Controls.Codes.KeyUp = sdl.GetKeyFromName(c.Controls.Names.KeyUp)
	if c.Controls.Codes.KeyUp == sdl.K_UNKNOWN {
		return fmt.Errorf("unknown key name (KeyUp): %q", c.Controls.Names.KeyUp)
	}
	c.Controls.Codes.KeyDown = sdl.GetKeyFromName(c.Controls.Names.KeyDown)
	if c.Controls.Codes.KeyDown == sdl.K_UNKNOWN {
		return fmt.Errorf("unknown key name (KeyDown): %q", c.Controls.Names.KeyDown)
	}
	c.Controls.Codes.KeyLeft = sdl.GetKeyFromName(c.Controls.Names.KeyLeft)
	if c.Controls.Codes.KeyLeft == sdl.K_UNKNOWN {
		return fmt.Errorf("unknown key name (KeyLeft): %q", c.Controls.Names.KeyLeft)
	}
	c.Controls.Codes.KeyRight = sdl.GetKeyFromName(c.Controls.Names.KeyRight)
	if c.Controls.Codes.KeyRight == sdl.K_UNKNOWN {
		return fmt.Errorf("unknown key name (KeyRight): %q", c.Controls.Names.KeyRight)
	}
	c.Controls.Codes.KeyFlag = sdl.GetKeyFromName(c.Controls.Names.KeyFlag)
	if c.Controls.Codes.KeyFlag == sdl.K_UNKNOWN {
		return fmt.Errorf("unknown key name (KeyFlag): %q", c.Controls.Names.KeyFlag)
	}
	c.Controls.Codes.KeyOpen = sdl.GetKeyFromName(c.Controls.Names.KeyOpen)
	if c.Controls.Codes.KeyOpen == sdl.K_UNKNOWN {
		return fmt.Errorf("unknown key name (KeyOpen): %q", c.Controls.Names.KeyOpen)
	}
	c.Controls.Codes.KeyReplay = sdl.GetKeyFromName(c.Controls.Names.KeyReplay)
	if c.Controls.Codes.KeyReplay == sdl.K_UNKNOWN {
		return fmt.Errorf("unknown key name (KeyOpen): %q", c.Controls.Names.KeyOpen)
	}
	return nil
}

var DefaultConfig = Config{
	Window: WindowConfig{
		Width:         1080,
		Height:        720,
		FPS:           30,
		Borderless:    false,
		IconFile:      "logo.png",
		FontFile:      "PressStart2P.ttf",
		ResourcesPath: "./assets",
	},
	Game: GameConfig{
		GridColumns:      30,
		GridRows:         30,
		BombPercent:      10,
		Lives:            3,
		WrongFlagPenalty: false,
	},
	Controls: GameControls{
		Names: ControlNames{
			KeyUp:     "up",
			KeyDown:   "down",
			KeyLeft:   "left",
			KeyRight:  "right",
			KeyFlag:   "f",
			KeyOpen:   "space",
			KeyReplay: "r",
		},
		Codes: ControlCodes{
			KeyUp:     sdl.K_UNKNOWN,
			KeyDown:   sdl.K_UNKNOWN,
			KeyLeft:   sdl.K_UNKNOWN,
			KeyRight:  sdl.K_UNKNOWN,
			KeyFlag:   sdl.K_UNKNOWN,
			KeyOpen:   sdl.K_UNKNOWN,
			KeyReplay: sdl.K_UNKNOWN,
		},
	},
}
