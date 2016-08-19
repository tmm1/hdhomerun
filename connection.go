package hdhomerun

import (
	"io"
	"net"
)

type Connection interface {
	Send(*Packet) error
	Recv() (*Packet, error)
}

type Connectable interface {
	Connect() error
}

type Closeable interface {
	Close() error
}

type Addressable interface {
	RemoteAddr() net.Addr
}

type IOConnection struct {
	encoder *Encoder
	decoder *Decoder
}

func NewIOConnection(rw io.ReadWriter) *IOConnection {
	return &IOConnection{
		encoder: NewEncoder(rw),
		decoder: NewDecoder(rw),
	}
}

func (conn *IOConnection) Send(p *Packet) error {
	return conn.encoder.Encode(p)
}

func (conn *IOConnection) Recv() (p *Packet, err error) {
	return conn.decoder.Next()
}

type TCPConnection struct {
	*net.TCPConn
	*IOConnection
	addr *net.TCPAddr
}

func NewTCPConnection(addr *net.TCPAddr) *TCPConnection {
	return &TCPConnection{
		addr: addr,
	}
}

func (conn *TCPConnection) Connect() (err error) {
	conn.TCPConn, err = net.DialTCP("tcp", nil, conn.addr)
	conn.IOConnection = NewIOConnection(conn.TCPConn)
	return err
}

func (conn *TCPConnection) RemoteAddr() net.Addr {
	return conn.addr
}
