package server

import (
	"Proyecto1-cc8-23002455/shared"
	"fmt"
	"math"
	"math/rand"
	"sort"
	"sync"
	"time"
    "strconv"
    "strings"
)

const (
	phaseLobby = iota
	phaseCountdown
	phasePlaying
	phasePostGame
)

type Player struct {
	Id         string
	Name       string
	conn       *shared.Conn
	sender     *outboundSender
	X, Y       float64
	DirX, DirY int
}

// A delayed client may skip obsolete state snapshots, but protocol control
// messages remain queued and ordered. The worker is the sole socket writer.
type outboundSender struct {
	conn   *shared.Conn
	mu     sync.Mutex
	cond   *sync.Cond
	queue  []outboundMessage
	closed bool
}
type outboundMessage struct {
	value interface{}
	state bool
}

func newOutboundSender(conn *shared.Conn) *outboundSender {
	s := &outboundSender{conn: conn}
	s.cond = sync.NewCond(&s.mu)
	go s.run()
	return s
}
func (s *outboundSender) send(value interface{}) {
	_, isState := value.(shared.StateMessage)
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.closed {
		return
	}
	if isState {
		for i := len(s.queue) - 1; i >= 0; i-- {
			if !s.queue[i].state {
				break
			}
			s.queue[i].value = value
			s.cond.Signal()
			return
		}
	}
	s.queue = append(s.queue, outboundMessage{value: value, state: isState})
	s.cond.Signal()
}
func (s *outboundSender) run() {
	for {
		s.mu.Lock()
		for len(s.queue) == 0 && !s.closed {
			s.cond.Wait()
		}
		if s.closed {
			s.mu.Unlock()
			return
		}
		message := s.queue[0]
		s.queue = s.queue[1:]
		s.mu.Unlock()
		if s.conn.WriteMessage(message.value) != nil {
			return
		}
	}
}
func (s *outboundSender) close() {
	s.mu.Lock()
	s.closed = true
	s.queue = nil
	s.cond.Broadcast()
	s.mu.Unlock()
}
func (p *Player) send(value interface{}) { p.sender.send(value) }

type lobby struct {
	mu                  sync.Mutex
	players             map[string]*Player
	nextID              int
	phase               int
	countdown           int
	cycle               uint64
	winner              string
	flagOwner           *string
	flagX, flagY        float64
	carrierWasInside    bool
	pendingInteractions []string
}

func newLobby() *lobby { return &lobby{players: make(map[string]*Player)} }

func (l *lobby) addPlayer(name string, conn *shared.Conn) (*Player, string) {
	l.mu.Lock()
	defer l.mu.Unlock()
	if l.phase != phaseLobby {
		return nil, shared.ErrGameStarted
	}
	if len(l.players) >= shared.MaxPlayers {
		return nil, shared.ErrLobbyFull
	}
	l.nextID++
	p := &Player{Id: fmt.Sprintf("p%d", l.nextID), Name: name, conn: conn, sender: newOutboundSender(conn)}
	l.players[p.Id] = p
	return p, ""
}

func (l *lobby) removePlayer(id string) {
	shouldLobby := false
	l.mu.Lock()
	p, exists := l.players[id]
	if !exists {
		l.mu.Unlock()
		return
	}
	if l.flagOwner != nil && *l.flagOwner == id {
		l.resetFlagLocked()
	}
	delete(l.players, id)
	p.sender.close()
	if l.phase == phaseCountdown && len(l.players) < shared.MinPlayers {
		l.phase, l.countdown = phaseLobby, 0
		l.cycle++
		shouldLobby = true
	} else if l.phase == phasePlaying && len(l.players) == 0 {
		l.resetRoundLocked()
		shouldLobby = true
	} else if l.phase == phaseLobby {
		shouldLobby = true
	}
	l.mu.Unlock()
	if shouldLobby {
		l.BroadcastLobby()
	}
}

func (l *lobby) snapshotLocked() []shared.LobbyPlayer {
	list := make([]shared.LobbyPlayer, 0, len(l.players))
	for _, p := range l.players {
		list = append(list, shared.LobbyPlayer{ID: p.Id, Name: p.Name})
	}
	sort.Slice(list, func(i, j int) bool { return list[i].ID < list[j].ID })
	return list
}

func (l *lobby) broadcast(v interface{}) {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.broadcastLocked(v)
}
func (l *lobby) broadcastLocked(v interface{}) {
	for _, p := range l.players {
		p.send(v)
	}
}
func (l *lobby) BroadcastLobby() {
	l.mu.Lock()
	msg := shared.LobbyMessage{Type: shared.TypeLobby, Players: l.snapshotLocked()}
	l.broadcastLocked(msg)
	l.mu.Unlock()
}
func (l *lobby) AcceptingPlayers() bool {
	l.mu.Lock()
	defer l.mu.Unlock()
	return l.phase == phaseLobby
}
func (l *lobby) PlayerCount() int      { l.mu.Lock(); defer l.mu.Unlock(); return len(l.players) }
func (l *lobby) IsPlaying() bool       { l.mu.Lock(); defer l.mu.Unlock(); return l.phase == phasePlaying }
func (l *lobby) IsFinished() bool      { l.mu.Lock(); defer l.mu.Unlock(); return l.phase == phasePostGame }
func (l *lobby) Winner() string        { l.mu.Lock(); defer l.mu.Unlock(); return l.winner }
func (l *lobby) CurrentCountdown() int { l.mu.Lock(); defer l.mu.Unlock(); return l.countdown }
func (l *lobby) FlagState() shared.FlagState {
	l.mu.Lock()
	defer l.mu.Unlock()
	return shared.FlagState{Owner: l.flagOwner, X: l.flagX, Y: l.flagY}
}

func (l *lobby) StartCountdown() {
	l.mu.Lock()
	if l.phase != phaseLobby || len(l.players) < shared.MinPlayers {
		l.mu.Unlock()
		return
	}
	l.phase, l.countdown = phaseCountdown, shared.CountdownSeconds
	l.cycle++
	cycle := l.cycle
	l.mu.Unlock()
	go l.runCountdown(cycle)
}
func (l *lobby) runCountdown(cycle uint64) {
	for seconds := shared.CountdownSeconds; seconds >= 1; seconds-- {
		l.mu.Lock()
		if l.phase != phaseCountdown || l.cycle != cycle || len(l.players) < shared.MinPlayers {
			l.mu.Unlock()
			return
		}
		l.countdown = seconds
		l.broadcastLocked(shared.CountdownMessage{Type: shared.TypeCountdown, Seconds: seconds})
		l.mu.Unlock()
		time.Sleep(time.Second)
	}
	l.mu.Lock()
	if l.phase != phaseCountdown || l.cycle != cycle || len(l.players) < shared.MinPlayers {
		l.mu.Unlock()
		return
	}
	l.countdown, l.phase = 0, phasePlaying
	l.spawnRoundLocked()
	l.broadcastLocked(shared.StartMessage{Type: shared.TypeStart})
	l.mu.Unlock()
	go l.runGameLoop(cycle)
}

func (l *lobby) spawnRoundLocked() {
	center := float64(shared.DefaultGameConfig.MapSize) / 2
	l.resetFlagLocked()
	for _, p := range l.players {
		angle := rand.Float64() * 2 * math.Pi
		radius := 350 + rand.Float64()*100
		p.X, p.Y = center+radius*math.Cos(angle), center+radius*math.Sin(angle)
		p.DirX, p.DirY = 0, 0
	}
	l.pendingInteractions = nil
}
func (l *lobby) resetFlagLocked() {
	l.flagOwner = nil
	l.flagX, l.flagY = 500, 500
	l.carrierWasInside = false
}
func (l *lobby) resetRoundLocked() {
	l.phase, l.countdown, l.winner = phaseLobby, 0, ""
	l.resetFlagLocked()
	l.pendingInteractions = nil
	for _, p := range l.players {
		p.DirX, p.DirY = 0, 0
	}
}

func (l *lobby) setPlayerDir(id string, dx, dy int) bool {
	l.mu.Lock()
	defer l.mu.Unlock()
	if l.phase != phasePlaying {
		return false
	}
	p, ok := l.players[id]
	if !ok {
		return false
	}
	p.DirX, p.DirY = dx, dy
	return true
}
func (l *lobby) queueInteract(id string) bool {
	l.mu.Lock()
	defer l.mu.Unlock()
	if l.phase != phasePlaying {
		return false
	}
	if _, ok := l.players[id]; !ok {
		return false
	}
	l.pendingInteractions = append(l.pendingInteractions, id)
	return true
}

func distance(x1, y1, x2, y2 float64) float64 { return math.Hypot(x2-x1, y2-y1) }
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
		players = append(players, shared.PlayerState{ID: p.Id, X: roundTenth(p.X), Y: roundTenth(p.Y)})
	}
	sort.Slice(players, func(i, j int) bool { return players[i].ID < players[j].ID })
	return shared.StateMessage{Type: shared.TypeState, Flag: shared.FlagState{Owner: l.flagOwner, X: roundTenth(l.flagX), Y: roundTenth(l.flagY)}, Players: players}
}
func roundTenth(v float64) float64 { return math.Floor(math.Abs(v)*10+0.5) / 10 * math.Copysign(1, v) }

func (l *lobby) runGameLoop(cycle uint64) {
	config := shared.DefaultGameConfig
	dt := 1.0 / float64(config.TickRate)
	ticker := time.NewTicker(time.Second / time.Duration(config.TickRate))
	defer ticker.Stop()
	for range ticker.C {
		l.mu.Lock()
		if l.phase != phasePlaying || l.cycle != cycle {
			l.mu.Unlock()
			return
		}
		for _, p := range l.players {
			if p.DirX == 0 && p.DirY == 0 {
				continue
			}
			m := math.Hypot(float64(p.DirX), float64(p.DirY))
			p.X = clampFloat(p.X+float64(p.DirX)/m*float64(config.Speed)*dt, 15, 985)
			p.Y = clampFloat(p.Y+float64(p.DirY)/m*float64(config.Speed)*dt, 15, 985)
		}
		if l.evaluateVictoryLocked() {
			l.mu.Unlock()
			go l.returnToLobbyAfterGame(cycle)
			return
		}
		l.processInteractionsLocked()
		l.broadcastLocked(l.buildStateLocked())
		l.mu.Unlock()
	}
}
func (l *lobby) evaluateVictoryLocked() bool {
	if l.flagOwner == nil {
		return false
	}
	holder, ok := l.players[*l.flagOwner]
	if !ok {
		l.resetFlagLocked()
		return false
	}
	l.flagX, l.flagY = holder.X, holder.Y
	inside := distance(holder.X, holder.Y, 500, 500) <= 315
	if l.carrierWasInside && !inside {
		l.phase, l.winner = phasePostGame, holder.Id
		l.broadcastLocked(shared.GameOverMessage{Type: shared.TypeGameOver, Winner: holder.Id})
		return true
	}
	if inside {
		l.carrierWasInside = true
	}
	return false
}
func (l *lobby) processInteractionsLocked() {
	pending := l.pendingInteractions
	l.pendingInteractions = nil
	for _, id := range pending {
		p, ok := l.players[id]
		if !ok {
			continue
		}
		if l.flagOwner == nil {
			if distance(p.X, p.Y, l.flagX, l.flagY) <= float64(shared.DefaultGameConfig.InteractRadius) {
				owner := id
				l.flagOwner = &owner
				l.flagX, l.flagY = p.X, p.Y
				l.carrierWasInside = distance(p.X, p.Y, 500, 500) <= 315
			}
			continue
		}
		if *l.flagOwner == id {
			continue
		}
		holder, ok := l.players[*l.flagOwner]
		if ok && distance(p.X, p.Y, holder.X, holder.Y) <= float64(shared.DefaultGameConfig.InteractRadius) {
			owner := id
			l.flagOwner = &owner
			l.flagX, l.flagY = p.X, p.Y
			l.carrierWasInside = distance(p.X, p.Y, 500, 500) <= 315
		}
	}
}
func (l *lobby) returnToLobbyAfterGame(cycle uint64) {
	time.Sleep(shared.PostGameSeconds * time.Second)
	l.mu.Lock()
	if l.phase != phasePostGame || l.cycle != cycle {
		l.mu.Unlock()
		return
	}
	l.resetRoundLocked()
	msg := shared.LobbyMessage{Type: shared.TypeLobby, Players: l.snapshotLocked()}
	l.broadcastLocked(msg)
	l.mu.Unlock()
}
func (l *lobby) GetPlayers() []*Player {
	l.mu.Lock()
	defer l.mu.Unlock()

	players := make([]*Player, 0, len(l.players))

	for _, p := range l.players {
		players = append(players, p)
	}

	sort.Slice(players, func(i, j int) bool {
		id1, _ := strconv.Atoi(strings.TrimPrefix(players[i].Id, "p"))
		id2, _ := strconv.Atoi(strings.TrimPrefix(players[j].Id, "p"))
		return id1 < id2
	})

	return players
}
func (l *lobby) ResetLobby() { l.mu.Lock(); l.cycle++; l.resetRoundLocked(); l.mu.Unlock() }
