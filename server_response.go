package kfnetwork

import (
	"fmt"
	"net"
)

func (s *Server) HandlePacket(client net.Conn, packet Packet) {
	switch packet.GetHeader().Type {
	case PacketTypeVersionRequest:
		s.HandleVersionRequest(client, packet.(VersionPacket))
	case PacketTypeLoginRequest:
		s.HandleLoginRequest(client, packet.(LoginRequestPacket))
	case PacketTypeCreateLobbyRequest:
		s.HandleCreateLobbyRequest(client, packet.(CreateLobbyRequestPacket))
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

func (s *Server) HandleCreateLobbyRequest(client net.Conn, packet CreateLobbyRequestPacket) error {
	player, e := s.FindPlayerByConnection(client)

	if e != nil {
		if s.Debug {
			logEntry := fmt.Sprintf("HandleCreateLobbyRequest: %s", e.Error())
			s.Log(logEntry)
		}

		return e
	}

	lobby := s.AddLobby(player, packet.Name)

	logEntry := fmt.Sprintf("Player %s created lobby %s (%s)", player.Name, lobby.name, lobby.ID())
	s.Log(logEntry)

	e = s.SendCreateLobbyResponse(player, lobby.ID())
	return e
}

func (s *Server) HandleLobbyChatRequest(client net.Conn, packet LobbyChatRequestPacket) error {
	return nil
}
