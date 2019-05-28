package kfnetwork

import "net"

type Observer interface {
	Notify(Event)
}

type Subject interface {
	AddObserver(Observer)
	RemoveObserver(Observer)
	Notify(Event)
}

type Event interface {
	GetClient() *net.Conn
	GetPacket() *Packet
}

type BaseEvent struct {
	connection *net.Conn
	packet     *Packet
}

type NetworkEvent struct {
	BaseEvent
}

func (b BaseEvent) GetClient() *net.Conn {
	return b.connection
}

func (b BaseEvent) GetPacket() *Packet {
	return b.packet
}
