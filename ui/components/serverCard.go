package components

import (
	"image"
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"

	"Proyecto1-cc8-23002455/client"
	"Proyecto1-cc8-23002455/ui/assets"
)

type ServerCard struct {
	X float64
	Y float64

	W float64
	H float64

	Server client.DiscoveredServer

	Hovered bool

	OnClick func(server client.DiscoveredServer)
}

func (c *ServerCard) Draw(screen *ebiten.Image) {

	bg := color.RGBA{35, 55, 95, 255}

	if c.Hovered {
		bg = color.RGBA{55, 85, 140, 255}
	}

	vector.FillRect(
		screen,
		float32(c.X),
		float32(c.Y),
		float32(c.W),
		float32(c.H),
		bg,
		false,
	)

	vector.StrokeRect(
		screen,
		float32(c.X),
		float32(c.Y),
		float32(c.W),
		float32(c.H),
		2,
		color.White,
		false,
	)

	op := &text.DrawOptions{}
	op.GeoM.Translate(c.X+20, c.Y+28)
	op.ColorScale.ScaleWithColor(color.White)

	text.Draw(
		screen,
		c.Server.Name,
		assets.MenuFont,
		op,
	)

	op2 := &text.DrawOptions{}
	op2.GeoM.Translate(c.X+20, c.Y+62)
	op2.ColorScale.ScaleWithColor(color.RGBA{210, 210, 210, 255})

	text.Draw(
		screen,
		c.Server.IP,
		assets.SmallFont,
		op2,
	)
}

func (c *ServerCard) Update() {

	x, y := ebiten.CursorPosition()

	c.Hovered = image.Pt(x, y).In(image.Rect(
		int(c.X),
		int(c.Y),
		int(c.X+c.W),
		int(c.Y+c.H),
	))

	if c.Hovered &&
		inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {

		if c.OnClick != nil {
			c.OnClick(c.Server)
		}
	}
}
