package socket

import (
	"encoding/json"
	"log"

	"github.com/gofiber/contrib/websocket"
)

type messageType struct {
	Content string `json:"content"`
}

var clients = make(map[*websocket.Conn]*client) // Map to store connected clients

type client struct {
	conn *websocket.Conn
}

func WSHandler(c *websocket.Conn) {
	defer func() {
		delete(clients, c)
	}()

	clients[c] = &client{conn: c}

	for {
		_, message, err := c.ReadMessage()
		if err != nil {
			log.Println("read error:", err)
			break
		}

		var msg messageType
		if err := json.Unmarshal(message, &msg); err != nil {
			log.Println("unmarshal error:", err)
			continue
		}
	}
}

func Broadcast(msg interface{}) {
	for client := range clients {
		if err := client.Conn.WriteJSON(msg); err != nil {
			log.Println("write error:", err)
			delete(clients, client)
			client.Conn.Close()
		}
	}
}
