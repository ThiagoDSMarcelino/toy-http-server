package request

import (
	"fmt"
	"io"
	"toy-http-server/internal/headers"
)

type parseState int

const (
	STATE_INIT         = 0
	STATE_REQUEST_LINE = 1
	STATE_HEADERS      = 2
	STATE_BODY         = 3
	STATE_DONE         = 4
)

type Request struct {
	RequestLine RequestLine
	Headers     headers.Headers
	state       parseState
}

func newRequest() *Request {
	return &Request{
		state:   STATE_INIT,
		Headers: headers.NewHeaders(),
	}
}

func (r *Request) parse(data []byte) (int, error) {
	read := 0

outer:
	for {
		switch r.state {
		case STATE_INIT:
			rl, n, err := parseRequestLine(data)
			if err != nil {
				return 0, err
			}

			if n == 0 {
				break outer
			}

			r.RequestLine = *rl
			r.state = STATE_HEADERS
			read += n

		case STATE_HEADERS:
			n, done, err := r.Headers.Parse(data[read:])
			if err != nil {
				return 0, err
			}

			if done {
				r.state = STATE_DONE
			}

			if n == 0 {
				break outer
			}

			read += n

		case STATE_DONE:
			break outer

		default:
			return read, fmt.Errorf("invalid state")
		}
	}

	return read, nil
}

func (r *Request) done() bool {
	return r.state == STATE_DONE
}

const BUFFER_SIZE = 4096

func RequestFromReader(r io.Reader) (*Request, error) {
	req := newRequest()
	buffer := make([]byte, BUFFER_SIZE)
	len := 0

	for !req.done() {
		readN, err := r.Read(buffer[len:])
		if err != nil {
			return nil, err
		}

		len += readN

		precessedN, err := req.parse(buffer[:len])
		if err != nil {
			return nil, err
		}

		copy(buffer, buffer[precessedN:len])
		len -= precessedN
	}

	return req, nil
}
