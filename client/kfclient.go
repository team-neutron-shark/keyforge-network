package main

import (
	"fmt"
	"net"
	"server"
)

func main() {
	packet := server.VersionPacket{}
	packet.Version = 0.01
	packet.Sequence = 1
	packet.Type = server.PacketTypeVersion

	connection, e := net.Dial("tcp", ":8888")

	if e != nil {
		fmt.Println(e)
		return
	}
	fmt.Println("writing packet")
	e = server.WritePacket(connection, packet)

	if e != nil {
		fmt.Println(e)
	}

}
