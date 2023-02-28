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
}

// hub maintains the set of active connections and broadcasts messages to the
// connections.
type Hub struct {
	// Registered connections.
	rooms map[string]map[*Client]bool

	// Register requests from the connections.
	register chan Subscription

	// Unregister requests from connections.
	unregister chan Subscription

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
		register:    make(chan Subscription),
		unregister:  make(chan Subscription),
		newRegister: make(chan bool),
		rooms:       make(map[string]map[*Client]bool),
		log:         log,
	}
}

func (h *Hub) Run() {
	for {
		select {
		case s := <-h.register:
			connections := h.rooms[s.room]
			if connections == nil {
				connections = make(map[*Client]bool)
				h.rooms[s.room] = connections
			}
			h.rooms[s.room][s.conn] = true
			h.newRegister <- true
		case s := <-h.unregister:

			connections := h.rooms[s.room]
			if connections != nil {
				if _, ok := connections[s.conn]; ok {
					delete(connections, s.conn)
					if len(connections) == 0 {
						delete(h.rooms, s.room)
					}
				}
			}
		}
	}
}

func (h *Hub) Read() {
	for {
		select {
		case <-h.newRegister:
			for room := range h.rooms {

				if len(h.rooms[room]) > 1 {
					continue
				}
				switch room {
				case "status":
					h.Status(room)
				case "ping":
					h.Ping(room)
				case "message":
					h.Message(room)
				default:
					h.NotFount(room)
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

		err := s.write(msg)
		if err != nil {
			h.err.Wrap(&err, "Error writing Send message", []interface{}{msg, room})
			continue
		}
	}

	return nil
}
