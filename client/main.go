package main
import (
	"Proyecto1-cc8-23002455/shared"
	"fmt"
	"net"
	"os"
	"time"
	"bufio"
)
func main() {
	server, err := discoverServer()
	if err != nil {
		fmt.Println("no se encontro servidor por broadcast, indique IP manualmente en el codigo")
		os.Exit(1)
	}
	fmt.Println("servidor encontrado:", server.Name, "en puerto", server.TCPPort)
	tcpConn, err := net.Dial("tcp", fmt.Sprintf("%s:%d", server.IP, server.TCPPort))
	if err != nil {
		fmt.Println("error al conectar por TCP:", err)
		os.Exit(1)
	}
	conn := shared.NewConn(tcpConn)
	state := newClientState()
	done := make(chan struct{})
	go readLoop(conn, state, done)
	stdin := bufio.NewScanner(os.Stdin)
	fmt.Print("ingresa tu nombre: ")
	stdin.Scan()
	name := stdin.Text()
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
			fmt.Println("start recibido, termina la simulacion de espera en el lobby")
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

type discoveredServer struct {
	shared.ServerInfoMessage
	IP string
}
func discoverServer() (*discoveredServer, error) {
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
	raw, remote, err := shared.ReadUDP(udpConn)
	if err != nil {
		return nil, err
	}
	var info shared.ServerInfoMessage
	if err := shared.DecodeMessage(raw, &info); err != nil {
		return nil, err
	}
	return &discoveredServer{ServerInfoMessage: info, IP: remote.IP.String()}, nil
}
