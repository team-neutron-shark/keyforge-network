package kfnetwork

import (
	"errors"
	"fmt"
	"net"
	"sync"
)

// Server - This type represent our server at a high level
type Server struct {
	Debug         bool
	TestString    string
	LogQueue      chan string
	PacketQueue   chan Packet
	Listener      net.Listener
	ListenerMutex sync.Mutex
	Clients       []*Player
	ClientMutex   sync.Mutex
	CardManager   *CardManager
	CardMutex     sync.Mutex
	Running       bool
}

// NewServer - Return a pointer to a newly created server
func NewServer(address string) *Server {
	server := new(Server)
	server.Debug = true
	server.Running = true
	server.LogQueue = make(chan string, 1024)

	server.CardManager = NewCardManager()
	e := server.CardManager.LoadFromFile("data/cards.json")

	if e != nil && server.Debug {
		logEntry := fmt.Sprintf("error loading card data: %s", e.Error())
		server.Log(logEntry)
	}

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
			if s.Debug {
				logEntry := fmt.Sprintf("Client connection accepted from remote address %s.", client.RemoteAddr())
				s.Log(logEntry)
			}

			// Handle accepted client
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
	s.ClientMutex.Lock()
	defer s.ClientMutex.Unlock()

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
			payload, _ := packet.GetPayload()
			logEntry := fmt.Sprintf("packet received from client %s!\nPayload: %s", client.RemoteAddr(), string(payload))
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
	client.Close()
}

func (s *Server) FindPlayerByConnection(connection net.Conn) (*Player, error) {
	for _, player := range s.Clients {
		if player.Client == connection {
			return player, nil
		}
	}

	return &Player{}, errors.New("could not find any players with the given connection")
}

func (s *Server) PlayerExists(client *Player) bool {
	s.ClientMutex.Lock()
	defer s.ClientMutex.Unlock()

	for _, player := range s.Clients {
		if player == client {
			return true
		}
	}

	return false
}

func (s *Server) AddPlayer(player *Player) {
	if !s.PlayerExists(player) {
		s.ClientMutex.Lock()
		defer s.ClientMutex.Unlock()
		s.Clients = append(s.Clients, player)
	}
}

func (s *Server) RemovePlayer(player *Player) {
	clients := []*Player{}

	s.ClientMutex.Lock()
	defer s.ClientMutex.Unlock()

	for _, client := range s.Clients {
		if client != player {
			clients = append(clients, client)
		}
	}

	s.Clients = clients
}

func (s *Server) SendErrorPacket(client net.Conn, message string) error {
	packet := ErrorPacket{}
	packet.Sequence = 0
	packet.Message = message

	e := WritePacket(client, packet)
	return e
}

func (s *Server) HandlePacket(client net.Conn, packet Packet) {
	switch packet.GetHeader().Type {
	case PacketTypeVersionRequest:
		s.HandleVersionRequest(client, packet.(VersionPacket))
	case PacketTypeLoginRequest:
		s.HandleLoginRequest(client, packet.(LoginRequestPacket))
	}
}

func (s *Server) HandleVersionRequest(client net.Conn, packet VersionPacket) {
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

func (s *Server) HandleLoginRequest(client net.Conn, packet LoginRequestPacket) {
	//TODO - add authentication logic; for now assume login succeeds
	player := NewPlayer()
	player.Name = packet.Name
	player.ID = packet.ID
	player.Client = client

	s.AddPlayer(player)
}

func (s *Server) HandleExitRequest(client net.Conn, packet ExitPacket) error {
	player, e := s.FindPlayerByConnection(client)

	if e != nil {
		return e
	}

	s.CloseConnection(player.Client)
	s.RemovePlayer(player)
	return nil
}
