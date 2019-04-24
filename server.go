package server

import (
	"fmt"
	keyforge "keyforge/game"
	"net"
	"sync"
)

// PlayerClient - This type holds both the keyforge player type along with
// the net.Conn object required for networked communication.
type PlayerClient struct {
	Client net.Conn
	keyforge.Player
}

// Server - This type represent our server at a high level
type Server struct {
	Debug         bool
	TestString    string
	LogQueue      chan string
	PacketQueue   chan Packet
	Listener      net.Listener
	ListenerMutex sync.Mutex
	Clients       []PlayerClient
	ClientMutex   sync.Mutex
	Running       bool
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
		fmt.Println(packet)
		if s.Debug {
			payload, _ := packet.GetPayload()
			logEntry := fmt.Sprintf("packet received from client %s!\nPayload: %s", client.RemoteAddr(), string(payload))
			s.Log(logEntry)
		}
	}
}

func (s *Server) CloseConnection(client net.Conn) {
	if s.Debug {
		logMessage := fmt.Sprintf("Closing remote connection for %s.", client.RemoteAddr())
		s.Log(logMessage)
	}
	client.Close()
}
