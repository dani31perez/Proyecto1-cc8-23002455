package server
import (
	"Proyecto1-cc8-23002455/shared"
	"fmt"
	"math"
	"math/rand"
	"sync"
	"time"
)
type Player struct {
	Id   string
	Name string
	conn *shared.Conn
	X    float64
	Y    float64
	DirX int
	DirY int
}
type lobby struct {
	mu        sync.Mutex
	players   map[string]*Player
	nextID    int
	counting  bool
	countdown int
	playing   bool
	finished  bool
	winner    string
	flagOwner *string
	flagX     float64
	flagY     float64
}
func newLobby() *lobby {
	return &lobby{players: make(map[string]*Player)}
}
func (l *lobby) addPlayer(name string, conn *shared.Conn) *Player {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.nextID++
	id := fmt.Sprintf("p%d", l.nextID)
	p := &Player{Id: id, Name: name, conn: conn}
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
		list = append(list, shared.LobbyPlayer{ID: p.Id, Name: p.Name})
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
func (l *lobby) broadcastLocked(v interface{}) {
	for _, p := range l.players {
		p.conn.WriteMessage(v)
	}
}
func (l *lobby) broadcastLobby() {
	msg := shared.LobbyMessage{Type: shared.TypeLobby, Players: l.snapshot()}
	l.broadcast(msg)
}
func (l *lobby) IsPlaying() bool {
	l.mu.Lock()
	defer l.mu.Unlock()
	return l.playing
}
func (l *lobby) IsFinished() bool {
	l.mu.Lock()
	defer l.mu.Unlock()
	return l.finished
}
func (l *lobby) Winner() string {
	l.mu.Lock()
	defer l.mu.Unlock()
	return l.winner
}
func (l *lobby) CurrentCountdown() int {
	l.mu.Lock()
	defer l.mu.Unlock()
	return l.countdown
}
func (l *lobby) FlagState() shared.FlagState {
	l.mu.Lock()
	defer l.mu.Unlock()
	return shared.FlagState{Owner: l.flagOwner, X: l.flagX, Y: l.flagY}
}
func (l *lobby) startCountdownOnce() {
	l.mu.Lock()
	if l.counting || l.playing || l.finished {
		l.mu.Unlock()
		return
	}
	l.counting = true
	l.mu.Unlock()
	go l.runCountdown()
}
func (l *lobby) runCountdown() {
	for seconds := 5; seconds >= 1; seconds-- {
		l.mu.Lock()
		l.countdown = seconds
		l.mu.Unlock()
		msg := shared.CountdownMessage{Type: shared.TypeCountdown, Seconds: seconds}
		l.broadcast(msg)
		fmt.Println("countdown enviado:", seconds)
		time.Sleep(1 * time.Second)
	}
	l.broadcast(shared.StartMessage{Type: shared.TypeStart})
	fmt.Println("start enviado, comienza la partida")
	l.mu.Lock()
	l.counting = false
	l.countdown = 0
	l.mu.Unlock()
	l.beginGame()
}
func (l *lobby) beginGame() {
	config := shared.DefaultGameConfig
	l.mu.Lock()
	l.playing = true
	center := float64(config.MapSize) / 2
	l.flagOwner = nil
	l.flagX = center
	l.flagY = center
	for _, p := range l.players {
		angle := rand.Float64() * 2 * math.Pi
		radius := float64(config.CircleRadius+config.PlayerRadius) + 5 + rand.Float64()*150
		p.X = center + radius*math.Cos(angle)
		p.Y = center + radius*math.Sin(angle)
		p.DirX = 0
		p.DirY = 0
	}
	l.mu.Unlock()
	go l.runGameLoop()
}
func clampDir(v int) int {
	if v > 1 {
		return 1
	}
	if v < -1 {
		return -1
	}
	return v
}
func (l *lobby) setPlayerDir(id string, dx, dy int) {
	l.mu.Lock()
	defer l.mu.Unlock()
	p, ok := l.players[id]
	if !ok {
		return
	}
	p.DirX = clampDir(dx)
	p.DirY = clampDir(dy)
}
func distance(x1, y1, x2, y2 float64) float64 {
	dx := x2 - x1
	dy := y2 - y1
	return math.Sqrt(dx*dx + dy*dy)
}
func (l *lobby) interact(id string) {
	l.mu.Lock()
	defer l.mu.Unlock()
	if !l.playing {
		return
	}
	p, ok := l.players[id]
	if !ok {
		return
	}
	radius := float64(shared.DefaultGameConfig.InteractRadius)
	if l.flagOwner == nil {
		if distance(p.X, p.Y, l.flagX, l.flagY) <= radius {
			owner := id
			l.flagOwner = &owner
			fmt.Println("bandera capturada por", id)
		}
		return
	}
	if *l.flagOwner == id {
		return
	}
	holder, ok := l.players[*l.flagOwner]
	if !ok {
		return
	}
	if distance(p.X, p.Y, holder.X, holder.Y) <= radius {
		owner := id
		l.flagOwner = &owner
		fmt.Println("bandera robada por", id, "a", holder.Id)
	}
}
func clampFloat(v, min, max float64) float64 {
	if v < min {
		return min
	}
	if v > max {
		return max
	}
	return v
}
func (l *lobby) buildStateLocked() shared.StateMessage {
	players := make([]shared.PlayerState, 0, len(l.players))
	for _, p := range l.players {
		players = append(players, shared.PlayerState{ID: p.Id, X: p.X, Y: p.Y})
	}
	return shared.StateMessage{
		Type:    shared.TypeState,
		Flag:    shared.FlagState{Owner: l.flagOwner, X: l.flagX, Y: l.flagY},
		Players: players,
	}
}
func (l *lobby) runGameLoop() {
	config := shared.DefaultGameConfig
	dt := 1.0 / float64(config.TickRate)
	ticker := time.NewTicker(time.Second / time.Duration(config.TickRate))
	defer ticker.Stop()
	minBound := float64(config.PlayerRadius)
	maxBound := float64(config.MapSize) - float64(config.PlayerRadius)
	center := float64(config.MapSize) / 2
	winDistance := float64(config.CircleRadius + config.PlayerRadius)
	for range ticker.C {
		l.mu.Lock()
		if l.finished {
			l.mu.Unlock()
			return
		}
		for _, p := range l.players {
			if p.DirX == 0 && p.DirY == 0 {
				continue
			}
			m := math.Sqrt(float64(p.DirX*p.DirX + p.DirY*p.DirY))
			nx := float64(p.DirX) / m
			ny := float64(p.DirY) / m
			p.X = clampFloat(p.X+nx*float64(config.Speed)*dt, minBound, maxBound)
			p.Y = clampFloat(p.Y+ny*float64(config.Speed)*dt, minBound, maxBound)
		}
		if l.flagOwner != nil {
			holder, ok := l.players[*l.flagOwner]
			if ok {
				l.flagX = holder.X
				l.flagY = holder.Y
				if distance(l.flagX, l.flagY, center, center) > winDistance {
					l.finished = true
					l.playing = false
					l.winner = holder.Id
					winnerName := holder.Name
					l.broadcastLocked(shared.GameOverMessage{Type: shared.TypeGameOver, Winner: holder.Id})
					fmt.Println("partida finalizada, gana:", holder.Id, winnerName)
					l.mu.Unlock()
					return
				}
			}
		}
		l.broadcastLocked(l.buildStateLocked())
		l.mu.Unlock()
	}
}
func (l *lobby) GetPlayers() []*Player {
	l.mu.Lock()
	defer l.mu.Unlock()
	var players []*Player
	for _, p := range l.players {
		players = append(players, p)
	}
	return players
}
