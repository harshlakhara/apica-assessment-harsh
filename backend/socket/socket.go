package socket

import (
	"fmt"
	"io"

	"golang.org/x/net/websocket"
)

type WebSocketServer struct {
	conns map[*websocket.Conn]bool
}

func NewSocketServer() *WebSocketServer {
	return &WebSocketServer{
		conns: make(map[*websocket.Conn]bool),
	}
}

func (s *WebSocketServer) HandleWS(ws *websocket.Conn) {
	fmt.Println("Incoming connection from: ", ws.RemoteAddr())

	s.conns[ws] = true

	s.readLoop(ws)
}

func (s *WebSocketServer) readLoop(ws *websocket.Conn) {
	buff := make([]byte, 1024)

	for {
		n, err := ws.Read(buff)
		if err != nil {
			if err == io.EOF {
				break
			}
			fmt.Println("Error reading message: ", err)
			continue
		}
		msg := buff[:n]

		fmt.Println(msg, " <--- received from ", ws.RemoteAddr())
	}
}

func (s *WebSocketServer) Broadcast(b []byte) {
	for conn := range s.conns {
		go conn.Write(b)
	}
}
