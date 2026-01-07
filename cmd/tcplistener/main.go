package main

import (
	"fmt"
	"log/slog"
	"net"
	"toy-http-server/internal/request"
)

func handleConnection(conn net.Conn) {
	defer conn.Close()

	r, err := request.RequestFromReader(conn)
	if err != nil {
		slog.Error("failed to parse request", "error", err)
		return
	}

	fmt.Printf("Request line:\n")
	fmt.Printf(" - Method: %s\n", r.RequestLine.Method)
	fmt.Printf(" - Target: %s\n", r.RequestLine.RequestTarget)
	fmt.Printf(" - Version: %s\n", r.RequestLine.HttpVersion)
	fmt.Printf("Headers:\n")
	r.Headers.ForEach(func(name, value string) {
		fmt.Printf(" - %s: %s\n", name, value)
	})

	fmt.Fprintf(conn, "HTTP/%s 200 OK\r\n\r\n", r.RequestLine.HttpVersion)
}

const NETWORK = "tcp"
const ADDRESS = ":8080"

func main() {
	listener, err := net.Listen(NETWORK, ADDRESS)
	if err != nil {
		slog.Error("failed to listen", "error", err)
		return
	}

	fmt.Printf("Listening on %s %s\n", NETWORK, ADDRESS)

	for {
		conn, err := listener.Accept()
		if err != nil {
			slog.Error("failed to accept connection", "error", err)
			continue
		}

		go handleConnection(conn)
	}
}
