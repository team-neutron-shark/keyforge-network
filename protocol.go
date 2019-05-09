package kfnetwork

// ProtocolVersion - This needs to be incremented when the server implements
// support for new features that old clients won't support.
const ProtocolVersion = 0.01

type PacketType uint16

const (
	PacketTypeExit uint16 = iota
	PacketTypeError
	PacketTypeVersionRequest
	PacketTypeVersionResponse
	PacketTypeLoginRequest
	PacketTypeLoginResponse
	PacketTypePlayerListRequest
	PacketTypePlayerListResponse
	PacketTypeGlobalChatRequest
	PacketTypeGlobalChatResponse
	PacketTypeCreateLobbyRequest
	PacketTypeCreateLobbyResponse
	PacketTypeJoinLobbyRequest
	PacketTypeJoinLobbyResponse
	PacketTypeLeaveLobbyRequest
	PacketTypeLeaveLobbyResponse
	PacketTypeKickLobbyRequest
	PacketTypeKickLobbyResponse
	PacketTypeBanLobbyRequest
	PacketTypeBanLobbyResponse
	PacketTypeConcedeGameRequest
	PacketTypeConcedeGameResponse
	PacketTypeSelectDeckRequest
	PacketTypeSelectDeckResponse
	PacketTypeReadyRequest
	PacketTypeReadyResponse
	PacketTypeStartGameRequest
	PacketTypeStartGameResponse
	PacketTypeLeaveGameRequest
	PacketTypeLeaveGameResponse
	PacketTypeUpdateGameState
	PacketTypeMulliganRequest
	PacketTypeMulliganResponse
	PacketTypeCardPileRequest
	PacketTypeCardPileResponse
	PacketTypeDrawCardRequest
	PacketTypeDrawCardResponse
	PacketTypePlayCardRequest
	PacketTypePlayCardResponse
	PacketTypeDiscardCardRequest
	PacketTypeDiscardCardResponse
	PacketTypeUseCardRequest
	PacketTypeUseCardResponse
)

type PileType uint8

const (
	CardPileDiscard uint8 = iota
	CardPileArchive
	CardPileHand
	CardPileDraw
)
