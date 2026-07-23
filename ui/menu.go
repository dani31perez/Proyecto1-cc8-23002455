package ui

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text"
	"golang.org/x/image/font/basicfont"
)

type Menu struct {

	manager *Manager

	server Button
	client Button
}

func NewMenu() *Menu {

	m := &Menu{}

	m.server = Button{
		X: 520,
		Y: 250,
		W: 240,
		H: 60,
		Text: "Servidor",
	}

	m.client = Button{
		X: 520,
		Y: 350,
		W: 240,
		H: 60,
		Text: "Cliente",
	}

	return m
}

func (m *Menu) Update() error {

	m.server.OnClick = func() {

		screen := NewServer()

		screen.manager = m.manager

		m.manager.Set(screen)
	}

	m.client.OnClick = func() {

		screen := NewClient()

		screen.manager = m.manager

		m.manager.Set(screen)
	}

	m.server.Update()
	m.client.Update()

	return nil
}

func (m *Menu) Draw(screen *ebiten.Image) {

	text.Draw(
		screen,
		"Capture The Flag",
		basicfont.Face7x13,
		500,
		150,
		nil,
	)

	m.server.Draw(screen)
	m.client.Draw(screen)
}