package ui

import (
	"image/color"

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

	return s
}

func (s *ServerScreen) Update() error {


	if server.CurrentLobby == nil {
		return nil
	}

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

		s.cards = append(
			s.cards,
			card,
		)

		y += 110
	}

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

	for i := range s.cards {
		s.cards[i].Draw(screen)
	}
}
