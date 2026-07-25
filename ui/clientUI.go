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

	input client.KeyListener

	loaded bool

	ipBox      components.TextBox
    portBox    components.TextBox
    connectBtn components.Button
}

func NewClient() *ClientScreen {

	c := &ClientScreen{}

	c.ipBox = components.TextBox{
		X:    300,
		Y:    230,
		W:    320,
		H:    40,
	}

	c.portBox = components.TextBox{
		X:    640,
		Y:    230,
		W:    100,
		H:    40,
		Text: "8889",
	}

	c.connectBtn = components.Button{
		X:    760,
		Y:    230,
		W:    170,
		H:    40,
		Text: "Conectar",
	}

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
			y := 340.0
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

	c.ipBox.Update()
	c.portBox.Update()

	c.connectBtn.OnClick = func() {
		server := client.DiscoveredServer{
			IP: c.ipBox.Text,
		}

		server.TCPPort = 8889

		go client.Run(server)
	}

	c.connectBtn.Update()

	c.input.Update()
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

	op2 := &text.DrawOptions{}
	op2.GeoM.Translate(300, 170)
	op2.ColorScale.ScaleWithColor(color.White)

	text.Draw(
		screen,
		"Conexion manual",
		assets.TitleFont,
		op2,
	)

	c.ipBox.Draw(screen)
	c.portBox.Draw(screen)
	c.connectBtn.Draw(screen)

	op3 := &text.DrawOptions{}
	op3.GeoM.Translate(300, 290)
	op3.ColorScale.ScaleWithColor(color.White)

	text.Draw(
		screen,
		"Servidores encontrados",
		assets.TitleFont,
		op3,
	)

	for i := range c.cards {
		c.cards[i].Draw(screen)
	}

	c.back.Draw(screen)
}
