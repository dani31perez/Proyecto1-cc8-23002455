package main

import (
	"log"
	"Proyecto1-cc8-23002455/ui"
	"Proyecto1-cc8-23002455/ui/assets"
	"github.com/hajimehoshi/ebiten/v2"
)

func main() {
	ebiten.SetWindowSize(ui.ScreenWidth, ui.ScreenHeight)
	ebiten.SetWindowTitle("Capture The Flag")

	assets.LoadFonts()
	game := ui.NewGame()

	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}
