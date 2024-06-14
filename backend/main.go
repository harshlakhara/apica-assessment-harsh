package main

import (
	"lru-cache/socket"
	"net/http"

	"golang.org/x/net/websocket"
)

func main() {
	// 	cache := cache.NewLRUCache(3, 5)
	// 	cache.Set("0", 0)
	// 	cache.Set("1", 1)
	// 	fmt.Println(cache.Get("0"))
	// 	time.Sleep(time.Second * 10)
	// 	fmt.Println(cache.Get("0"))
	server := socket.NewSocketServer()
	http.Handle("/ws", websocket.Handler(server.HandleWS))
	http.ListenAndServe(":3000", nil)
}
