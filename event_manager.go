package kfnetwork

import "sync"

var eventManagerOnce sync.Once
var eventManagerSingleton *EventManager

type EventManager struct {
	observers []Observer
}

func Events() *EventManager {
	eventManagerOnce.Do(func() {
		eventManagerSingleton = new(EventManager)
		eventManagerSingleton.AddObserver(Logger())
	})

	return eventManagerSingleton
}

// AddObserver - Adds an observer to the list of observers. Observers would be
// types such as loggers, packet responders, and other "classes" that need to
// concern themselves with incoming packet events.
func (e *EventManager) AddObserver(observer Observer) {
	e.observers = append(e.observers, observer)
}

// RemoveObserver - Removes an observer from the list of observers. Generally
// one wouldn't want to do this in practice, but it's there if we need it.
func (e *EventManager) RemoveObserver(observer Observer) {
	observers := []Observer{}

	for _, o := range e.observers {
		if o != observer {
			observers = append(observers, o)
		}
	}

	e.observers = observers
}

// NotifyObservers - Notifies observers that a network event has occured.
func (e *EventManager) NotifyObservers(event Event) {
	for _, observer := range e.observers {
		observer.Notify(event)
	}
}

func (e *EventManager) Notify(event Event) {
	switch event.(type) {
	case NetworkEvent:
		Logger().Log("EventManager: network event received.")
	}
}
