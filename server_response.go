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
	case PacketTypeGlobalChatRequest:
		s.HandleGlobalChatRequest(client, packet.(GlobalChatRequestPacket))
	case PacketTypePlayerListRequest:
		s.HandlePlayerListRequest(client, packet.(PlayerListRequestPacket))
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

func (s *Server) HandlePlayerListRequest(client net.Conn, packet PlayerListRequestPacket) error {
	player, e := s.FindPlayerByConnection(client)

	if e != nil {
		if s.Debug {
			logEntry := fmt.Sprintf("HandlePlayerListRequest: %s", e.Error())
			s.Log(logEntry)
		}
		return e
	}

	playerList := PlayerList{}

	for _, p := range s.Clients {
		entry := PlayerListEntry{}
		entry.ID = p.ID
		entry.Name = p.Name

		playerList.Players = append(playerList.Players, entry)
	}

	playerList.Count = uint(len(playerList.Players))

	e = s.SendPlayerListResponse(player, playerList)

	if e != nil {
		return e
	}

	logEntry := fmt.Sprintf("Player %s requested the player list", player.Name)
	s.Log(logEntry)
	return nil
}

func (s *Server) HandleLobbyChatRequest(client net.Conn, packet LobbyChatRequestPacket) error {
	return nil
}

func (s *Server) HandleGlobalChatRequest(client net.Conn, packet GlobalChatRequestPacket) error {
	player, e := s.FindPlayerByConnection(client)

	if e != nil {
		return e
	}

	for _, p := range s.Clients {
		s.SendGlobalChatResponse(p, player.Name, packet.Message)
	}

	logEntry := fmt.Sprintf("(Global Chat) %s: %s", player.Name, packet.Message)
	s.Log(logEntry)
	return nil
}
