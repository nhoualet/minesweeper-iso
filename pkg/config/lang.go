package config

type MainMenuLang struct {
	NewGame    string `yaml:"new-game"`
	ResumeGame string `yaml:"resume-game"`
	Settings   string `yaml:"settings"`
	Exit       string `yaml:"exit"`
}

type SettingsLang struct {
	Columns string `yaml:"columns"`
	Rows    string `yaml:"rows"`
	Bombs   string `yaml:"bombs"`
	Lives   string `yaml:"lives"`
	// LivesNoLimit string `yaml:"livesNoLimit"`
}

type GameLang struct {
	TilesHidden    string `yaml:"tiles-hidden"`
	FlagsUsed      string `yaml:"flags-used"`
	bombsRemaining string `yaml:"bombs-remaining"`
	livesLeft      string `yaml:"lives-left"`

	MovedTo             string `yaml:"moved-to"`
	BombExploded        string `yaml:"bomb-exploded"`
	OpenedSingleTile    string `yaml:"opened-single-tile"`
	OpenedMultipleTiles string `yaml:"opened-multiple-tiles"`
	FlagSetWrong        string `yaml:"flag-set-wrong"`
	FlagSet             string `yaml:"flag-set-unspecified"`
	FlagUnset           string `yaml:"flag-unset"`

	GameLost       string `yaml:"game-lost"`
	GameWon        string `yaml:"game-won"`
	IncorrectFlags string `yaml:"incorrect-flags"`

	MainMenuBtn string `yaml:"main-menu"`
	SettingsBtn string `yaml:"settings-menu"`
}

type LangConfig struct {
	MainMenu MainMenuLang `yaml:"main-menu"`
	Settings SettingsLang `yaml:"settings-menu"`
	Game     GameLang     `yaml:"game"`
}

var DefaultLang = LangConfig{
	MainMenu: MainMenuLang{
		NewGame:    "New Game",
		ResumeGame: "Continue",
		Settings:   "Exit",
	},
	Settings: SettingsLang{
		Columns: "Grid columns: %d",
		Rows:    "Grid rows: %d",
		Bombs:   "Bombs: %d (%d of the tiles)",
		Lives:   "Lives: %d",
	},
	Game: GameLang{
		TilesHidden:         "Tiles hidden: %d/%d",
		FlagsUsed:           "Flags used: %d",
		bombsRemaining:      "Bombs remaining: %d",
		livesLeft:           "Lives left: %d/%d",
		MovedTo:             "Moved to tile @%d;%d",
		BombExploded:        "Bomb exploded @%d;%d",
		OpenedSingleTile:    "Opened tile @%d;%d",
		OpenedMultipleTiles: "Opened %d tiles from @%d;%d",
		FlagSetWrong:        "Wrong flag set on tile @%d;%d",
		FlagSet:             "Set flag @%d;%d",
		FlagUnset:           "Unset flag @%d;%d",
		GameLost:            "Game lost (no lives left) press [%s] to replay",
		GameWon:             "Game won, press [%s] to replay",
		IncorrectFlags:      "At least one flag isn't right",
		MainMenuBtn:         "Main menu",
		SettingsBtn:         "Settings",
	},
}

func (c *LangConfig) Check() error {
	return nil
}
