package ui

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text/v2"

	"Proyecto1-cc8-23002455/client"
	"Proyecto1-cc8-23002455/ui/assets"
	"Proyecto1-cc8-23002455/ui/components"
)

type ClientScreen struct {
	manager *Manager

	back components.Button

	cards []components.ServerCard

	servers []client.DiscoveredServer

	loaded bool
}

func NewClient() *ClientScreen {

	c := &ClientScreen{}

	c.back = components.Button{
		X:    20,
		Y:    20,
		W:    180,
		H:    50,
		Text: "Regresar",
	}

	return c
}

func (c *ClientScreen) Update() error {

	if !c.loaded {
		servers, err := client.DiscoverServer()
		if err == nil {
			c.servers = servers
			y := 220.0
			for _, server := range servers {
				card := components.ServerCard{
					X:      300,
					Y:      y,
					W:      900,
					H:      90,
					Server: server,
					OnClick: func(s client.DiscoveredServer) {
						go client.Run(s)
					},
				}
				c.cards = append(
					c.cards,
					card,
				)
				y += 110
			}
		}
		c.loaded = true
	}

	for i := range c.cards {
		c.cards[i].Update()
	}

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
		assets.TitleFont,
		op,
	)

	for i := range c.cards {
		c.cards[i].Draw(screen)
	}

	c.back.Draw(screen)
}
