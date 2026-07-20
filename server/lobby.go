package main
import (
	"Proyecto1-cc8-23002455/shared"
	"fmt"
	"sync"
	"time"
)
type player struct {
	id   string
	name string
	conn *shared.Conn
}
type lobby struct {
	mu       sync.Mutex
	players  map[string]*player
	nextID   int
	counting bool
}
func newLobby() *lobby {
	return &lobby{players: make(map[string]*player)}
}
func (l *lobby) addPlayer(name string, conn *shared.Conn) *player {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.nextID++
	id := fmt.Sprintf("p%d", l.nextID)
	p := &player{id: id, name: name, conn: conn}
	l.players[id] = p
	return p
}
func (l *lobby) removePlayer(id string) {
	l.mu.Lock()
	defer l.mu.Unlock()
	delete(l.players, id)
}
func (l *lobby) snapshot() []shared.LobbyPlayer {
	l.mu.Lock()
	defer l.mu.Unlock()
	list := make([]shared.LobbyPlayer, 0, len(l.players))
	for _, p := range l.players {
		list = append(list, shared.LobbyPlayer{ID: p.id, Name: p.name})
	}
	return list
}
func (l *lobby) broadcast(v interface{}) {
	l.mu.Lock()
	defer l.mu.Unlock()
	for _, p := range l.players {
		p.conn.WriteMessage(v)
	}
}
func (l *lobby) broadcastLobby() {
	msg := shared.LobbyMessage{Type: shared.TypeLobby, Players: l.snapshot()}
	l.broadcast(msg)
}
func (l *lobby) startCountdownOnce() {
	l.mu.Lock()
	if l.counting {
		l.mu.Unlock()
		return
	}
	l.counting = true
	l.mu.Unlock()
	go l.runCountdown()
}
func (l *lobby) runCountdown() {
	for seconds := 5; seconds >= 1; seconds-- {
		msg := shared.CountdownMessage{Type: shared.TypeCountdown, Seconds: seconds}
		l.broadcast(msg)
		fmt.Println("countdown enviado:", seconds)
		time.Sleep(1 * time.Second)
	}
	l.broadcast(shared.StartMessage{Type: shared.TypeStart})
	fmt.Println("start enviado, fin de la simulacion de espera en el lobby")
	l.mu.Lock()
	l.counting = false
	l.mu.Unlock()
}
