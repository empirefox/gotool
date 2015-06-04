package ws

import (
	"github.com/golang/glog"
	"github.com/gorilla/websocket"
)

func chanFromWs(ws *websocket.Conn) chan []byte {
	defer func() {
		if err := recover(); err != nil {
			glog.Errorln(err)
		}
	}()
	c := make(chan []byte, 64)

	go func() {
		for {
			_, b, err := ws.ReadMessage()
			if len(b) > 0 {
				c <- b
			}
			if err != nil {
				c <- nil
				break
			}
		}
	}()

	return c
}

func Pipe(conn1 *websocket.Conn, conn2 *websocket.Conn) {
	chan1 := chanFromWs(conn1)
	chan2 := chanFromWs(conn2)

	for {
		select {
		case b1 := <-chan1:
			if b1 == nil {
				return
			} else {
				conn2.WriteMessage(websocket.TextMessage, b1)
			}
		case b2 := <-chan2:
			if b2 == nil {
				return
			} else {
				conn1.WriteMessage(websocket.TextMessage, b2)
			}
		}
	}
}
