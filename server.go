package kfnetwork

import (
	"errors"
	"fmt"
	keyforge "keyforge/game"
	"net"
	"sync"
)

const PlayerLimit = 2

// PlayerClient - This type holds both the keyforge player type along with
// the net.Conn object required for networked communication.
type PlayerClient struct {
	Client    net.Conn
	Connected bool
	keyforge.Player
}

type SpectatorClient struct {
	Client net.Conn
	Name   string
}

// Server - This type represent our server at a high level
type Server struct {
	Debug          bool
	TestString     string
	LogQueue       chan string
	PacketQueue    chan Packet
	Listener       net.Listener
	ListenerMutex  sync.Mutex
	Players        []PlayerClient
	PlayerMutex    sync.Mutex
	Spectators     []SpectatorClient
	SpectatorMutex sync.Mutex
	Running        bool
}

// NewServer - Return a pointer to a newly created server
func NewServer(address string) *Server {
	server := new(Server)
	server.Debug = true
	server.Running = true
	server.LogQueue = make(chan string, 1024)

	// Start the listen loop on the specified address
	go server.ListenLoop(address)

	if server.Debug {
		server.Log("Server successfully started.")
	}

	return server
}

// Listen - Listen for incoming connections
func (s *Server) Listen(address string) error {
	var e error

	s.Listener, e = net.Listen("tcp4", address)

	if e != nil {
		return e
	}

	if s.Debug {
		logEntry := fmt.Sprintf("Listener started on address %s.", address)
		s.Log(logEntry)
	}

	return nil
}

// ListenLoop - Listen on the specified address and accept incoming clients.
// Spins off a new goroutine to handle the incoming client.
func (s *Server) ListenLoop(address string) {
	s.Listen(address)

	for s.Running {
		client, e := s.Accept()

		// If we encounter an error, close the client.
		if e != nil {
			client.Close()
		} else {
			// Handle incoming connection.
			if s.Debug {
				logEntry := fmt.Sprintf("Client connection accepted from remote address %s.", client.RemoteAddr())
				s.Log(logEntry)
			}

			// Handle accepted client
			s.ReceiveVersionPacket(client)
			s.ReceiveLoginPacket(client)
			go s.ReadLoop(client)
		}
	}
}

// Accept - Accept incoming connections
func (s *Server) Accept() (net.Conn, error) {
	var e error

	connection, e := s.Listener.Accept()

	if e != nil {
		return connection, e
	}

	return connection, nil
}

// HandleConnection - Process incoming connections after being accepted.
func (s *Server) HandleConnection(client net.Conn) {
	//s.ClientMutex.Lock()
	//defer s.ClientMutex.Unlock()

}

func (s *Server) Stop() {
	s.Running = false

	s.ListenerMutex.Lock()
	defer s.ListenerMutex.Unlock()

	s.Listener.Close()
}

func (s *Server) Log(message string) {
	logMessage := fmt.Sprintf("[LOG] %s", message)
	s.LogQueue <- logMessage
}

func (s *Server) GetLogs() []string {
	logs := []string{}

	for len(s.LogQueue) > 0 {
		logs = append(logs, <-s.LogQueue)
	}

	return logs
}

func (s *Server) PrintLogs() {
	for _, log := range s.GetLogs() {
		fmt.Println(log)
	}
}

func (s *Server) ReadLoop(client net.Conn) {
	for s.Running {
		packet, e := ReadPacket(client)

		if e != nil {
			logEntry := fmt.Sprintf("ReadPacket: %s", e.Error())
			s.Log(logEntry)
			s.CloseConnection(client)
			return
		}

		if s.Debug {
			logEntry := fmt.Sprintf("Packet received from client %s.", client.RemoteAddr())
			s.Log(logEntry)
		}

		s.HandlePacket(client, packet)
	}
}

func (s *Server) CloseConnection(client net.Conn) {
	if s.Debug {
		logMessage := fmt.Sprintf("Closing remote connection for %s.", client.RemoteAddr())
		s.Log(logMessage)
	}

	// Mark player as disconnected. This avoids reading from the client in
	// the main read loop.
	player, e := s.FindPlayerByConnection(client)

	if e == nil {
		player.Connected = false
	}

	// Close the connection.
	client.Close()
}

func (s *Server) SendErrorPacket(client net.Conn, message string) error {
	packet := ErrorPacket{}
	packet.Sequence = 0
	packet.Message = message

	e := WritePacket(client, packet)
	return e
}

func (s *Server) FindPlayerByConnection(client net.Conn) (PlayerClient, error) {
	playerClient := PlayerClient{}

	for _, player := range s.Players {
		if client == player.Client {
			return player, nil
		}
	}

	return playerClient, errors.New("no player found")
}

func (s *Server) HandlePacket(client net.Conn, packet Packet) {
	switch packet.GetHeader().Type {
	case PacketTypeVersion:
		s.HandleVersionPacket(client, packet.(VersionPacket))
	case PacketTypeLogin:
		s.HandleLoginPacket(client, packet.(LoginPacket))
	}
}

func (s *Server) HandleVersionPacket(client net.Conn, packet VersionPacket) {
	debugString := fmt.Sprintf("HandleVersionPacket: %+v", packet)
	s.Log(debugString)
	if packet.Version != ProtocolVersion {
		if s.Debug {
			logEntry := fmt.Sprintf("Client %s sent a version packet with a mismatching version.", client.RemoteAddr())
			s.Log(logEntry)
		}
		s.SendErrorPacket(client, "Protocol version mismatch.")
		s.CloseConnection(client)
	}
}

func (s *Server) HandleLoginPacket(client net.Conn, packet LoginPacket) {
	if len(s.Players) < 2 {
		s.PlayerMutex.Lock()
		defer s.PlayerMutex.Unlock()

		player := PlayerClient{}
		player.Name = packet.Name
		s.Players = append(s.Players, player)

		if s.Debug {
			logEntry := fmt.Sprintf("Player %s has joined the game.", player.Name)
			s.Log(logEntry)
		}
	} else {
		s.SpectatorMutex.Lock()
		defer s.SpectatorMutex.Unlock()

		spectator := SpectatorClient{}
		spectator.Name = packet.Name
		s.Spectators = append(s.Spectators, spectator)

		if s.Debug {
			logEntry := fmt.Sprintf("Spectator %s has joined the game.")
			s.Log(logEntry)
		}
	}
}

func (s *Server) ReceiveLoginPacket(client net.Conn) error {
	packet, e := ReadPacket(client)

	if e != nil {
		return e
	}

	if packet.GetHeader().Type != PacketTypeLogin {
		logEntry := fmt.Sprintf("Expected login packet from %s; received packet of a different type.", client.RemoteAddr())
		s.Log(logEntry)

		s.SendErrorPacket(client, "No login packet received.")
		s.CloseConnection(client)
	}

	s.HandlePacket(client, packet)

	return nil
}

func (s *Server) ReceiveVersionPacket(client net.Conn) error {
	packet, e := ReadPacket(client)

	if e != nil {
		return e
	}

	if packet.GetHeader().Type != PacketTypeVersion {
		logEntry := fmt.Sprintf("Expected version packet from %s; received packet of a different type.", client.RemoteAddr())
		s.Log(logEntry)

		s.SendErrorPacket(client, "No version packet received.")
		s.CloseConnection(client)
	}

	s.HandlePacket(client, packet)

	return nil
}
