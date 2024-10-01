package main

import (
	"fmt"
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

type Message struct {
	Command string `json:"command"`
	Key     string `json:"key"`
	Value   string `json:"value"`
}

var store = make(map[string]string)
var mutex = &sync.Mutex{}

func handleWebSocket(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		fmt.Println("Error while upgrading:", err)
		return
	}
	defer conn.Close()

	for {
		var msg Message
		err := conn.ReadJSON(&msg)
		if err != nil {
			fmt.Println("Error reading JSON:", err)
			break
		}

		mutex.Lock()
		switch msg.Command {
		case "ADD":
			store[msg.Key] = msg.Value
			fmt.Printf("Added: %s = %s\n", msg.Key, msg.Value)
		}

		// Отправка обновленного map клиенту
		conn.WriteJSON(store)
		mutex.Unlock()
	}
}

func main() {
	http.HandleFunc("/ws", handleWebSocket)
	// Замените "192.168.1.100" на ваш IP-адрес
	fmt.Println("WebSocket server started at :185.102.139.168:8084")
	if err := http.ListenAndServe("185.102.139.168:8084", nil); err != nil {
		panic("Error starting server: " + err.Error())
	}
}
