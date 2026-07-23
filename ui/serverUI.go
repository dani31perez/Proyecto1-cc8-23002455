package ui

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
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
		W: 220,
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

	op := &text.DrawOptions{}

	op.GeoM.Translate(
		500,
		120,
	)

	op.ColorScale.ScaleWithColor(
		color.RGBA{80,220,80,255},
	)

	text.Draw(
		screen,
		"Modo Servidor",
		TitleFont,
		op,
	)

	op2 := &text.DrawOptions{}

	op2.GeoM.Translate(
		460,
		200,
	)

	op2.ColorScale.ScaleWithColor(
		color.RGBA{80,170,255,255},
	)

	text.Draw(
		screen,
		"Aqui posteriormente aparecera el lobby.",
		SmallFont,
		op2,
	)

	s.back.Draw(screen)
}