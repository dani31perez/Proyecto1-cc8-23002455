package ui

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
  	"image/color"

	"Proyecto1-cc8-23002455/client"
	"Proyecto1-cc8-23002455/server"
)

type Menu struct {

	manager *Manager

	server Button
	client Button
}

func NewMenu() *Menu {

	m := &Menu{}

	m.server = Button{
		X: 550,
		Y: 280,
		W: 380,
		H: 80,
		Text: "Servidor",
	}

	m.client = Button{
		X: 550,
		Y: 400,
		W: 380,
		H: 80,
		Text: "Cliente",
	}

	return m
}

func (m *Menu) Update() error {

	m.server.OnClick = func() {

		go server.Run()

		screen := NewServer()
		screen.manager = m.manager

		m.manager.Set(screen)
	}

	m.client.OnClick = func() {

		go client.Run()

		screen := NewClient()
		screen.manager = m.manager

		m.manager.Set(screen)
	}

	m.server.Update()
	m.client.Update()

	return nil
}

func (m *Menu) Draw(screen *ebiten.Image) {

	op := &text.DrawOptions{}

	op.GeoM.Translate(
		400,
		120,
	)

	op.ColorScale.ScaleWithColor(
		color.RGBA{255,220,50,255},
	)

	text.Draw(
		screen,
		"Capture The Flag",
		TitleFont,
		op,
	)

	m.server.Draw(screen)
	m.client.Draw(screen)
}