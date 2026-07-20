package main
import (
	"Proyecto1-cc8-23002455/shared"
	"fmt"
	"net"
	"os"
	"time"
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
	join := shared.JoinMessage{Type: shared.TypeJoin, V: 1, Name: "jugador1"}
	conn.WriteMessage(join)
	for {
		raw, err := conn.ReadMessage()
		if err != nil {
			fmt.Println("desconectado del servidor:", err)
			return
		}
		msgType, err := shared.PeekType(raw)
		if err != nil {
			continue
		}
		fmt.Println("mensaje del servidor:", msgType)
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
