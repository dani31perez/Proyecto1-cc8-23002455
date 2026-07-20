package main
import (
	"Proyecto1-cc8-23002455/shared"
	"fmt"
	"net"
	"os"
)
func main() {
	tcpPort := 8889
	tcpListener, err := net.Listen("tcp", fmt.Sprintf(":%d", tcpPort))
	if err != nil {
		fmt.Println("error al abrir socket TCP:", err)
		os.Exit(1)
	}
	go runDiscovery(tcpPort)
	fmt.Println("servidor escuchando TCP en puerto", tcpPort)
	for {
		conn, err := tcpListener.Accept()
		if err != nil {
			fmt.Println("error al aceptar conexion:", err)
			continue
		}
		go handleClient(conn)
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
func handleClient(netConn net.Conn) {
	conn := shared.NewConn(netConn)
	defer netConn.Close()
	for {
		raw, err := conn.ReadMessage()
		if err != nil {
			fmt.Println("cliente desconectado:", err)
			return
		}
		msgType, err := shared.PeekType(raw)
		if err != nil {
			conn.WriteMessage(shared.ErrorMessage{Type: shared.TypeError, Reason: "INVALID_JSON"})
			continue
		}
		fmt.Println("mensaje recibido:", msgType)
	}
}
