package client

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

	dx, dy := 0, 0

	if ebiten.IsKeyPressed(ebiten.KeyW) || ebiten.IsKeyPressed(ebiten.KeyUp) {
		dy = -1
	}
	if ebiten.IsKeyPressed(ebiten.KeyS) || ebiten.IsKeyPressed(ebiten.KeyDown) {
		dy = 1
	}
	if ebiten.IsKeyPressed(ebiten.KeyA) || ebiten.IsKeyPressed(ebiten.KeyLeft) {
		dx = -1
	}
	if ebiten.IsKeyPressed(ebiten.KeyD) || ebiten.IsKeyPressed(ebiten.KeyRight) {
		dx = 1
	}

	if inpututil.IsKeyJustPressed(ebiten.KeyE) {
		fmt.Println("[ACCION] Enviando 'interact' al servidor")
		k.conn.WriteMessage(shared.InteractMessage{Type: shared.TypeInteract})
	}

	if dx != k.lastX || dy != k.lastY {
		fmt.Printf("[INTENCIÓN] Cambio detectado -> Enviando: X=%d, Y=%d\n", dx, dy)
		sendInput(k.conn, dx, dy)

		k.conn.WriteMessage(shared.InputMessage{Type: shared.TypeInput,Dir: shared.Direction{X: dx,Y: dy,}})
		
		k.lastX = dx
		k.lastY = dy
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