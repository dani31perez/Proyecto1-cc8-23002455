package ui

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"Proyecto1-cc8-23002455/client"
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
	if inpututil.IsKeyJustPressed(ebiten.KeyF11) {
		ebiten.SetFullscreen(!ebiten.IsFullscreen())
	}
	client.Notify.Update()
	return g.manager.Update()
}

func (g *Game) Draw(screen *ebiten.Image) {
	g.manager.Draw(screen)
	client.Notify.Draw(screen)
	
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return ScreenWidth, ScreenHeight
}