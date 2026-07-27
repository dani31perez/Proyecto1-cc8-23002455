//go:build !windows

package server

import (
	"context"
	"net"
	"syscall"
)

// Unix platforms support both address and port reuse; the latter is attempted
// where the kernel offers it.
func listenDiscoveryUDP() (*net.UDPConn, error) {
	const soReusePort = 0x0f
	listener := net.ListenConfig{Control: func(_, _ string, control syscall.RawConn) error {
		return control.Control(func(fd uintptr) {
			_ = syscall.SetsockoptInt(int(fd), syscall.SOL_SOCKET, syscall.SO_REUSEADDR, 1)
			_ = syscall.SetsockoptInt(int(fd), syscall.SOL_SOCKET, soReusePort, 1)
		})
	}}
	packetConn, err := listener.ListenPacket(context.Background(), "udp4", ":8888")
	if err != nil {
		return nil, err
	}
	return packetConn.(*net.UDPConn), nil
}
