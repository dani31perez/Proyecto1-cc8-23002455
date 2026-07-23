package server
import (
	"Proyecto1-cc8-23002455/shared"
	"fmt"
	"net"
	"os"
)
func Run() {
	tcpPort := 8889
	tcpListener, err := net.Listen("tcp", fmt.Sprintf(":%d", tcpPort))
	if err != nil {
		fmt.Println("error al abrir socket TCP:", err)
		os.Exit(1)
	}
	l := newLobby()
	go runDiscovery(tcpPort)
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

func runDiscovery(tcpPort int) {
	addr := &net.UDPAddr{Port: shared.DiscoveryPort, IP: net.IPv4zero}
	udpConn, err := net.ListenUDP("udp", addr)
	if err != nil {
		fmt.Println("error al abrir socket UDP:", err)
		return
	}
	for {
		raw, remote, err := shared.ReadUDP(udpConn)
		if err != nil {
			fmt.Println("error al leer UDP:", err)
			continue
		}
		msgType, err := shared.PeekType(raw)
		if err != nil {
			continue
		}
		if msgType != shared.TypeDiscover {
			continue
		}
		info := shared.ServerInfoMessage{Type: shared.TypeServerInfo, V: 1, Name: "servidor-ctf", TCPPort: tcpPort, State: "lobby", Players: 0}
		shared.SendUDP(udpConn, remote, info)
	}
}
func handleClient(netConn net.Conn, l *lobby) {
	conn := shared.NewConn(netConn)
	defer netConn.Close()
	var p *player
	for p == nil {
		raw, err := conn.ReadMessage()
		if err != nil {
			fmt.Println("cliente desconectado antes de unirse:", err)
			return
		}
		msgType, err := shared.PeekType(raw)
		if err != nil {
			conn.WriteMessage(shared.ErrorMessage{Type: shared.TypeError, Reason: shared.ErrInvalidJSON})
			continue
		}
		if msgType != shared.TypeJoin {
			conn.WriteMessage(shared.ErrorMessage{Type: shared.TypeError, Reason: shared.ErrNotJoined})
			continue
		}
		var join shared.JoinMessage
		if err := shared.DecodeMessage(raw, &join); err != nil {
			conn.WriteMessage(shared.ErrorMessage{Type: shared.TypeError, Reason: shared.ErrInvalidField})
			continue
		}
		if join.V != 1 {
			conn.WriteMessage(shared.ErrorMessage{ Type:   shared.TypeError, Reason: shared.ErrVersionMismatch})
			return
		}
		if join.Name == "" || len(join.Name) > 20 {
			conn.WriteMessage(shared.ErrorMessage{Type: shared.TypeError, Reason: shared.ErrNameInvalid})
			continue
		}
		p = l.addPlayer(join.Name, conn)
		fmt.Println("join recibido, v:", join.V, "name:", join.Name, "asignado id:", p.id)
	}
	defer l.removePlayer(p.id)
	welcome := shared.WelcomeMessage{Type: shared.TypeWelcome, PlayerID: p.id, Config: shared.GameConfig{MapSize: 1000, CircleRadius: 300, PlayerRadius: 15, InteractRadius: 40, Speed: 200, TickRate: 20}}
	conn.WriteMessage(welcome)
	l.broadcastLobby()
	l.startCountdownOnce()
	for {
		raw, err := conn.ReadMessage()
		if err != nil {
			fmt.Println("jugador desconectado:", p.id, err)
			l.broadcastLobby()
			return
		}
		msgType, err := shared.PeekType(raw)
		if err != nil {
			conn.WriteMessage(shared.ErrorMessage{Type: shared.TypeError, Reason: shared.ErrInvalidJSON})
			continue
		}
		
		switch msgType {

			case shared.TypeInput:
				if err := conn.ValidateFields(raw, shared.InputMessage{}); err != nil {
					continue
				}
				var msg shared.InputMessage

				if err := shared.DecodeMessage(raw, &msg); err != nil {
					fmt.Println("error al decodificar input:", err)
					continue
				}

				fmt.Println(
					"input recibido de",
					p.id,
					"direccion:",
					msg.Dir,
				)

			case shared.TypeInteract:
				if err := conn.ValidateFields(raw, shared.InteractMessage{}); err != nil {
					continue
				}
				var msg shared.InteractMessage

				if err := shared.DecodeMessage(raw, &msg); err != nil {
					fmt.Println("error al decodificar interact:", err)
					continue
				}

				fmt.Println(
					"interact recibido de",
					p.id,
					"target:",
					msg,
				)

			default:
				conn.WriteMessage(shared.ErrorMessage{
					Type:   shared.TypeError,
					Reason: shared.ErrUnknownType,
				})
		}
	}
}
