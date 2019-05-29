package kfnetwork

type Observer interface {
	Notify(Event)
}

type Subject interface {
	AddObserver(Observer)
	RemoveObserver(Observer)
	NotifyObservers(Event)
}
