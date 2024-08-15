package client_notifier

import "sync"

type ClientNotifier interface {
	RegisterClient(id string) <-chan Event
	UnregisterClient(id string)
	NotifyClient(id string, e Event)
}

type client struct {
	id        string
	eventChan chan Event
}

func NewBasicNotifier() ClientNotifier {
	return &basicNotifier{
		m:       &sync.RWMutex{},
		clients: map[string]client{},
	}
}

type basicNotifier struct {
	m       *sync.RWMutex
	clients map[string]client
}

func (p *basicNotifier) RegisterClient(id string) <-chan Event {
	p.m.Lock()
	defer p.m.Unlock()

	eventChan := make(chan Event)
	p.clients[id] = client{
		id:        id,
		eventChan: eventChan,
	}

	return eventChan
}

func (p *basicNotifier) UnregisterClient(id string) {
	p.m.Lock()
	defer p.m.Unlock()

	delete(p.clients, id)
}

func (p *basicNotifier) NotifyClient(id string, e Event) {
	p.m.RLock()
	defer p.m.RUnlock()

	client, ok := p.clients[id]
	if !ok {
		return
	}
	client.eventChan <- e
}
