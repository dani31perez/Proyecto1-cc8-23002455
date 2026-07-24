package ui

import (
	"github.com/hajimehoshi/ebiten/v2"
)
const (
	ScreenWidth  = 1480
	ScreenHeight = 900
)

type Game struct {
	manager *Manager
}

func NewGame() *Game {

	menu := NewMenu()

	manager := NewManager(menu)

	menu.manager = manager

	return &Game{
		manager: manager,
	}
}

func (g *Game) Update() error {
	return g.manager.Update()
}

func (g *Game) Draw(screen *ebiten.Image) {
	g.manager.Draw(screen)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return ScreenWidth, ScreenHeight
}