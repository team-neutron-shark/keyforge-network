package kfnetwork

type Lobby struct {
	id      string
	host    *Player
	players []*Player
	game    *Game
	name    string
}

func NewLobby() *Lobby {
	lobby := new(Lobby)
	return lobby
}

func (l *Lobby) Players() []*Player {
	return l.players
}

func (l *Lobby) AddPlayer(player *Player) {
	if len(l.players) < 2 && !l.PlayerExists(player) {
		l.players = append(l.players, player)
	}
}

func (l *Lobby) RemovePlayer(player *Player) {
	players := []*Player{}

	for _, p := range l.players {
		if p != player {
			players = append(players, p)
		}
	}

	l.players = players
}

func (l *Lobby) PlayerExists(player *Player) bool {
	for _, p := range l.players {
		if p == player {
			return true
		}
	}

	return false
}

func (l *Lobby) ID() string {
	return l.id
}

func (l *Lobby) SetID(id string) {
	l.id = id
}

func (l *Lobby) Host() *Player {
	return l.host
}

func (l *Lobby) SetHost(player *Player) {
	l.host = player
}
