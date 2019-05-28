package kfnetwork

type Observer interface {
	Notify(Packet)
}

type Subject interface {
	AddObserver(Observer)
	RemoveObserver(Observer)
	Notify(Packet)
}
