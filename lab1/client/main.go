package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"net"
)

type Msg struct {
	Data string `json:"data"`
}

type Results struct {
	Data string `json:"data"`
}

func main() {

	conn, err := net.Dial("tcp", "185.102.139.168:8081")
	if err != nil {
		fmt.Println("Error connecting to server:", err)
		return
	}
	defer conn.Close()

	var a string
	fmt.Print("Enter your message: ")
	fmt.Scan(&a)

	data, err := json.Marshal(Msg{Data: a})
	if err != nil {
		fmt.Println("Error encoding JSON:", err)
		return
	}

	_, err = conn.Write(append(data, '\n'))
	if err != nil {
		fmt.Println("Error sending data:", err)
		return
	}

	response, err := bufio.NewReader(conn).ReadString('\n')
	if err != nil {
		fmt.Println("Error reading response:", err)
		return
	}

	var result Results
	err = json.Unmarshal([]byte(response), &result)
	if err != nil {
		fmt.Println("Error decoding JSON:", err)
		return
	}

	// Выводим результаты
	fmt.Printf("Message: %s\n", result.Data)
}
