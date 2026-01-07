package request

import (
	"fmt"
	"io"
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
	RequestLine *RequestLine
	state       parseState
}

func newRequest() *Request {
	return &Request{
		state: STATE_INIT,
	}
}

func (r *Request) parse(data []byte) (int, error) {
	read := 0

	switch r.state {
	case STATE_INIT:
		rl, n, err := parseRequestLine(data)
		if err != nil {
			return 0, err
		}

		if n == 0 {
			break
		}

		r.RequestLine = rl
		r.state = STATE_DONE
		read += n

	case STATE_DONE:
		break

	default:
		return read, fmt.Errorf("invalid state")
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
