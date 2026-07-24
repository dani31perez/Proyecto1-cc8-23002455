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
}
type Client struct {
    State *clientState

    conn *shared.Conn

    Servers []DiscoveredServer
}
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