package socket

import (
	"time"

	"github.com/gorilla/websocket"
)

func (s *Subscription) SendOne() {
	switch s.room {
	case "status":
		s.Status()
	default:
		s.NotFount()
	}
}

func (s *Subscription) Status() {
	var err error
	defer s.err.Wrap(&err, "h.Status", "")
	for {
		if s.closed {
			return
		}
		err = s.conn.write(Message{
			Code:   1,
			Status: "status",
			Data:   "next",
		})
		if err != nil {
			break
		}
		time.Sleep(4 * time.Second)
	}
	return
}

func (s *Subscription) NotFount() {
	var err error
	defer s.err.Wrap(&err, "h.Status", "")

	err = s.conn.write(Message{
		Code:   websocket.CloseInvalidFramePayloadData,
		Status: "not fount",
		Data:   "no such room available",
	})

	s.CloseOne()

	return
}
