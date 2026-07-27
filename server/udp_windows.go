//go:build windows

package server

import (
	"context"
	"net"
	"syscall"
)

// Windows does not expose SO_REUSEPORT. SO_REUSEADDR is the equivalent option
// available here; this is set before bind so rapid restarts can share UDP:8888.
func listenDiscoveryUDP() (*net.UDPConn, error) {
	listener := net.ListenConfig{Control: func(_, _ string, control syscall.RawConn) error {
		return control.Control(func(fd uintptr) {
			_ = syscall.SetsockoptInt(syscall.Handle(fd), syscall.SOL_SOCKET, syscall.SO_REUSEADDR, 1)
		})
	}}
	packetConn, err := listener.ListenPacket(context.Background(), "udp4", ":8888")
	if err != nil {
		return nil, err
	}
	return packetConn.(*net.UDPConn), nil
}
