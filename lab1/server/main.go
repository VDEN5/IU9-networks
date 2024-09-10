package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"net"
	"strings"
)

type Msg struct {
	Data string `json:"data"`
}

type Results struct {
	Data string `json:"data"`
}

func main() {

	fmt.Println("Launching server...")

	ln, err := net.Listen("tcp", "185.102.139.168:8081")
	if err != nil {
		fmt.Println("Error setting up listener:", err)
		return
	}

	conn, err := ln.Accept()
	if err != nil {
		fmt.Println("Error accepting connection:", err)
		return
	}

	for {
		message, err := bufio.NewReader(conn).ReadString('\n')
		if err != nil {
			fmt.Println("Error reading message:", err)
			return
		}

		fmt.Println("Message Received:", string(message))

		message = strings.TrimSpace(message)

		var res Msg
		err = json.Unmarshal([]byte(message), &res)
		if err != nil {
			fmt.Println("Error decoding JSON:", err)
			continue
		}

		fmt.Printf("To: %s\n", res.Data)

		result := Results{
			Data: res.Data,
		}

		response, err := json.Marshal(result)
		if err != nil {
			fmt.Println("Error encoding JSON:", err)
			continue
		}

		conn.Write([]byte(string(response) + "\n"))
	}
}
