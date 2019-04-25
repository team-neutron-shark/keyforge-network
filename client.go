package kfnetwork

import "net"

type Client struct {
	Connection net.Conn
	Sequence   uint16
}

func NewClient() *Client {
	client := new(Client)
	return client
}

func (c *Client) Connect(address string) error {
	var e error
	c.Connection, e = net.Dial("tcp", address)

	if e != nil {
		return e
	}

	return nil
}

func (c *Client) SendVersion() error {
	packet := VersionPacket{}
	packet.Sequence = c.Sequence
	packet.Type = PacketTypeVersion
	packet.Version = ProtocolVersion

	e := WritePacket(c.Connection, packet)
	c.Sequence++
	return e
}

func (c *Client) SendExit() error {
	packet := ExitPacket{}
	packet.Sequence = c.Sequence
	packet.Type = PacketTypeExit

	e := WritePacket(c.Connection, packet)
	c.Sequence++
	return e
}

func (c *Client) SendLogin() error {
	packet := LoginPacket{}
	packet.Sequence = c.Sequence
	packet.Type = PacketTypeLogin
	packet.Name = "test"
	packet.ID = GenerateUUID()

	e := WritePacket(c.Connection, packet)
	return e
}
