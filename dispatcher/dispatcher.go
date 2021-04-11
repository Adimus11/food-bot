package dispatcher

import "sync"

type Dispatcher struct {
	clients    map[*Client]bool
	broadcast  chan []byte
	register   chan *Client
	unregister chan *Client

	mu sync.RWMutex
}

func NewDispatcher() *Dispatcher {
	dis := &Dispatcher{
		broadcast:  make(chan []byte),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		clients:    make(map[*Client]bool),
	}
	go dis.run()
	return dis
}

func (d *Dispatcher) run() {
	for {
		select {
		case client := <-d.register:
			d.mu.Lock()
			d.clients[client] = true
			d.mu.Unlock()
		case client := <-d.unregister:
			d.mu.Lock()
			if _, ok := d.clients[client]; ok {
				delete(d.clients, client)
				close(client.send)
			}
			d.mu.Unlock()
		case message := <-d.broadcast:
			d.mu.RLock()
			for client := range d.clients {
				select {
				case client.send <- message:
				default:
					close(client.send)
					delete(d.clients, client)
				}
			}
			d.mu.RUnlock()
		}
	}
}
