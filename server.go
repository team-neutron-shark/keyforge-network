package kfnetwork

import (
	"fmt"
	"net"
	"sync"
)

// Server - This type represent our server at a high level
type Server struct {
	ClientMutex   sync.Mutex
	CardManager   *CardManager
	Debug         bool
	Listener      net.Listener
	ListenerMutex sync.Mutex
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
	server.AddObserver(Events())

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
			observers = append(observers, o)
		}
	}

	s.observers = observers
}

// NotifyObservers - Notifies observers that a network event has occured.
func (s *Server) NotifyObservers(event Event) {
	for _, observer := range s.observers {
		observer.Notify(event)
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

func (s *Server) Stop() {
	s.Running = false
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
			player, e := Players().FindPlayerByConnection(client)
			if e == nil {
				Players().RemovePlayer(player)
			}
			s.CloseConnection(client)
			return
		}

		s.NotifyObservers(NetworkEvent{connection: &client, packet: &packet})
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

func (s *Server) SendErrorPacket(client net.Conn, message string) error {
	packet := ErrorPacket{}
	packet.Type = PacketTypeError
	packet.Message = message

	e := WritePacket(client, packet)
	return e
}

func (s *Server) SendVersionResponse(client net.Conn) error {
	packet := VersionPacket{}
	packet.Type = PacketTypeVersionResponse
	packet.Version = ProtocolVersion

	e := WritePacket(client, packet)
	return e
}

func (s *Server) SendCreateLobbyResponse(player *Player, id string) error {
	packet := CreateLobbyResponsePacket{}
	packet.Type = PacketTypeCreateLobbyResponse
	packet.ID = id

	e := WritePacket(player.Client, packet)
	return e
}

func (s *Server) SendPlayerListResponse(player *Player, list PlayerList) error {
	packet := PlayerListResponsePacket{}
	packet.Type = PacketTypePlayerListResponse
	packet.Count = list.Count
	packet.Players = list.Players

	e := WritePacket(player.Client, packet)
	return e
}

func (s *Server) SendLobbyListResponse(player *Player, list LobbyList) error {
	packet := LobbyListResponsePacket{}
	packet.Type = PacketTypeLobbyListResponse
	packet.Count = list.Count
	packet.Lobbies = list.Lobbies

	e := WritePacket(player.Client, packet)
	return e
}

func (s *Server) SendGlobalChatResponse(player *Player, name, message string) error {
	packet := GlobalChatResponsePacket{}
	packet.Type = PacketTypeGlobalChatResponse
	packet.Name = name
	packet.Message = message

	e := WritePacket(player.Client, packet)
	return e
}

func (s *Server) SendJoinLobbyResponse(player *Player, name, id string, success bool) error {
	packet := JoinLobbyResponsePacket{}
	packet.Type = PacketTypeJoinLobbyResponse
	packet.Name = name
	packet.ID = id
	packet.Success = success

	e := WritePacket(player.Client, packet)
	return e
}

func (s *Server) SendLeaveLobbyResponse(player *Player, name, id string, success bool) error {
	packet := JoinLobbyResponsePacket{}
	packet.Type = PacketTypeLeaveLobbyResponse
	packet.Name = name
	packet.ID = id
	packet.Success = success

	e := WritePacket(player.Client, packet)
	return e
}

func (s *Server) SendLobbyKickResponse(player *Player, target string, success bool) error {
	packet := LobbyKickResponsePacket{}
	packet.Type = PacketTypeKickLobbyResponse
	packet.Target = target
	packet.Success = success

	e := WritePacket(player.Client, packet)
	return e
}

func (s *Server) SendLobbyChatResponse(player *Player, name string, message string) error {
	packet := LobbyChatResponsePacket{}
	packet.Type = PacketTypeLobbyChatResponse
	packet.Name = name
	packet.Message = message

	e := WritePacket(player.Client, packet)
	return e
}
