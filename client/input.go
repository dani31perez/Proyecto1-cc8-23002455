package main

import (
	"fmt"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"Proyecto1-cc8-23002455/shared"
)

type keyListener struct {
	conn *shared.Conn
}

func (k *keyListener) Update() error {
	// Salir con Q
	if inpututil.IsKeyJustPressed(ebiten.KeyQ) {
		return ebiten.Termination
	}

	dx, dy := 0, 0
	pressed := false

	// Evaluamos qué teclas se acaban de presionar
	if inpututil.IsKeyJustPressed(ebiten.KeyW) {
		dy = -1
		pressed = true
	} else if inpututil.IsKeyJustPressed(ebiten.KeyS) {
		dy = 1
		pressed = true
	}

	if inpututil.IsKeyJustPressed(ebiten.KeyA) {
		dx = -1
		pressed = true
	} else if inpututil.IsKeyJustPressed(ebiten.KeyD) {
		dx = 1
		pressed = true
	}

	if pressed {
		sendInput(k.conn, dx, dy)
	}

	return nil
}

func (k *keyListener) Draw(screen *ebiten.Image) {}

func (k *keyListener) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return 200, 200
}

func startInputLoop(conn *shared.Conn) {
	ebiten.SetWindowSize(200, 200)
	ebiten.SetWindowTitle("Input WASD")
	
	game := &keyListener{conn: conn}

	if err := ebiten.RunGame(game); err != nil && err != ebiten.Termination {
		fmt.Println("error en el listener de teclas:", err)
	}
}

func sendInput(conn *shared.Conn, dx, dy int) error {
	msg := shared.InputMessage{
		Type: shared.TypeInput,
		Dir: shared.Direction{
			X: dx,
			Y: dy,
		},
	}
	return conn.WriteMessage(msg)
}