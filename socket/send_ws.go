package socket

import (
	"fmt"
	"time"
	"websocket-service/pkg/errors"

	"github.com/gorilla/websocket"
)

// broadcast to everyone in the room
func (h *Hub) SendRoom() {
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

func (h *Hub) Status(room string) (err error) {
	defer errors.WrapCheck(&err, "h.Status")
	for {
		if len(h.rooms[room]) == 0 {
			return nil
		}
		err = h.Send(Message{
			Code:   1,
			Status: "status",
			Data:   "next",
		}, room)
		if err != nil {
			// h.log.Error("Error sending status message: ", log)
			break
		}
		time.Sleep(4 * time.Second)
	}
	return err
}

func (h *Hub) Ping(room string) (err error) {
	defer errors.WrapCheck(&err, "h.Ping")

	var count int
	for {
		fmt.Println(h.rooms[room])
		if len(h.rooms[room]) == 0 {
			return nil
		}
		err = h.Send(Message{
			Code:   websocket.PongMessage,
			Status: "pong",
			Data:   "pong",
		}, room)
		if err != nil {
			break
		}
		time.Sleep(4 * time.Second)
		count++
	}
	return err
}

func (h *Hub) Message(room string) (err error) {
	defer errors.WrapCheck(&err, "h.Message")

	// for {
	// 	if len(h.rooms[room]) == 0 {
	// 		return nil
	// 	}
	// 	err = h.Send(Message{
	// 		Code:   websocket.PongMessage,
	// 		Status: "Message",
	// 		Data:   "Message",
	// 	}, room)
	// 	if err != nil {
	// 		break
	// 	}
	// 	time.Sleep(1 * time.Second)
	// }
	return err
}

func (h *Hub) NotFount(room string) (err error) {
	defer errors.WrapCheck(&err, "h.NotFount")
	fmt.Println("NotFount")

	if len(h.rooms[room]) == 0 {
		return nil
	}
	err = h.Send(Message{
		Code:   websocket.CloseInvalidFramePayloadData,
		Status: "not fount",
		Data:   "no such room available",
	}, room)

	for s := range h.rooms[room] {
		s.CloseHub(h)
	}

	time.Sleep(1 * time.Second)
	return err
}
