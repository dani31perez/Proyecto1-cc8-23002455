package components

import (
	"image"
	"image/color"
	"strings"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/text/v2"

	"Proyecto1-cc8-23002455/ui/assets"
)

type TextBox struct {
	X, Y float64
	W, H float64

	Text   string
	Active bool
}

func (t *TextBox) Update() {

	// Click para activar/desactivar
	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
		x, y := ebiten.CursorPosition()

		r := image.Rect(
		int(t.X),
		int(t.Y),
		int(t.X+t.W),
		int(t.Y+t.H),
)

t.Active = image.Pt(x, y).In(r)
	}

	if !t.Active {
		return
	}

	// Escribir caracteres
	for _, r := range ebiten.AppendInputChars(nil) {
		if r >= 32 && r <= 126 {
			t.Text += string(r)
		}
	}

	// Borrar
	if inpututil.IsKeyJustPressed(ebiten.KeyBackspace) && len(t.Text) > 0 {
		t.Text = t.Text[:len(t.Text)-1]
	}

	// Evitar textos muy largos
	if len(t.Text) > 30 {
		t.Text = t.Text[:30]
	}
}

func (t *TextBox) Draw(screen *ebiten.Image) {

	ebitenutil.DrawRect(
		screen,
		t.X,
		t.Y,
		t.W,
		t.H,
		color.RGBA{40, 40, 40, 255},
	)

	ebitenutil.DrawRect(screen, t.X, t.Y, t.W, 2, color.White)
	ebitenutil.DrawRect(screen, t.X, t.Y+t.H-2, t.W, 2, color.White)
	ebitenutil.DrawRect(screen, t.X, t.Y, 2, t.H, color.White)
	ebitenutil.DrawRect(screen, t.X+t.W-2, t.Y, 2, t.H, color.White)

	txt := t.Text

	if t.Active {
		if len(txt)%2 == 0 {
			txt += "|"
		}
	}

	op := &text.DrawOptions{}
	op.GeoM.Translate(t.X+10, t.Y+13)

	text.Draw(
		screen,
		strings.TrimSpace(txt),
		assets.SmallFont,
		op,
	)
}