package assets

import (
	"bytes"
	"log"
	"os"

	"github.com/hajimehoshi/ebiten/v2/text/v2"
)

var (
	TitleFont *text.GoTextFace
	MenuFont  *text.GoTextFace
	SmallFont *text.GoTextFace
)

func LoadFonts() {

	data, err := os.ReadFile("assets/PressStart2P-Regular.ttf")
	if err != nil {
		log.Fatal(err)
	}

	source, err := text.NewGoTextFaceSource(bytes.NewReader(data))
	if err != nil {
		log.Fatal(err)
	}

	TitleFont = &text.GoTextFace{
		Source: source,
		Size:   42,
	}

	MenuFont = &text.GoTextFace{
		Source: source,
		Size:   22,
	}

	SmallFont = &text.GoTextFace{
		Source: source,
		Size:   16,
	}
}