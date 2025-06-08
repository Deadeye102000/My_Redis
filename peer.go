package main

import (
	// "fmt"
	"net"
	// "log/slog"
)

const (
	IAC      = 255
	DO       = 253
	LINEMODE = 34
)

type Peer struct {
	conn  net.Conn
	msgCh chan []byte
}

func filterTelnetNegotiation(data []byte) []byte {
	result := make([]byte, 0, len(data))
	i := 0
	for i < len(data) {
		if data[i] == IAC && i+2 < len(data) {
			// Skip IAC + CMD + OPTION
			i += 3
		} else {
			result = append(result, data[i])
			i++
		}
	}
	return result
}

func NewPeer(conn net.Conn, msgCh chan []byte) *Peer {
	return &Peer{
		conn:  conn,
		msgCh: msgCh,
	}
}
func (p *Peer) readLoop() error {
	_, _ = p.conn.Write([]byte{IAC, DO, LINEMODE})
	buf := make([]byte, 1024)
	for {
		n, err := p.conn.Read(buf)
		if err != nil {
			// slog.Error("read error", "err", err)
			// p.conn.Close()
			return err
		}
		// fmt.Println(string(buf[:n]))
		msgBuf := make([]byte, n)
		cleaned := filterTelnetNegotiation(buf[:n])

		copy(msgBuf, cleaned)
		// msg := string(buf[:n])

		p.msgCh <- msgBuf
	}
}
