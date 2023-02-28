package socket

import (
	"websocket-service/pkg/errors"
	"websocket-service/pkg/logger"
)

type Request struct {
	Type    string `json:"type"`
	Content string `json:"data"`
	To      string `json:"to"`
}

type Subscription struct {
	// Clinet
	conn   *Client
	room   string
	closed bool
	err    errors.Error
}

// hub maintains the set of active connections and broadcasts messages to the
// connections.
type Hub struct {
	// Registered connections.
	rooms map[string]map[*Subscription]bool

	// Register requests from the connections.
	register chan *Subscription

	// Unregister requests from connections.
	unregister chan *Subscription

	// logger is used to log messages
	log logger.LoggerI

	// to find out who the new user is
	newRegister chan bool

	// to write the error to the log
	err errors.Error
}

func NewHub(log logger.LoggerI) *Hub {
	return &Hub{
		err:         *errors.NewError(log, "Hub", ""),
		register:    make(chan *Subscription),
		unregister:  make(chan *Subscription),
		newRegister: make(chan bool),
		rooms:       make(map[string]map[*Subscription]bool),
		log:         log,
	}
}

func (h *Hub) Run() {
	for {
		select {
		case s := <-h.register:
			connections := h.rooms[s.room]
			if connections == nil {
				connections = make(map[*Subscription]bool)
				h.rooms[s.room] = connections
			}
			h.rooms[s.room][s] = true
			h.newRegister <- true
		case s := <-h.unregister:

			connections := h.rooms[s.room]
			if connections != nil {
				if _, ok := connections[s]; ok {
					delete(connections, s)
					if len(connections) == 0 {
						delete(h.rooms, s.room)
					}
				}
			}
		}
	}
}

func (h *Hub) Send(msg Message, room string) error {
	for s := range h.rooms[room] {
		if s.closed {
			continue
		}

		err := s.conn.write(msg)
		if err != nil {
			h.err.Wrap(&err, "Error writing Send message", []interface{}{msg, room})
			continue
		}
	}

	return nil
}
