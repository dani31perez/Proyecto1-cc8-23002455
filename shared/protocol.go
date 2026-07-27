package shared

import (
	"bufio"
	"encoding/json"
	"errors"
	"net"
	"sync"
)

const (
	DiscoveryPort  = 8888
	MessageMaxSize = 64 * 1024 // Includes the trailing newline on TCP.
)

var ErrFrameTooLarge = errors.New("message exceeds protocol limit")

// Conn implements the protocol's one-JSON-object-per-line TCP framing.
// writeMu prevents concurrent broadcasts from interleaving bytes on the socket.
type Conn struct {
	Reader  *bufio.Reader
	Writer  net.Conn
	writeMu sync.Mutex
}

func NewConn(c net.Conn) *Conn {
	return &Conn{Reader: bufio.NewReaderSize(c, MessageMaxSize+1), Writer: c}
}

func (c *Conn) ReadMessage() ([]byte, error) {
	line, err := c.Reader.ReadSlice('\n')
	if err == bufio.ErrBufferFull {
		return nil, ErrFrameTooLarge
	}
	if err != nil {
		return nil, err
	}
	if len(line) > MessageMaxSize {
		return nil, ErrFrameTooLarge
	}
	line = line[:len(line)-1] // remove \n
	if len(line) > 0 && line[len(line)-1] == '\r' {
		line = line[:len(line)-1] // tolerate CRLF, as permitted by the standard
	}
	return line, nil
}

func (c *Conn) WriteMessage(v interface{}) error {
	data, err := json.Marshal(v)
	if err != nil {
		return err
	}
	data = append(data, '\n')
	if len(data) > MessageMaxSize {
		return ErrFrameTooLarge
	}
	c.writeMu.Lock()
	defer c.writeMu.Unlock()
	for len(data) > 0 {
		n, err := c.Writer.Write(data)
		if err != nil {
			return err
		}
		data = data[n:]
	}
	return nil
}

func DecodeMessage(raw []byte, v interface{}) error { return json.Unmarshal(raw, v) }

func PeekType(raw []byte) (string, error) {
	var message TypeOnly
	if err := json.Unmarshal(raw, &message); err != nil {
		return "", err
	}
	if message.Type == "" {
		return "", errors.New("message has no type")
	}
	return message.Type, nil
}

func SendUDP(conn *net.UDPConn, addr *net.UDPAddr, v interface{}) error {
	data, err := json.Marshal(v)
	if err != nil {
		return err
	}
	if len(data) > MessageMaxSize {
		return ErrFrameTooLarge
	}
	_, err = conn.WriteToUDP(data, addr)
	return err
}

func ReadUDP(conn *net.UDPConn) ([]byte, *net.UDPAddr, error) {
	buf := make([]byte, MessageMaxSize+1)
	n, addr, err := conn.ReadFromUDP(buf)
	if err != nil {
		return nil, nil, err
	}
	if n > MessageMaxSize {
		return nil, addr, ErrFrameTooLarge
	}
	return buf[:n], addr, nil
}
