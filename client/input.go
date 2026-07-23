package main

import (
	"fmt"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"Proyecto1-cc8-23002455/shared"
	"time"
)

type keyListener struct {
	conn *shared.Conn

	lastX int
	lastY int

	lastSend time.Time
}

func (k *keyListener) Update() error {

	if inpututil.IsKeyJustPressed(ebiten.KeyQ) {
		return ebiten.Termination
	}

	dx, dy := 0, 0

	if ebiten.IsKeyPressed(ebiten.KeyW) {
		dy = -1
	}

	if ebiten.IsKeyPressed(ebiten.KeyS) {
		dy = 1
	}

	if ebiten.IsKeyPressed(ebiten.KeyA) {
		dx = -1
	}

	if ebiten.IsKeyPressed(ebiten.KeyD) {
		dx = 1
	}


	// Detecta cambio de dirección
	changed := dx != k.lastX || dy != k.lastY


	// Limita envío a 20 veces por segundo
	canSend := time.Since(k.lastSend) >= 50*time.Millisecond


	if (dx != 0 || dy != 0) && canSend {
		sendInput(k.conn, dx, dy)

		k.lastX = dx
		k.lastY = dy
		k.lastSend = time.Now()
	}


	// Si soltó todas las teclas manda detenerse
	if dx == 0 && dy == 0 && (k.lastX != 0 || k.lastY != 0) {

		sendInput(k.conn, 0, 0)

		k.lastX = 0
		k.lastY = 0
		k.lastSend = time.Now()
	}


	// Si cambió dirección manda inmediatamente
	if changed && (dx != 0 || dy != 0) {

		sendInput(k.conn, dx, dy)

		k.lastX = dx
		k.lastY = dy
		k.lastSend = time.Now()
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
	fmt.Printf("ANTES DE ENVIAR: %+v\n", msg)
	return conn.WriteMessage(msg)
}