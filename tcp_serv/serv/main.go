package main

import (
	"fmt"
	"io"
	"net"
)

func handleConnection(conn net.Conn) {
	defer conn.Close()
	fmt.Printf("Connection with %v succesful \n", conn.RemoteAddr().String())
	buff := make([]byte, 1024)
	for {
		n, err := conn.Read(buff)
		if err != nil {
			if err == io.EOF {
				fmt.Printf("Client %v disconnected \n", conn.RemoteAddr().String())
				return
			}
			fmt.Printf("Error with readind data: %v \n", err)
			return
		}
		msg := buff[:n]
		fmt.Printf("Got massage : %s \n", string(msg))
		answer := string(msg) + " from_server"
		_, err = conn.Write([]byte(answer))
		if err != nil {
			fmt.Printf("Failed with sendind response: %s \n", err)
		}
	}

}

func main() {
	listener, err := net.Listen("tcp", ":8000")
	if err != nil {
		fmt.Printf("Error listening to the adress: %v", err)
	}
	defer listener.Close()
	fmt.Println("Server starts succesfuly")
	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Printf("Error with handling connection: %v", err)
			continue
		}
		handleConnection(conn)
	}
}
