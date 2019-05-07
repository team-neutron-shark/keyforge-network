package kfnetwork

import (
	"errors"
	"fmt"
	"net"
	"sync"
)

// Server - This type represent our server at a high level
type Server struct {
	Clients       []*Player
	ClientMutex   sync.Mutex
	CardManager   *CardManager
	Debug         bool
	Listener      net.Listener
	ListenerMutex sync.Mutex
	Lobbies       []*Lobby
	LogQueue      chan string
	PacketQueue   chan Packet
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
		if e != nil && client != nil {
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

func (s *Server) AddLobby(creator *Player) *Lobby {
	lobby := NewLobby()
	lobby.SetID(GenerateUUID())
	lobby.AddPlayer(creator)
	lobby.SetHost(creator)

	s.Lobbies = append(s.Lobbies, lobby)

	return lobby
}

func (s *Server) RemoveLobby(lobby *Lobby) {
	lobbies := []*Lobby{}

	for _, l := range s.Lobbies {
		if l != lobby {
			lobbies = append(lobbies, l)
		}
	}

	s.Lobbies = lobbies
}

func (s *Server) FindLobbyByID(id string) (*Lobby, error) {
	for _, lobby := range s.Lobbies {
		if lobby.ID() == id {
			return lobby, nil
		}
	}

	return &Lobby{}, errors.New("no lobby found with the given ID")
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
			logEntry := fmt.Sprintf("packet received from client %s - Payload: %s", client.RemoteAddr(), string(payload))
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
	logEntry := fmt.Sprintf("Adding player %s.", player.Name)
	s.Log(logEntry)
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
	packet.Type = PacketTypeError
	packet.Message = message

	e := WritePacket(client, packet)
	return e
}

func (s *Server) SendCreateLobbyResponse(player *Player, id string) error {
	packet := CreateLobbyResponsePacket{}
	packet.Type = PacketTypeCreateLobbyResponse
	packet.Sequence = 0
	packet.ID = id

	e := WritePacket(player.Client, packet)
	return e
}
