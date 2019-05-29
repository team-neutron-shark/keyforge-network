package kfnetwork

import "net"

type Event interface{}

type NetworkEvent struct {
	connection *net.Conn
	packet     *Packet
}
