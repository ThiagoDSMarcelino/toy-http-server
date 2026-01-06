package main

import (
	"bytes"
	"fmt"
	"io"
	"log/slog"
	"net"
	"strings"
)

const BUFFER_SIZE = 8
const LINE_SEPARATOR = '\n'
const NETWORK = "tcp"
const ADDRESS = ":8080"

func readLinesFromChannel(f io.ReadCloser) <-chan string {
	out := make(chan string, 1)

	go func() {
		defer close(out)
		defer f.Close()

		line := strings.Builder{}
		data := make([]byte, BUFFER_SIZE)
		for {
			n, err := f.Read(data)
			if err != nil {
				break
			}

			chunk := data[:n]
			idx := bytes.IndexByte(chunk, LINE_SEPARATOR)

			// no line separator found
			if idx == -1 {
				line.Write(chunk)
				continue
			}

			line.Write(chunk[:idx])
			out <- line.String()
			line.Reset()

			for {
				chunk = chunk[idx+1:]
				idx = bytes.IndexByte(chunk, LINE_SEPARATOR)

				if idx == -1 {
					line.Write(chunk)
					break
				}

				line.Write(chunk[:idx])
				out <- line.String()
				line.Reset()
			}
		}

		if line.Len() > 0 {
			out <- line.String()
		}
	}()

	return out
}

func handleConnection(conn net.Conn) {
	defer conn.Close()

	fmt.Printf("Accepted connection from %s\n", conn.RemoteAddr().String())

	for line := range readLinesFromChannel(conn) {
		fmt.Printf("Received line: %s\n", line)
	}

	fmt.Printf("Connection closed from %s\n", conn.RemoteAddr().String())
}

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
