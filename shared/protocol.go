package shared
import (
	"bufio"
	"encoding/json"
	"fmt"
	"net"
)
const DiscoveryPort = 8888
const MessageMaxSize = 64 * 1024
type Conn struct {
	Reader *bufio.Reader
	Writer net.Conn
}
func NewConn(c net.Conn) *Conn {
	return &Conn{Reader: bufio.NewReaderSize(c, MessageMaxSize), Writer: c}
}
func (c *Conn) ReadMessage() ([]byte, error) {
	line, err := c.Reader.ReadBytes('\n')
	if err != nil {
		return nil, err
	}
	return line[:len(line)-1], nil
}
func (c *Conn) WriteMessage(v interface{}) error {
	data, err := json.Marshal(v)
	if err != nil {
		return err
	}
	data = append(data, '\n')
	_, err = c.Writer.Write(data)
	return err
}
func PeekType(raw []byte) (string, error) {
	var t TypeOnly
	if err := json.Unmarshal(raw, &t); err != nil {
		return "", err
	}
	if t.Type == "" {
		return "", fmt.Errorf("mensaje sin campo type")
	}
	return t.Type, nil
}
func DecodeMessage(raw []byte, v interface{}) error {
	return json.Unmarshal(raw, v)
}
func SendUDP(conn *net.UDPConn, addr *net.UDPAddr, v interface{}) error {
	data, err := json.Marshal(v)
	if err != nil {
		return err
	}
	_, err = conn.WriteToUDP(data, addr)
	return err
}
func ReadUDP(conn *net.UDPConn) ([]byte, *net.UDPAddr, error) {
	buf := make([]byte, MessageMaxSize)
	n, addr, err := conn.ReadFromUDP(buf)
	if err != nil {
		return nil, nil, err
	}
	return buf[:n], addr, nil
}
