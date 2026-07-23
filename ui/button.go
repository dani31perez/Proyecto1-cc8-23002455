package ui

import (
	"image"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/text"
	"golang.org/x/image/font/basicfont"
)

type Button struct {
	X float64
	Y float64

	W float64
	H float64

	Text string

	OnClick func()
}

func (b *Button) Draw(screen *ebiten.Image) {

	ebitenutil.DrawRect(screen,
		b.X,
		b.Y,
		b.W,
		b.H,
		nil,
	)

	text.Draw(
		screen,
		b.Text,
		basicfont.Face7x13,
		int(b.X+20),
		int(b.Y+b.H/2),
		nil,
	)
}

func (b *Button) Update() {

	if !ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft) {
		return
	}

	x, y := ebiten.CursorPosition()

	if image.Pt(x, y).In(image.Rect(
		int(b.X),
		int(b.Y),
		int(b.X+b.W),
		int(b.Y+b.H),
	)) {

		if b.OnClick != nil {
			b.OnClick()
		}
	}
}