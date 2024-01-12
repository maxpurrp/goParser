package main

import (
	"fmt"
	"net"
	"time"
)

func main() {
	conn, err := net.Dial("tcp", ":8000")
	if err != nil {
		fmt.Println("Error connection:", err)
	}
	for i := 0; i < 5; i++ {
		message := "Hello Max"
		fmt.Println("Sending Message:", message)

		_, err = conn.Write([]byte(message))
		if err != nil {
			fmt.Println("Error sending message:", err)
		}
		buff := make([]byte, 1024)
		n, err := conn.Read(buff)
		if err != nil {
			fmt.Printf("error with reading anser: %v \n", err)
		}
		result := string(buff[:n])
		fmt.Printf("Got answer: %s \n", result)
		time.Sleep(1 * time.Second)
	}
	conn.Close()
}
