package ui

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text"
	"golang.org/x/image/font/basicfont"
)

type ClientScreen struct {

	manager *Manager

	back Button
}

func NewClient() *ClientScreen {

	c := &ClientScreen{}

	c.back = Button{
		X: 20,
		Y: 20,
		W: 180,
		H: 50,
		Text: "Regresar",
	}

	return c
}

func (c *ClientScreen) Update() error {

	if ebiten.IsKeyPressed(ebiten.KeyEscape) {

		menu := NewMenu()

		menu.manager = c.manager

		c.manager.Set(menu)
	}

	c.back.OnClick = func() {

		menu := NewMenu()

		menu.manager = c.manager

		c.manager.Set(menu)
	}

	c.back.Update()

	return nil
}

func (c *ClientScreen) Draw(screen *ebiten.Image) {

	text.Draw(
		screen,
		"Modo Cliente",
		basicfont.Face7x13,
		500,
		120,
		nil,
	)

	text.Draw(
		screen,
		"Aqui posteriormente aparecera la busqueda de servidores.",
		basicfont.Face7x13,
		330,
		200,
		nil,
	)

	c.back.Draw(screen)
}