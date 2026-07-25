package client
import (
	"Proyecto1-cc8-23002455/shared"
	"fmt"
	"net"
	"time"
)

var CurrentConn *shared.Conn
func Run(server DiscoveredServer, name string) {
	tcpConn, err := net.Dial("tcp", fmt.Sprintf("%s:%d", server.IP, server.TCPPort))
	if err != nil {
		fmt.Println("error al conectar por TCP:", err)
		return
	}
	conn := shared.NewConn(tcpConn)
	CurrentConn = conn
	state := newClientState()
	CurrentState = state
	done := make(chan struct{})
	go readLoop(conn, state, done)
	join := shared.JoinMessage{Type: shared.TypeJoin, V: 1, Name: name}
	conn.WriteMessage(join)
	<-done
}

func readLoop(conn *shared.Conn, state *clientState, done chan struct{}) {
	defer close(done)
	for {
		raw, err := conn.ReadMessage()
		if err != nil {
			fmt.Println("desconectado del servidor:", err)
			return
		}
		msgType, err := shared.PeekType(raw)
		if err != nil {
			fmt.Println("mensaje invalido recibido:", err)
			continue
		}
		switch msgType {
		case shared.TypeWelcome:
			var msg shared.WelcomeMessage
			if err := shared.DecodeMessage(raw, &msg); err != nil {
				fmt.Println("error al decodificar welcome:", err)
				continue
			}
			state.setWelcome(msg.PlayerID, msg.Config)
			fmt.Println("welcome recibido, player_id:", msg.PlayerID, "config:", msg.Config)
		case shared.TypeLobby:
			var msg shared.LobbyMessage
			if err := shared.DecodeMessage(raw, &msg); err != nil {
				fmt.Println("error al decodificar lobby:", err)
				continue
			}
			state.setLobby(msg.Players)
			fmt.Println("lobby recibido, jugadores:", msg.Players)
		case shared.TypeCountdown:
			var msg shared.CountdownMessage
			if err := shared.DecodeMessage(raw, &msg); err != nil {
				fmt.Println("error al decodificar countdown:", err)
				continue
			}
			state.setCountdown(msg.Seconds)
			fmt.Println("countdown recibido, segundos:", msg.Seconds)
		case shared.TypeStart:
			state.setStarted()
			fmt.Println("start recibido, termina la simulacion de espera en el lobby")
		case shared.TypeState:
			var msg shared.StateMessage
			if err := shared.DecodeMessage(raw, &msg); err != nil {
				fmt.Println("error al decodificar state:", err)
				continue
			}

			state.setGameState(msg.Flag, msg.Players)

		case shared.TypeGameOver:
			var msg shared.GameOverMessage
			if err := shared.DecodeMessage(raw, &msg); err != nil {
				fmt.Println("error al decodificar gameover:", err)
				continue
			}

			state.setGameOver(msg.Winner)
			fmt.Println("game over recibido. ganador:", msg.Winner)			
		case shared.TypeError:
			var msg shared.ErrorMessage
			if err := shared.DecodeMessage(raw, &msg); err != nil {
				fmt.Println("error al decodificar error:", err)
				continue
			}
			fmt.Println("error recibido del servidor:", msg.Reason)
		default:
			fmt.Println("tipo de mensaje no manejado todavia:", msgType)
		}
	}
}


func DiscoverServer() ([]DiscoveredServer, error) {
	localAddr := &net.UDPAddr{Port: 0}
	udpConn, err := net.ListenUDP("udp", localAddr)
	if err != nil {
		return nil, err
	}
	defer udpConn.Close()
	
	broadcastAddr := &net.UDPAddr{IP: net.IPv4bcast, Port: shared.DiscoveryPort}
	discover := shared.DiscoverMessage{Type: shared.TypeDiscover, V: 1}
	if err := shared.SendUDP(udpConn, broadcastAddr, discover); err != nil {
		return nil, err
	}
	udpConn.SetReadDeadline(time.Now().Add(2 * time.Second))

	found := make(map[string]DiscoveredServer)
	for {
		raw, remote, err := shared.ReadUDP(udpConn)
		if err != nil {
			break
		}
		var info shared.ServerInfoMessage
		if err := shared.DecodeMessage(raw, &info); err != nil {
			continue
		}
		key := fmt.Sprintf("%s:%d", remote.IP.String(), info.TCPPort)
		found[key] = DiscoveredServer{
			ServerInfoMessage: info,
			IP:                remote.IP.String(),
		}
	}

	servers := make([]DiscoveredServer, 0, len(found))
	for _, s := range found {
		servers = append(servers, s)
	}
	return servers, nil
}
