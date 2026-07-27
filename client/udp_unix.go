//go:build !windows

package client

import (
	"context"
	"net"
	"syscall"
)

func openDiscoverySocket() (*net.UDPConn, error) {
	listener := net.ListenConfig{Control: func(_, _ string, control syscall.RawConn) error {
		return control.Control(func(fd uintptr) { _ = syscall.SetsockoptInt(int(fd), syscall.SOL_SOCKET, syscall.SO_BROADCAST, 1) })
	}}
	packetConn, err := listener.ListenPacket(context.Background(), "udp4", ":0")
	if err != nil {
		return nil, err
	}
	return packetConn.(*net.UDPConn), nil
}
