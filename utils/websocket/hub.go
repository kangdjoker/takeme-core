package websocket

type Message struct {
	Data []byte
	Room string
}

type subscription struct {
	conn *connection
	room string
}

// hub maintains the set of active connections and broadcasts messages to the
// connections.
type hub struct {
	// Registered connections.
	Room map[string]map[*connection]bool

	// Inbound messages from the connections.
	Broadcast chan Message

	// Register requests from the connections.
	Register chan subscription

	// Unregister requests from connections.
	Unregister chan subscription
}

var Hub = hub{
	Broadcast:  make(chan Message),
	Register:   make(chan subscription),
	Unregister: make(chan subscription),
	Room:       make(map[string]map[*connection]bool),
}

func (h *hub) Run() {
	for {
		select {
		case s := <-h.Register:
			connections := h.Room[s.room]
			if connections == nil {
				connections = make(map[*connection]bool)
				h.Room[s.room] = connections
			}
			h.Room[s.room][s.conn] = true
		case s := <-h.Unregister:
			connections := h.Room[s.room]
			if connections != nil {
				if _, ok := connections[s.conn]; ok {
					delete(connections, s.conn)
					close(s.conn.send)
					if len(connections) == 0 {
						delete(h.Room, s.room)
					}
				}
			}
		case m := <-h.Broadcast:
			connections := h.Room[m.Room]
			for c := range connections {
				select {
				case c.send <- m.Data:
				default:
					close(c.send)
					delete(connections, c)
					if len(connections) == 0 {
						delete(h.Room, m.Room)
					}
				}
			}
		}
	}
}
