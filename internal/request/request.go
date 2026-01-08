package request

import (
	"fmt"
	"io"
	"toy-http-server/internal/body"
	"toy-http-server/internal/headers"
	"toy-http-server/internal/requestline"
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
	RequestLine requestline.RequestLine
	Headers     headers.Headers
	Body        []byte
	state       parseState
}

func newRequest() *Request {
	return &Request{
		state:   STATE_INIT,
		Headers: headers.NewHeaders(),
	}
}

const CONTENT_LENGTH_HEADER = "Content-Length"

var INVALID_CONTENT_LENGTH_ERROR = fmt.Errorf("invalid Content-Length value")

func (r *Request) parse(data []byte) (int, error) {
	read := 0

outer:
	for {
		data = data[read:]

		switch r.state {
		case STATE_INIT:
			rl, n, err := requestline.Parse(data)
			if err != nil {
				return read, err
			}

			if n == 0 {
				break outer
			}

			r.RequestLine = *rl
			r.state = STATE_HEADERS
			read += n

		case STATE_HEADERS:
			n, done, err := r.Headers.Parse(data)
			if err != nil {
				return read, err
			}

			read += n

			if done {
				// Technically, there are cases where a body can be present without
				// a Content-Length header (e.g., Transfer-Encoding: chunked), but
				// for simplicity, we only check for Content-Length here.
				if r.Headers.Contains(CONTENT_LENGTH_HEADER) {
					r.state = STATE_BODY
				} else {
					r.state = STATE_DONE
				}

				break outer
			}

			if n == 0 {
				break outer
			}

		case STATE_BODY:
			cl, ok := r.Headers.GetInt(CONTENT_LENGTH_HEADER, 0)
			if !ok {
				return read, INVALID_CONTENT_LENGTH_ERROR
			}

			_, _, err := body.Parse(data, cl)
			if err != nil {
				return read, err
			}

			r.state = STATE_DONE

		case STATE_DONE:
			break outer

		default:
			panic("invalid state")
		}
	}

	return read, nil
}

func (r *Request) done() bool {
	return r.state == STATE_DONE
}

const BUFFER_SIZE = 4096

func FromReader(r io.Reader) (*Request, error) {
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
