package client

import (
	"Proyecto1-cc8-23002455/shared"
	"fmt"
	"net"
	"sort"
	"strconv"
	"strings"
	"time"
)

var CurrentConn *shared.Conn

func Run(server DiscoveredServer, name string) {
	tcpConn, err := net.Dial("tcp", net.JoinHostPort(server.IP, strconv.Itoa(server.TCPPort)))
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
	_ = conn.WriteMessage(shared.JoinMessage{Type: shared.TypeJoin, V: shared.ProtocolVersion, Name: name})
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
			continue
		}
		switch msgType {
		case shared.TypeWelcome:
			var msg shared.WelcomeMessage
			if shared.DecodeMessage(raw, &msg) == nil {
				state.setWelcome(msg.PlayerID, msg.Config)
			}
		case shared.TypeLobby:
			var msg shared.LobbyMessage
			if shared.DecodeMessage(raw, &msg) == nil {
				state.setLobby(msg.Players)
			}
		case shared.TypeCountdown:
			var msg shared.CountdownMessage
			if shared.DecodeMessage(raw, &msg) == nil {
				state.setCountdown(msg.Seconds)
			}
		case shared.TypeStart:
			state.setStarted()
		case shared.TypeState:
			var msg shared.StateMessage
			if shared.DecodeMessage(raw, &msg) == nil {
				state.setGameState(msg.Flag, msg.Players)
			}
		case shared.TypeGameOver:
			var msg shared.GameOverMessage
			if shared.DecodeMessage(raw, &msg) == nil {
				state.setGameOver(msg.Winner)
			}
		case shared.TypeError:
			var msg shared.ErrorMessage
			if shared.DecodeMessage(raw, &msg) == nil {
				Notify.Show("ERROR RECIBIDO DEL SERVIDOR:", msg.Reason)
			}
		}
	}
}

// DiscoverServer sends both protocol-required broadcasts: the limited
// broadcast and every IPv4 subnet broadcast available on this host.
func DiscoverServer() ([]DiscoveredServer, error) {
	targets := map[string]*net.UDPAddr{}
	addTarget := func(ip net.IP) { targets[ip.String()] = &net.UDPAddr{IP: ip, Port: shared.DiscoveryPort} }
	addTarget(net.IPv4bcast)
	for _, ip := range subnetBroadcasts() {
		addTarget(ip)
	}
	addresses := make([]*net.UDPAddr, 0, len(targets))
	for _, addr := range targets {
		addresses = append(addresses, addr)
	}
	return discover(addresses)
}

// DirectServer accepts either IP:puerto or an IP plus the port from the UI.
func DirectServer(address, portText string) (DiscoveredServer, error) {
	address = strings.TrimSpace(address)
	portText = strings.TrimSpace(portText)
	host := address
	port := 0
	if parsedHost, parsedPort, err := net.SplitHostPort(address); err == nil {
		host = parsedHost
		port, err = strconv.Atoi(parsedPort)
		if err != nil {
			return DiscoveredServer{}, err
		}
	} else {
		var err error
		port, err = strconv.Atoi(portText)
		if err != nil {
			return DiscoveredServer{}, fmt.Errorf("puerto TCP inválido")
		}
	}
	if net.ParseIP(host) == nil || port < 1 || port > 65535 {
		return DiscoveredServer{}, fmt.Errorf("dirección IP:puerto inválida")
	}
	return DiscoveredServer{ServerInfoMessage: shared.ServerInfoMessage{Type: shared.TypeServerInfo, V: shared.ProtocolVersion, TCPPort: port}, IP: host}, nil
}

func discover(addresses []*net.UDPAddr) ([]DiscoveredServer, error) {
	udpConn, err := openDiscoverySocket()
	if err != nil {
		return nil, err
	}
	defer udpConn.Close()
	message := shared.DiscoverMessage{Type: shared.TypeDiscover, V: shared.ProtocolVersion}
	for _, address := range addresses {
		if err := shared.SendUDP(udpConn, address, message); err != nil {
			return nil, err
		}
	}
	_ = udpConn.SetReadDeadline(time.Now().Add(2 * time.Second))
	found := make(map[string]DiscoveredServer)
	for {
		raw, remote, err := shared.ReadUDP(udpConn)
		if err != nil {
			break
		}
		var info shared.ServerInfoMessage
		if shared.DecodeMessage(raw, &info) != nil || info.Type != shared.TypeServerInfo || info.V != shared.ProtocolVersion || info.TCPPort < 1 || info.TCPPort > 65535 || (info.State != "lobby" && info.State != "playing") {
			continue
		}
		key := net.JoinHostPort(remote.IP.String(), strconv.Itoa(info.TCPPort))
		found[key] = DiscoveredServer{ServerInfoMessage: info, IP: remote.IP.String()}
	}
	servers := make([]DiscoveredServer, 0, len(found))
	for _, server := range found {
		servers = append(servers, server)
	}
	sort.Slice(servers, func(i, j int) bool {
		return net.JoinHostPort(servers[i].IP, strconv.Itoa(servers[i].TCPPort)) < net.JoinHostPort(servers[j].IP, strconv.Itoa(servers[j].TCPPort))
	})
	return servers, nil
}

func subnetBroadcasts() []net.IP {
	var result []net.IP
	interfaces, err := net.Interfaces()
	if err != nil {
		return result
	}
	for _, iface := range interfaces {
		if iface.Flags&net.FlagUp == 0 || iface.Flags&net.FlagLoopback != 0 {
			continue
		}
		addresses, err := iface.Addrs()
		if err != nil {
			continue
		}
		for _, address := range addresses {
			ipNet, ok := address.(*net.IPNet)
			if !ok {
				continue
			}
			ip := ipNet.IP.To4()
			mask := ipNet.Mask
			if ip == nil || len(mask) != net.IPv4len {
				continue
			}
			broadcast := make(net.IP, net.IPv4len)
			for i := range broadcast {
				broadcast[i] = ip[i] | ^mask[i]
			}
			result = append(result, broadcast)
		}
	}
	return result
}
