package ui

import (
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"

	"Proyecto1-cc8-23002455/client"
	"Proyecto1-cc8-23002455/server"
	"Proyecto1-cc8-23002455/shared"
	"Proyecto1-cc8-23002455/ui/assets"
	"Proyecto1-cc8-23002455/ui/components"
)

const mapAreaSize = 760.0

type PlayScreen struct {
	manager   *Manager
	role      string
	input     client.KeyListener
	names     map[string]string
	backLobby components.Button
}

func NewPlay(role string) *PlayScreen {
	p := &PlayScreen{
		role:  role,
		names: make(map[string]string),
	}

	p.backLobby = components.Button{
		X:    ScreenWidth/2 - 130,
		Y:    ScreenHeight - 80,
		W:    350,
		H:    50,
		Text: "Volver al lobby",
	}

	return p
}
func (p *PlayScreen) mapOrigin() (float64, float64) {
	x := (ScreenWidth - mapAreaSize) / 2
	y := (ScreenHeight-mapAreaSize)/2 + 30
	return x, y
}
func (p *PlayScreen) toScreen(lx, ly float64) (float64, float64) {
	ox, oy := p.mapOrigin()
	scale := mapAreaSize / float64(shared.DefaultGameConfig.MapSize)
	return ox + lx*scale, oy + ly*scale
}
func (p *PlayScreen) Update() error {
	if p.role == "client" && client.CurrentState != nil && !client.CurrentState.Started() {
		clientScreen := NewClient()
		clientScreen.manager = p.manager
		clientScreen.loaded = true
		p.manager.Set(clientScreen)
		return nil
	}
	if p.role == "server" && server.CurrentLobby != nil && !server.CurrentLobby.IsPlaying() && !server.CurrentLobby.IsFinished() {
		serverScreen := NewServer()
		serverScreen.manager = p.manager
		p.manager.Set(serverScreen)
		return nil
	}
	if p.role == "client" && client.CurrentState != nil {
		p.input.Update()
		for _, lp := range client.CurrentState.LobbyPlayers() {
			p.names[lp.ID] = lp.Name
		}
	}
	if p.role == "server" && server.CurrentLobby != nil {
		if server.CurrentLobby.IsFinished() {
			p.backLobby.OnClick = func() {
				server.CurrentLobby.ResetLobby()
				server.CurrentLobby.BroadcastLobby()
				serverScreen := NewServer()
				serverScreen.manager = p.manager
				p.manager.Set(serverScreen)
			}
			p.backLobby.Update()
		}
		for _, sp := range server.CurrentLobby.GetPlayers() {
			p.names[sp.Id] = sp.Name
		}
	}

	if p.role == "client" &&
		client.CurrentState != nil &&
		client.CurrentState.Winner() != "" {

		p.backLobby.OnClick = func() {

			client.CurrentState.Reset()

			clientScreen := NewClient()
			clientScreen.manager = p.manager

			p.manager.Set(clientScreen)
		}

		p.backLobby.Update()
	}
	if ebiten.IsKeyPressed(ebiten.KeyEscape) {
		menu := NewMenu()
		menu.manager = p.manager
		p.manager.Set(menu)
	}
	return nil
}
func (p *PlayScreen) drawTitle(screen *ebiten.Image, title string, clr color.Color) {
	op := &text.DrawOptions{}
	op.GeoM.Translate(40, 30)
	op.ColorScale.ScaleWithColor(clr)
	text.Draw(screen, title, assets.MenuFont, op)
}
func (p *PlayScreen) drawPlayer(screen *ebiten.Image, id, name string, x, y float64, clr color.Color) {
	scale := mapAreaSize / float64(shared.DefaultGameConfig.MapSize)
	sx, sy := p.toScreen(x, y)
	radius := float32(shared.DefaultGameConfig.PlayerRadius) * float32(scale) * 1.8
	if radius < 10 {
		radius = 10
	}
	vector.FillCircle(screen, float32(sx), float32(sy), radius, clr, true)
	width, _ := text.Measure(name, assets.SmallFont, 0)
	op := &text.DrawOptions{}
	op.GeoM.Translate(sx-width/2, sy-float64(radius)-20)
	op.ColorScale.ScaleWithColor(color.White)
	text.Draw(screen, name, assets.SmallFont, op)
}
func (p *PlayScreen) drawFlag(screen *ebiten.Image, x, y float64, held bool) {
	scale := mapAreaSize / float64(shared.DefaultGameConfig.MapSize)
	sx, sy := p.toScreen(x, y)
	var radius float32
	if held {
		radius = float32(shared.DefaultGameConfig.PlayerRadius) * float32(scale) * 0.7
		if radius < 5 {
			radius = 5
		}
	} else {
		radius = float32(shared.DefaultGameConfig.PlayerRadius) * float32(scale) * 0.9
		if radius < 6 {
			radius = 6
		}
	}
	vector.FillCircle(screen, float32(sx), float32(sy), radius, color.RGBA{230, 170, 30, 255}, true)
}
func (p *PlayScreen) drawBoundary(screen *ebiten.Image) {
	center := float64(shared.DefaultGameConfig.MapSize) / 2
	scale := mapAreaSize / float64(shared.DefaultGameConfig.MapSize)
	sx, sy := p.toScreen(center, center)
	radius := float32(shared.DefaultGameConfig.CircleRadius) * float32(scale)
	vector.StrokeCircle(screen, float32(sx), float32(sy), radius, 2, color.RGBA{120, 120, 200, 255}, true)
}
func (p *PlayScreen) drawWinner(screen *ebiten.Image, winnerLabel string) {
	op := &text.DrawOptions{}
	op.GeoM.Translate(float64(ScreenWidth)/2-220, float64(ScreenHeight)/2-20)
	op.ColorScale.ScaleWithColor(color.RGBA{255, 220, 50, 255})
	text.Draw(screen, "Gana: "+winnerLabel, assets.TitleFont, op)
}
func (p *PlayScreen) Draw(screen *ebiten.Image) {
	p.drawBoundary(screen)
	if p.role == "server" {
		p.drawTitle(screen, "Modo Servidor - Partida en curso", color.RGBA{80, 220, 80, 255})
		if server.CurrentLobby == nil {
			return
		}
		flag := server.CurrentLobby.FlagState()
		for _, sp := range server.CurrentLobby.GetPlayers() {
			p.drawPlayer(screen, sp.Id, sp.Name, sp.X, sp.Y, color.RGBA{80, 180, 255, 255})
		}
		p.drawFlag(screen, flag.X, flag.Y, flag.Owner != nil)
		if server.CurrentLobby.IsFinished() {
			winnerID := server.CurrentLobby.Winner()
			label := winnerID
			if name, ok := p.names[winnerID]; ok {
				label = name
			}
			p.drawWinner(screen, label)
			p.backLobby.Draw(screen)
		}
		return
	}
	p.drawTitle(screen, "Modo Cliente - Jugando", color.White)
	if client.CurrentState == nil {
		return
	}
	myID := client.CurrentState.PlayerID()
	flag, players := client.CurrentState.GameState()
	for _, ps := range players {
		clr := color.RGBA{80, 180, 255, 255}
		if ps.ID == myID {
			clr = color.RGBA{90, 230, 130, 255}
		}
		name := p.names[ps.ID]
		if name == "" {
			name = ps.ID
		}
		p.drawPlayer(screen, ps.ID, name, ps.X, ps.Y, clr)
	}
	p.drawFlag(screen, flag.X, flag.Y, flag.Owner != nil)
	if winner := client.CurrentState.Winner(); winner != "" {
		label := winner
		if name, ok := p.names[winner]; ok {
			label = name
		}
		p.drawWinner(screen, label)
		p.backLobby.Draw(screen)
	}
	instructions := "WASD/flechas mueven, E toma o roba la bandera"
	op := &text.DrawOptions{}
	op.GeoM.Translate(40, float64(ScreenHeight)-30)
	op.ColorScale.ScaleWithColor(color.RGBA{210, 210, 210, 255})
	text.Draw(screen, instructions, assets.SmallFont, op)
}
