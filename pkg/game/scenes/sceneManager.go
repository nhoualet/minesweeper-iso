package scenes

import (
	"errors"
	"fmt"
	"minesweeper/pkg/config"
	"minesweeper/pkg/game/rendering"

	"github.com/veandco/go-sdl2/sdl"
)

type SceneManager struct {
	currentScene      interface{ Scene }
	namePreviousScene string
	nameCurrentScene  string
	scenes            map[string]Scene
	defaultSceneName  string
	IsRunning         *bool
	config            config.Config
	renderer          rendering.CustomRenderer
}

func (sm *SceneManager) AddScene(scene Scene, id string) {
	sm.scenes[id] = scene
}
func (sm *SceneManager) AddDefaultScene(scene Scene, id string) {
	sm.scenes[id] = scene
	sm.defaultSceneName = id
}

func NewSceneManager(isRunning *bool, renderer rendering.CustomRenderer, config config.Config) (*SceneManager, error) {
	sm := &SceneManager{
		scenes:            map[string]Scene{},
		currentScene:      nil,
		namePreviousScene: "",
		defaultSceneName:  "",
		IsRunning:         isRunning,
		config:            config,
		renderer:          renderer,
	}
	return sm, nil
}

func (sm *SceneManager) GetConfig() config.Config {
	return sm.config
}

func (sm *SceneManager) SetConfig(cfg config.Config) {
	sm.config = cfg
}

func (sm *SceneManager) GetScene(name string) (Scene, error) {
	scene, found := sm.scenes[name]
	if !found {
		return nil, errors.New("getScene: scene not found")
	}
	return scene, nil
}

func (sm *SceneManager) SetScene(name string, reload bool) error {
	newScene, found := sm.scenes[name]
	if !found {
		return errors.New("setScene: scene not found")
	}
	if sm.currentScene != nil {
		fmt.Printf("%q unload\n", sm.nameCurrentScene)
		sm.currentScene.Unload()
		sm.namePreviousScene = sm.nameCurrentScene
	}
	sm.currentScene = newScene
	sm.nameCurrentScene = name
	fmt.Printf("%q load\n", sm.nameCurrentScene)
	sm.currentScene.Enter(reload)
	return nil
}
func (sm *SceneManager) SetLastScene(reload bool) error {
	return sm.SetScene(sm.namePreviousScene, reload)
}

func (sm *SceneManager) SetSceneDefault(reload bool) error {
	return sm.SetScene(sm.defaultSceneName, reload)
}
func (sm *SceneManager) ExitCurrentScene() {
	if sm.currentScene != nil {
		sm.currentScene.Exit()
	}
}

func (sm *SceneManager) Quit() {
	fmt.Println("bye bye")
	*sm.IsRunning = false
}

func (sm *SceneManager) ProcessEvent(e sdl.Event) {
	if sm.currentScene == nil {
		sm.Quit()
		return
	}
	if sm.currentScene.ProcessEvent(e) == EventProcessed {
		return
	}
	switch t := e.(type) {
	case *sdl.KeyboardEvent:
		keyCode := t.Keysym.Sym
		pressed := t.State == sdl.PRESSED
		if pressed {
			if keyCode == sdl.K_ESCAPE {
				if sm.nameCurrentScene == sm.defaultSceneName {
					sm.Quit()
				}
			}
		}

	case *sdl.WindowEvent:
		if t.Event == sdl.WINDOWEVENT_SIZE_CHANGED {
			sm.currentScene.ProcessResize(t.Data1, t.Data2)
		}
	}

}
func (sm *SceneManager) Update(deltaMS uint64) {
	if sm.currentScene == nil {
		sm.Quit()
		return
	}
	sm.currentScene.Update(deltaMS)
}
func (sm *SceneManager) Draw(renderer rendering.CustomRenderer) {
	if sm.currentScene == nil {
		sm.Quit()
		return
	}
	if sm.currentScene.NeedsRedraw() {
		sm.renderer.SDLrenderer.SetDrawColor(20, 20, 20, sdl.ALPHA_OPAQUE)
		sm.renderer.SDLrenderer.Clear()
		sm.currentScene.Draw(renderer)
		sm.renderer.SDLrenderer.Present()
	}
}
