package client

import (
	"image/color"
	"time"

	"Proyecto1-cc8-23002455/ui/assets"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

type Notification struct {
	Title   string
	Message string

	Visible bool
	until   time.Time
}

var Notify Notification

func (n *Notification) Show(title, message string, seconds int) {
	n.Title = title
	n.Message = message
	n.Visible = true
	n.until = time.Now().Add(time.Duration(seconds) * time.Second)
}

func (n *Notification) Update() {
	if n.Visible && time.Now().After(n.until) {
		n.Visible = false
	}
}

func (n *Notification) Draw(screen *ebiten.Image) {
	if !n.Visible {
		return
	}

	w := float32(650)
	h := float32(150)

	x := float32((1480 - int(w)) / 2)
	y := float32(50)

	vector.DrawFilledRect(
		screen,
		x,
		y,
		w,
		h,
		color.RGBA{40, 40, 40, 240},
		false,
	)

	vector.StrokeRect(
		screen,
		x,
		y,
		w,
		h,
		3,
		color.RGBA{220, 70, 70, 255},
		false,
	)

	title := &text.DrawOptions{}
	title.GeoM.Translate(float64(x+20), float64(y+35))
	title.ColorScale.ScaleWithColor(color.White)

	text.Draw(
		screen,
		n.Title,
		assets.MenuFont,
		title,
	)

	msg := &text.DrawOptions{}
	msg.GeoM.Translate(float64(x+20), float64(y+90))
	msg.ColorScale.ScaleWithColor(color.RGBA{255, 210, 210, 255})

	text.Draw(
		screen,
		n.Message,
		assets.SmallFont,
		msg,
	)
}