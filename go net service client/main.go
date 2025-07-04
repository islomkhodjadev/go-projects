package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"strings"
)

func main() {
	ln, err := net.Listen("tcp", ":8080")
	if err != nil {
		log.Fatal(err)
	}
	defer ln.Close()

	for {
		conn, err := ln.Accept()

		if err != nil {
			log.Fatal(err)
		}

		go handleConnection(conn)
	}

}

func handleConnection(conn net.Conn) {
	defer conn.Close()

	reader := bufio.NewReader(conn)
	line, err := reader.ReadString('\n')
	if err != nil {
		fmt.Fprintf(conn, "Error in reading command %v\n", err)
	}

	parts := strings.SplitN(strings.TrimSpace(line), " ", 2)
	if len(parts) != 2 {
		fmt.Fprintf(conn, "Invalid command format. Expected format: COMMAND RESOURCE \n")
		return
	}
	command := parts[0]
	resourse := parts[1]

	log.Printf("Recieved command %s the resourse %s", command, resourse)
	switch command {
	case "GET":
		handleGet(conn, resourse)
	default:
		fmt.Fprintf(conn, "Invalid command %s\n", command)
	}
}

func handleGet(conn net.Conn, resourse string) {
	fmt.Fprintf(conn, "GET command for resource %s\n", resourse)
}
