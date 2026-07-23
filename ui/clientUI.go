package ui

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
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

	op := &text.DrawOptions{}

	op.GeoM.Translate(
		500,
		120,
	)

	op.ColorScale.ScaleWithColor(color.White)

	text.Draw(
		screen,
		"Modo Cliente",
		TitleFont,
		op,
	)

	op2 := &text.DrawOptions{}

	op2.GeoM.Translate(
		320,
		200,
	)

	op2.ColorScale.ScaleWithColor(color.White)

	text.Draw(
		screen,
		"Aqui posteriormente aparecera la busqueda de servidores.",
		SmallFont,
		op2,
	)

	c.back.Draw(screen)
}