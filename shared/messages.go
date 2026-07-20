package shared
const (
	TypeDiscover   = "discover"
	TypeServerInfo = "server_info"
	TypeJoin       = "join"
	TypeInput      = "input"
	TypeInteract   = "interact"
	TypeWelcome    = "welcome"
	TypeLobby      = "lobby"
	TypeCountdown  = "countdown"
	TypeStart      = "start"
	TypeState      = "state"
	TypeGameOver   = "game_over"
	TypeError      = "error"
)
type TypeOnly struct {
	Type string `json:"type"`
}
type DiscoverMessage struct {
	Type string `json:"type"`
	V    int    `json:"v"`
}
type ServerInfoMessage struct {
	Type    string `json:"type"`
	V       int    `json:"v"`
	Name    string `json:"name"`
	TCPPort int    `json:"tcp_port"`
	State   string `json:"state"`
	Players int    `json:"players"`
}
type JoinMessage struct {
	Type string `json:"type"`
	V    int    `json:"v"`
	Name string `json:"name"`
}
type Direction struct {
	X int `json:"x"`
	Y int `json:"y"`
}
type InputMessage struct {
	Type string    `json:"type"`
	Dir  Direction `json:"dir"`
}
type InteractMessage struct {
	Type string `json:"type"`
}
type GameConfig struct {
	MapSize        int `json:"map_size"`
	CircleRadius   int `json:"circle_radius"`
	PlayerRadius   int `json:"player_radius"`
	InteractRadius int `json:"interact_radius"`
	Speed          int `json:"speed"`
	TickRate       int `json:"tick_rate"`
}
type WelcomeMessage struct {
	Type     string     `json:"type"`
	PlayerID string     `json:"player_id"`
	Config   GameConfig `json:"config"`
}
type LobbyPlayer struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}
type LobbyMessage struct {
	Type    string        `json:"type"`
	Players []LobbyPlayer `json:"players"`
}
type CountdownMessage struct {
	Type    string `json:"type"`
	Seconds int    `json:"seconds"`
}
type StartMessage struct {
	Type string `json:"type"`
}
type FlagState struct {
	Owner *string `json:"owner"`
	X     float64 `json:"x"`
	Y     float64 `json:"y"`
}
type PlayerState struct {
	ID string  `json:"id"`
	X  float64 `json:"x"`
	Y  float64 `json:"y"`
}
type StateMessage struct {
	Type    string        `json:"type"`
	Flag    FlagState     `json:"flag"`
	Players []PlayerState `json:"players"`
}
type GameOverMessage struct {
	Type   string `json:"type"`
	Winner string `json:"winner"`
}
type ErrorMessage struct {
	Type   string `json:"type"`
	Reason string `json:"reason"`
}
