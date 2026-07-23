package ui

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text"
	"golang.org/x/image/font/basicfont"
)

type ServerScreen struct {

	manager *Manager

	back Button
}

func NewServer() *ServerScreen {

	s := &ServerScreen{}

	s.back = Button{
		X: 20,
		Y: 20,
		W: 180,
		H: 50,
		Text: "Regresar",
	}

	return s
}

func (s *ServerScreen) Update() error {

	if ebiten.IsKeyPressed(ebiten.KeyEscape) {
		menu := NewMenu()
		menu.manager = s.manager
		s.manager.Set(menu)
	}

	s.back.OnClick = func() {

		menu := NewMenu()

		menu.manager = s.manager

		s.manager.Set(menu)
	}

	s.back.Update()

	return nil
}

func (s *ServerScreen) Draw(screen *ebiten.Image) {

	text.Draw(
		screen,
		"Modo Servidor",
		basicfont.Face7x13,
		500,
		120,
		nil,
	)

	text.Draw(
		screen,
		"Aqui posteriormente aparecera el lobby.",
		basicfont.Face7x13,
		400,
		200,
		nil,
	)

	s.back.Draw(screen)
}