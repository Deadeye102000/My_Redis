package main

import (
	"My_Redis/client"
	"context"
	"log"

	"log/slog"
	"net"
	"time"
)

const defaultListenAddr = ":5001"

type Config struct {
	ListenAddr string
}

type Server struct {
	Config
	peers     map[*Peer]bool
	ln        net.Listener
	addPeerCh chan *Peer
	quitCh    chan struct{}
	msgch     chan []byte
}

func NewServer(cfg Config) *Server {
	if len(cfg.ListenAddr) == 0 {
		cfg.ListenAddr = defaultListenAddr
	}

	return &Server{
		Config:    cfg,
		peers:     make(map[*Peer]bool),
		addPeerCh: make(chan *Peer),
		quitCh:    make(chan struct{}),
		msgch:     make(chan []byte),
	}
}

func (s *Server) Start() error {
	ln, err := net.Listen("tcp", s.ListenAddr)
	if err != nil {
		return err

	}
	s.ln = ln

	go s.loop()

	slog.Info("server running", "listenAddr", s.ListenAddr)

	return s.acceptLoop()
	// return nil
}

func (s *Server) handleRawMessafge(rawMsg []byte) error {

	cmd, err := parseCommand(string(rawMsg))
	if err != nil {
		return err
	}
	switch cmd.(type) {
	case SetCommand:
		slog.Info("Somebody wants to set a key in the hash table")
	}
	return nil
}

func (s *Server) loop() {
	for {
		select {
		case rawMsg := <-s.msgch:
			if err := s.handleRawMessafge(rawMsg); err != nil {
				slog.Error("raw message error", "err", err)
			}
			// fmt.Println(rawMsg)
		case <-s.quitCh:
			// slog.Info("server shutting down")
			return
		case peer := <-s.addPeerCh:
			s.peers[peer] = true

			// default:
			// 	fmt.Println("foo")
		}
	}
}

func (s *Server) acceptLoop() error {
	for {
		conn, err := s.ln.Accept()
		if err != nil {
			slog.Error("accept errpr", "err", err)
			continue
		}
		go s.handleConn(conn)
	}
}

func (s *Server) handleConn(conn net.Conn) {
	peer := NewPeer(conn, s.msgch)
	s.addPeerCh <- peer
	slog.Info("new peer connected", "remoteAddr", conn.RemoteAddr())

	if err := peer.readLoop(); err != nil {
		slog.Error("error reading from peer", "err", err, "remoteAddr", conn.RemoteAddr())
	}
	// peer.readLoop()
	conn.Close()
	slog.Info("peer disconnected", "remoteAddr", conn.RemoteAddr())
}

func main() {
	go func() {
		server := NewServer(Config{})
		log.Fatal(server.Start())
	}()
	time.Sleep(time.Second) // wait for the server to start

	c := client.NewClient("localhost:5001")

	// select {} // we are blocking here so the program does not exit
	// if err != nil {
	// 	log.Fatal("error creating client", err)
	// }
	if err := c.Set(context.TODO(), "foo", "bar"); err != nil {
		log.Fatal(err)
	}

	select {}
}
