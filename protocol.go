package kfnetwork

// ProtocolVersion - This needs to be incremented when the server implements
// support for new features that old clients won't support.
const ProtocolVersion = 0.01

const (
	PacketTypeExit = iota
	PacketTypeError
	PacketTypeVersion
	PacketTypeLogin
	PacketTypeSendDeck
	PacketTypeSendCardd
)
