package components

import (
	"image"
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"

	"Proyecto1-cc8-23002455/ui/assets"
)

type PlayerCard struct {

	X float64
	Y float64

	W float64
	H float64

	Name string

	Hovered bool
}


func (p *PlayerCard) Draw(screen *ebiten.Image) {

	bg := color.RGBA{35, 55, 95, 255}

	if p.Hovered {
		bg = color.RGBA{55, 85, 140, 255}
	}

	vector.FillRect(
		screen,
		float32(p.X),
		float32(p.Y),
		float32(p.W),
		float32(p.H),
		bg,
		false,
	)

	vector.StrokeRect(
		screen,
		float32(p.X),
		float32(p.Y),
		float32(p.W),
		float32(p.H),
		2,
		color.White,
		false,
	)

	op := &text.DrawOptions{}

	op.GeoM.Translate(
		p.X+20,
		p.Y+30,
	)

	op.ColorScale.ScaleWithColor(
		color.White,
	)

	text.Draw(
		screen,
		p.Name,
		assets.MenuFont,
		op,
	)
}


func (p *PlayerCard) Update() {

	x, y := ebiten.CursorPosition()

	p.Hovered = image.Pt(x,y).In(
		image.Rect(
			int(p.X),
			int(p.Y),
			int(p.X+p.W),
			int(p.Y+p.H),
		),
	)

}