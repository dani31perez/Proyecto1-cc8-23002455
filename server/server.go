package server

import (
	"Proyecto1-cc8-23002455/shared"
	"encoding/json"
	"errors"
	"fmt"
	"net"
	"strings"
	"unicode"
)

var CurrentLobby *lobby

func Run() {
	const tcpPort = 8889
	tcpListener, err := net.Listen("tcp", fmt.Sprintf(":%d", tcpPort))
	if err != nil {
		fmt.Println("error al abrir socket TCP:", err)
		return
	}
	defer tcpListener.Close()
	l := newLobby()
	CurrentLobby = l
	go runDiscovery(tcpPort, l)
	fmt.Println("servidor escuchando TCP en puerto", tcpPort)
	for {
		conn, err := tcpListener.Accept()
		if err != nil {
			fmt.Println("error al aceptar conexion:", err)
			continue
		}
		go handleClient(conn, l)
	}
}

// Go opens UDP sockets with address reuse support on the platforms where it is
// meaningful. ListenUDP is deliberately kept separate from the game TCP socket.
func runDiscovery(tcpPort int, l *lobby) {
	udpConn, err := listenDiscoveryUDP()
	if err != nil {
		fmt.Println("error al abrir socket UDP:", err)
		return
	}
	defer udpConn.Close()
	for {
		raw, remote, err := shared.ReadUDP(udpConn)
		if err != nil {
			continue // UDP discovery failures are intentionally silent.
		}
		var discover shared.DiscoverMessage
		if json.Unmarshal(raw, &discover) != nil || discover.Type != shared.TypeDiscover || discover.V != shared.ProtocolVersion {
			continue
		}
		state := "playing"
		if l.AcceptingPlayers() {
			state = "lobby"
		}
		_ = shared.SendUDP(udpConn, remote, shared.ServerInfoMessage{
			Type: shared.TypeServerInfo, V: shared.ProtocolVersion, Name: "servidor-ctf",
			TCPPort: tcpPort, State: state, Players: l.PlayerCount(),
		})
	}
}

func handleClient(netConn net.Conn, l *lobby) {
	conn := shared.NewConn(netConn)
	defer netConn.Close()
	var player *Player
	for player == nil {
		raw, err := conn.ReadMessage()
		if err != nil {
			if errors.Is(err, shared.ErrFrameTooLarge) {
				sendError(conn, shared.ErrMessageTooLarge)
			}
			return
		}
		msgType, code := messageType(raw)
		if code != "" {
			sendError(conn, code)
			continue
		}
		if msgType != shared.TypeJoin {
			sendError(conn, shared.ErrNotJoined)
			continue
		}
		join, code := decodeJoin(raw)
		if code != "" {
			sendError(conn, code)
			continue
		}
		if join.V != shared.ProtocolVersion {
			sendError(conn, shared.ErrVersionMismatch)
			return
		}
		name, ok := validName(join.Name)
		if !ok {
			sendError(conn, shared.ErrNameInvalid)
			continue
		}
		p, code := l.addPlayer(name, conn)
		if code != "" {
			sendError(conn, code)
			if code == shared.ErrLobbyFull || code == shared.ErrGameStarted {
				return
			}
			continue
		}
		player = p
		player.send(shared.WelcomeMessage{Type: shared.TypeWelcome, PlayerID: player.Id, Config: shared.DefaultGameConfig})
		l.BroadcastLobby()
	}
	defer l.removePlayer(player.Id)
	for {
		raw, err := conn.ReadMessage()
		if err != nil {
			if errors.Is(err, shared.ErrFrameTooLarge) {
				sendError(conn, shared.ErrMessageTooLarge)
			}
			return
		}
		msgType, code := messageType(raw)
		if code != "" {
			sendPlayerError(player, code)
			continue
		}
		switch msgType {
		case shared.TypeJoin:
			sendPlayerError(player, shared.ErrInvalidPhase)
		case shared.TypeInput:
			input, code := decodeInput(raw)
			if code != "" {
				sendPlayerError(player, code)
				continue
			}
			if !l.setPlayerDir(player.Id, input.Dir.X, input.Dir.Y) {
				sendPlayerError(player, shared.ErrInvalidPhase)
			}
		case shared.TypeInteract:
			if code := decodeInteract(raw); code != "" {
				sendPlayerError(player, code)
				continue
			}
			if !l.queueInteract(player.Id) {
				sendPlayerError(player, shared.ErrInvalidPhase)
			}
		default:
			sendPlayerError(player, shared.ErrUnknownType)
		}
	}
}

func sendError(conn *shared.Conn, reason string) {
	_ = conn.WriteMessage(shared.ErrorMessage{Type: shared.TypeError, Reason: reason})
}

func sendPlayerError(player *Player, reason string) {
	player.send(shared.ErrorMessage{Type: shared.TypeError, Reason: reason})
}

func objectFields(raw []byte) (map[string]json.RawMessage, string) {
	var fields map[string]json.RawMessage
	if err := json.Unmarshal(raw, &fields); err != nil || fields == nil {
		return nil, shared.ErrInvalidJSON
	}
	return fields, ""
}

func messageType(raw []byte) (string, string) {
	fields, code := objectFields(raw)
	if code != "" {
		return "", code
	}
	value, ok := fields["type"]
	if !ok {
		return "", shared.ErrMissingField
	}
	var messageType string
	if json.Unmarshal(value, &messageType) != nil || messageType == "" {
		return "", shared.ErrInvalidField
	}
	return messageType, ""
}

func decodeJoin(raw []byte) (shared.JoinMessage, string) {
	fields, code := objectFields(raw)
	if code != "" {
		return shared.JoinMessage{}, code
	}
	for _, name := range []string{"type", "v", "name"} {
		if _, ok := fields[name]; !ok {
			return shared.JoinMessage{}, shared.ErrMissingField
		}
	}
	var message shared.JoinMessage
	if json.Unmarshal(raw, &message) != nil || message.Type != shared.TypeJoin {
		return shared.JoinMessage{}, shared.ErrInvalidField
	}
	return message, ""
}

func decodeInput(raw []byte) (shared.InputMessage, string) {
	fields, code := objectFields(raw)
	if code != "" {
		return shared.InputMessage{}, code
	}
	for _, name := range []string{"type", "dir"} {
		if _, ok := fields[name]; !ok {
			return shared.InputMessage{}, shared.ErrMissingField
		}
	}
	var direction map[string]json.RawMessage
	if json.Unmarshal(fields["dir"], &direction) != nil || direction == nil {
		return shared.InputMessage{}, shared.ErrInvalidField
	}
	for _, name := range []string{"x", "y"} {
		if _, ok := direction[name]; !ok {
			return shared.InputMessage{}, shared.ErrMissingField
		}
	}
	var message shared.InputMessage
	if json.Unmarshal(raw, &message) != nil || message.Type != shared.TypeInput || message.Dir.X < -1 || message.Dir.X > 1 || message.Dir.Y < -1 || message.Dir.Y > 1 {
		return shared.InputMessage{}, shared.ErrInvalidField
	}
	return message, ""
}

func decodeInteract(raw []byte) string {
	fields, code := objectFields(raw)
	if code != "" {
		return code
	}
	if _, ok := fields["type"]; !ok {
		return shared.ErrMissingField
	}
	var message shared.InteractMessage
	if json.Unmarshal(raw, &message) != nil || message.Type != shared.TypeInteract {
		return shared.ErrInvalidField
	}
	return ""
}

func validName(name string) (string, bool) {
	name = strings.TrimSpace(name)
	if name == "" || len([]rune(name)) > shared.NameMaxLength {
		return "", false
	}
	for _, char := range name {
		if unicode.IsControl(char) || char == '\r' || char == '\n' {
			return "", false
		}
	}
	return name, true
}
