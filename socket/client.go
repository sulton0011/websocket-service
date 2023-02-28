package socket

import (
	"context"
	"net/http"
	"time"

	"websocket-service/config"
	"websocket-service/pkg/helper"
	"websocket-service/pkg/logger"
	"websocket-service/pkg/security"

	"github.com/gorilla/websocket"
)

type Message struct {
	Code   int    `json:"code"`
	Status string `json:"status"`
	Data   any    `json:"data"`
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

// connection is an middleman between the websocket connection and the hub.
type Client struct {
	// The websocket connection.
	ws     *websocket.Conn
	log    logger.LoggerI
	ctx    context.Context
	user   *security.TokenInfo
	closed bool
}

func NewSubscription(ws *websocket.Conn, log logger.LoggerI, ctx context.Context, user *security.TokenInfo, room string) Subscription {
	return Subscription{
		conn: &Client{
			ws:     ws,
			log:    log,
			ctx:    ctx,
			user:   user,
			closed: false,
		},
		room:   room,
		closed: false,
	}
}

// readPump pumps messages from the websocket connection to the hub.
func (s *Subscription) readPump(h *Hub) {
	c := s.conn
	defer func() {
		s.Close(h)
	}()
	c.ws.SetReadLimit(config.MaxMessageSize)

	// SetReadDeadline sets the read deadline on the underlying network connection. After a read has timed out,
	// the websocket connection state is corrupt and all future reads will return an error.
	// A zero value for t means reads will not time out.
	// c.ws.SetReadDeadline(time.Now().Add(config.PongWait))

	c.ws.SetPongHandler(func(string) error { c.ws.SetReadDeadline(time.Now().Add(config.PongWait)); return nil })
	for {
		var msg struct {
			Status string `json:"status"`
		}

		err := c.ws.ReadJSON(&msg)

		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway) {
				c.log.Error("!!!Errof", logger.Error(err))
			}
			break
		}
	}
}

// write writes a message with the given message type and payload.
func (c *Client) write(payload any) error {
	c.ws.SetWriteDeadline(time.Now().Add(config.WriteWait))
	return c.ws.WriteJSON(payload)
}

func (s *Subscription) Close(h *Hub) {
	s.closed = true
	s.conn.closed = true
	s.conn.ws.Close()
	delete(h.rooms[s.room], s.conn)
	h.unregister <- *s
}

// serveWs handles websocket requests from the peer.
func (h *Hub) ServeWs(w http.ResponseWriter, r *http.Request, ctx context.Context) {
	req := helper.GetValueContext(ctx)
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		h.err.Wrap(&err, "ServeWs", req)
		return
	}

	s := NewSubscription(ws, h.log, ctx, req, ctx.Value("room").(string))

	h.register <- s
	go s.readPump(h)
}
