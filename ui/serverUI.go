package ui

import (
	"image/color"
	"strconv"

	"Proyecto1-cc8-23002455/server"
	"Proyecto1-cc8-23002455/ui/assets"
	"Proyecto1-cc8-23002455/ui/components"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
)

type ServerScreen struct {
	manager *Manager
	back    components.Button
	cards   []components.PlayerCard
	startBtn components.Button
	scroll  float64
}

func NewServer() *ServerScreen {

	s := &ServerScreen{}

	s.back = components.Button{
		X:    20,
		Y:    20,
		W:    220,
		H:    50,
		Text: "Regresar",
	}

	s.startBtn = components.Button{
		X:    ScreenWidth - 260,
		Y:    120,
		W:    220,
		H:    60,
		Text: "Iniciar",
	}

	return s
}

func (s *ServerScreen) Update() error {

	_, wheelY := ebiten.Wheel()

	s.scroll -= wheelY * 40

	if s.scroll < 0 {
		s.scroll = 0
	}

	maxScroll := float64(len(s.cards))*110 - 350
	if maxScroll < 0 {
		maxScroll = 0
	}

	if s.scroll > maxScroll {
		s.scroll = maxScroll
	}
	if server.CurrentLobby == nil {
		return nil
	}

	if server.CurrentLobby.IsPlaying() {
		play := NewPlay("server")
		play.manager = s.manager
		s.manager.Set(play)
		return nil
	}

	s.cards = nil
	players := server.CurrentLobby.GetPlayers()
	y := 220.0

	for _, p := range players {

		card := components.PlayerCard{
			X:      (ScreenWidth - 700) / 2,
			Y:      y,
			W:      700,
			H:      90,
			Name:   p.Name + " " + p.Id,
		}

		card.Update()
		s.cards = append(
			s.cards,
			card,
		)

		y += 110
	}

	s.startBtn.OnClick = func() {
		server.CurrentLobby.StartCountdown()
	}

	s.startBtn.Update()

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
		color.RGBA{80, 220, 80, 255},
	)

	text.Draw(
		screen,
		"Modo Servidor",
		assets.TitleFont,
		op,
	)

	s.startBtn.Draw(screen)

	for i := range s.cards {

		card := s.cards[i]

		card.Y -= s.scroll

		if card.Y+card.H < 220 {
			continue
		}

		if card.Y > ScreenHeight-80 {
			continue
		}

		card.Draw(screen)
	}

	s.back.Draw(screen)
	if server.CurrentLobby != nil {
		if seconds := server.CurrentLobby.CurrentCountdown(); seconds > 0 {
			op2 := &text.DrawOptions{}
			op2.GeoM.Translate(500, ScreenHeight-90)
			op2.ColorScale.ScaleWithColor(color.RGBA{255, 220, 50, 255})
			text.Draw(screen, "la partida inicia en...", assets.MenuFont, op2)

			op3 := &text.DrawOptions{}
			op3.GeoM.Translate(500, ScreenHeight-60)
			op3.ColorScale.ScaleWithColor(color.RGBA{255, 220, 50, 255})
			text.Draw(screen, strconv.Itoa(seconds), assets.TitleFont, op3)
		}
	}
}
