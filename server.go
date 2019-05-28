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
	observers     []Observer
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
		Logger().Warn(logEntry)
	}

	// Add Observers
	server.AddObserver(Logger())

	// Start the listen loop on the specified address
	go server.ListenLoop(address)

	if server.Debug {
		Logger().Log("Server successfully started.")
	}

	return server
}

// AddObserver - Adds an observer to the list of observers. Observers would be
// types such as loggers, packet responders, and other "classes" that need to
// concern themselves with incoming packet events.
func (s *Server) AddObserver(observer Observer) {
	s.observers = append(s.observers, observer)
}

// RemoveObserver - Removes an observer from the list of observers. Generally
// one wouldn't want to do this in practice, but it's there if we need it.
func (s *Server) RemoveObserver(observer Observer) {
	observers := []Observer{}

	for _, o := range s.observers {
		if o != observer {
			observers = append(observers)
		}
	}

	s.observers = observers
}

// Notify - Notifies observers that a network event has occured.
func (s *Server) Notify(packet Packet) {
	for _, observer := range s.observers {
		observer.Notify(packet)
	}
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
		Logger().Log(logEntry)
	}

	return nil
}

// ListenLoop - Listen on the specified address and accept incoming clients.
// Spins off a new goroutine to handle the incoming client.
func (s *Server) ListenLoop(address string) {
	e := s.Listen(address)

	// If we can't listen on the port print an error, halt the server, and return.
	if e != nil {
		logEntry := fmt.Sprintf("Unable to listen on address %s: %s", address, e.Error())
		Logger().Error(logEntry)
		s.Running = false
		return
	}

	for s.Running {
		client, e := s.Accept()

		// If we encounter an error, close the client.
		if e != nil && client != nil {
			client.Close()
		} else {
			if s.Debug {
				logEntry := fmt.Sprintf("Client connection accepted from remote address %s.", client.RemoteAddr())
				Logger().Log(logEntry)
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

func (s *Server) AddLobby(creator *Player, name string) *Lobby {
	lobby := NewLobby()
	lobby.SetID(GenerateUUID())
	lobby.AddPlayer(creator)
	lobby.SetHost(creator)
	lobby.SetName(name)

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

func (s *Server) FindLobbyByName(name string) (*Lobby, error) {
	for _, lobby := range s.Lobbies {
		if lobby.name == name {
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
	Logger().Log(message)
}

func (s *Server) ReadLoop(client net.Conn) {
	for s.Running {
		packet, e := ReadPacket(client)

		if e != nil {
			logEntry := fmt.Sprintf("ReadPacket: %s", e.Error())
			Logger().Error(logEntry)
			player, e := s.FindPlayerByConnection(client)
			if e == nil {
				s.RemovePlayer(player)
			}
			s.CloseConnection(client)
			return
		}

		s.Notify(packet)
		//if s.Debug {
		//	payload, _ := GetPacketPayload(packet)
		//	logEntry := fmt.Sprintf("packet received from client %s - Payload: %s", client.RemoteAddr(), string(payload))
		//	Logger().Log(logEntry)
		//}
		s.HandlePacket(client, packet)
	}
}

func (s *Server) CloseConnection(client net.Conn) {
	if s.Debug {
		logMessage := fmt.Sprintf("Closing remote connection for %s.", client.RemoteAddr())
		Logger().Log(logMessage)
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
	logEntry := fmt.Sprintf("Player %s logged in.", player.Name)
	Logger().Log(logEntry)

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

func (s *Server) SendPlayerListResponse(player *Player, list PlayerList) error {
	packet := PlayerListResponsePacket{}
	packet.Type = PacketTypePlayerListResponse
	packet.Sequence = 0
	packet.Count = list.Count
	packet.Players = list.Players

	e := WritePacket(player.Client, packet)
	return e
}

func (s *Server) SendLobbyListResponse(player *Player, list LobbyList) error {
	packet := LobbyListResponsePacket{}
	packet.Type = PacketTypeLobbyListResponse
	packet.Sequence = 0
	packet.Count = list.Count
	packet.Lobbies = list.Lobbies

	e := WritePacket(player.Client, packet)
	return e
}

func (s *Server) SendGlobalChatResponse(player *Player, name, message string) error {
	packet := GlobalChatResponsePacket{}
	packet.Type = PacketTypeGlobalChatResponse
	packet.Sequence = 0
	packet.Name = name
	packet.Message = message

	e := WritePacket(player.Client, packet)
	return e
}

func (s *Server) SendJoinLobbyResponse(player *Player, name, id string, success bool) error {
	packet := JoinLobbyResponsePacket{}
	packet.Type = PacketTypeJoinLobbyResponse
	packet.Sequence = 0
	packet.Name = name
	packet.ID = id
	packet.Success = success

	e := WritePacket(player.Client, packet)
	return e
}
