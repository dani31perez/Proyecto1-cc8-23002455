package ui

import (
	"image"
	"image/color"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
)

type Button struct {
	X float64
	Y float64

	W float64
	H float64

	Hovered bool 
	Text string

	OnClick func()
}

func (b *Button) Draw(screen *ebiten.Image) {

	buttonColor := color.RGBA{35, 140, 255, 255}

	if b.Hovered {
		buttonColor = color.RGBA{60, 170, 255, 255}
	}
	vector.FillRect(
		screen,
		float32(b.X),
		float32(b.Y),
		float32(b.W),
		float32(b.H),
		buttonColor,
		false,
	)

	vector.FillRect(
		screen,
		float32(b.X),
		float32(b.Y),
		float32(b.W),
		2,
		color.White,
		false,
	)
	vector.FillRect(
		screen,
		float32(b.X),
		float32(b.Y+b.H-2),
		float32(b.W),
		2,
		color.White,
		false,
	)
	vector.FillRect(
		screen,
		float32(b.X),
		float32(b.Y),
		2,
		float32(b.H),
		color.White,
		false,
	)
	vector.FillRect(
		screen,
		float32(b.X+b.W-2),
		float32(b.Y),
		2,
		float32(b.H),
		color.White,
		false,
	)

	width, height := text.Measure(
		b.Text,
		MenuFont,
		0,
	)

	x := b.X + (b.W-width)/2
	y := b.Y + (b.H-height)/2

	op := &text.DrawOptions{}
	op.GeoM.Translate(x, y)
	op.ColorScale.ScaleWithColor(color.White)

	text.Draw(
		screen,
		b.Text,
		MenuFont,
		op,
	)
}

func (b *Button) Update() {
	x, y := ebiten.CursorPosition()

	inside := image.Pt(x, y).In(image.Rect(
		int(b.X),
		int(b.Y),
		int(b.X+b.W),
		int(b.Y+b.H),
	))

	b.Hovered = inside

	if inside && inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
		if b.OnClick != nil {
			b.OnClick()
		}
	}
}