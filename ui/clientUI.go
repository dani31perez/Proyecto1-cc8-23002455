package ui

import (
	"image/color"
	"strconv"

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

	nameBox           components.TextBox
	ipBox             components.TextBox
	portBox           components.TextBox
	connectBtn        components.Button
	reloadBtn         components.Button
	scroll            float64
}

func NewClient() *ClientScreen {

	c := &ClientScreen{}

	c.nameBox = components.TextBox{
		X: 300,
		Y: 230,
		W: 400,
		H: 40,
		Text: "dani",
	}

	c.ipBox = components.TextBox{
		X: 300,
		Y: 330,
		W: 400,
		H: 40,
	}

	c.portBox = components.TextBox{
		X:    720,
		Y:    330,
		W:    100,
		H:    40,
		Text: "8889",
	}

	c.connectBtn = components.Button{
		X:    840,
		Y:    330,
		W:    250,
		H:    50,
		Text: "IP:puerto",
	}

	c.reloadBtn = components.Button{
		X:    820,
		Y:    420,
		W:    500,
		H:    50,
		Text: "Recargar servidores",
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

func (c *ClientScreen) playerName() string {
	if c.nameBox.Text == "" {
		return "Jugador"
	}
	return c.nameBox.Text
}

func (c *ClientScreen) refreshServers() {
	servers, err := client.DiscoverServer()
	if err != nil {
		return
	}
	c.setServers(servers)
}

func (c *ClientScreen) setServers(servers []client.DiscoveredServer) {
	c.cards = nil
	c.servers = servers
	y := 500.0
	for _, srv := range servers {
		server := srv
		card := components.ServerCard{
			X:      300,
			Y:      y,
			W:      900,
			H:      90,
			Server: server,
			OnClick: func(s client.DiscoveredServer) {
				go client.Run(s, c.playerName())
			},
		}
		c.cards = append(
			c.cards,
			card,
		)
		y += 110
	}
}

func (c *ClientScreen) Update() error {

	_, wheelY := ebiten.Wheel()

	c.scroll -= wheelY * 40

	if c.scroll < 0 {
		c.scroll = 0
	}

	maxScroll := float64(len(c.cards))*110 - 250
	if maxScroll < 0 {
		maxScroll = 0
	}

	if c.scroll > maxScroll {
		c.scroll = maxScroll
	}
	if client.CurrentState != nil && client.CurrentState.Started() {
		play := NewPlay("client")
		play.manager = c.manager
		c.manager.Set(play)
		return nil
	}

	if !c.loaded {
		c.refreshServers()
		c.loaded = true
	}

	for i := range c.cards {
		c.cards[i].Update()
	}

	c.nameBox.Update()
	c.ipBox.Update()
	c.portBox.Update()

	c.connectBtn.OnClick = func() {
		server, err := client.DirectServer(c.ipBox.Text, c.portBox.Text)
		if err == nil {
			go client.Run(server, c.playerName())
		}
	}

	c.connectBtn.Update()

	c.reloadBtn.OnClick = func() {
		c.refreshServers()
	}

	c.reloadBtn.Update()

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

	op.ColorScale.ScaleWithColor(
		color.RGBA{80, 220, 80, 255},
	)

	text.Draw(
		screen,
		"Modo Cliente",
		assets.TitleFont,
		op,
	)

	opn := &text.DrawOptions{}
	opn.GeoM.Translate(300, 200)
	opn.ColorScale.ScaleWithColor(color.White)

	text.Draw(
		screen,
		"Tu nombre",
		assets.MenuFont,
		opn,
	)

	c.nameBox.Draw(screen)

	op2 := &text.DrawOptions{}
	op2.GeoM.Translate(300, 300)
	op2.ColorScale.ScaleWithColor(color.White)

	text.Draw(
		screen,
		"Manual: IP para UDP; IP:puerto para TCP",
		assets.MenuFont,
		op2,
	)

	c.ipBox.Draw(screen)
	c.portBox.Draw(screen)
	c.connectBtn.Draw(screen)

	op3 := &text.DrawOptions{}
	op3.GeoM.Translate(300, 440)
	op3.ColorScale.ScaleWithColor(color.White)

	text.Draw(
		screen,
		"Servidores encontrados",
		assets.MenuFont,
		op3,
	)

	c.reloadBtn.Draw(screen)

	for i := range c.cards {

		card := c.cards[i]

		card.Y -= c.scroll

		if card.Y+card.H < 490 {
			continue
		}

		if card.Y > ScreenHeight-80 {
			continue
		}

		card.Draw(screen)
	}

	c.back.Draw(screen)

	if client.CurrentState != nil {
		op4 := &text.DrawOptions{}
		op4.GeoM.Translate(150, ScreenHeight-100)
		op4.ColorScale.ScaleWithColor(color.RGBA{120, 230, 140, 255})
		text.Draw(screen, "Conectado, esperando a que el servidor inicie la partida", assets.MenuFont, op4)

		y := float64(ScreenHeight - 70)
		for _, lp := range client.CurrentState.LobbyPlayers() {
			opp := &text.DrawOptions{}
			opp.GeoM.Translate(150, y)
			opp.ColorScale.ScaleWithColor(color.White)
			text.Draw(screen, lp.Name+"  ("+lp.ID+")", assets.SmallFont, opp)
			y += 24
		}

		if seconds := client.CurrentState.Countdown(); seconds > 0 {
			client.Notify.Show("Inicia en", strconv.Itoa(seconds), 1)
		}
	}
}
