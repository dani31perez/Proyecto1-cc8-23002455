package client
import (
	"Proyecto1-cc8-23002455/shared"
	"sync"
)
type DiscoveredServer struct {
	shared.ServerInfoMessage
	IP string
}
type clientState struct {
	mu        sync.Mutex
	playerID  string
	config    shared.GameConfig
	lobby     []shared.LobbyPlayer
	countdown int
	started   bool
	flag      shared.FlagState
	players   []shared.PlayerState
	winner    string
}
var CurrentState *clientState
func newClientState() *clientState {
	return &clientState{}
}
func (s *clientState) setWelcome(playerID string, config shared.GameConfig) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.playerID = playerID
	s.config = config
}
func (s *clientState) setLobby(players []shared.LobbyPlayer) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.lobby = players
}
func (s *clientState) setCountdown(seconds int) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.countdown = seconds
}
func (s *clientState) setStarted() {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.started = true
}
func (s *clientState) setGameState(flag shared.FlagState, players []shared.PlayerState) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.flag = flag
	s.players = players
}
func (s *clientState) setGameOver(winner string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.winner = winner
}
func (s *clientState) PlayerID() string {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.playerID
}
func (s *clientState) Config() shared.GameConfig {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.config
}
func (s *clientState) LobbyPlayers() []shared.LobbyPlayer {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.lobby
}
func (s *clientState) Countdown() int {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.countdown
}
func (s *clientState) Started() bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.started
}
func (s *clientState) GameState() (shared.FlagState, []shared.PlayerState) {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.flag, s.players
}
func (s *clientState) Winner() string {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.winner
}
func (s *clientState) Reset() {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.playerID = ""
	s.config = shared.GameConfig{}
	s.lobby = nil
	s.countdown = 0
	s.started = false
	s.flag = shared.FlagState{}
	s.players = nil
	s.winner = ""
}