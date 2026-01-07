package request

import (
	"bytes"
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

type RequestLine struct {
	HttpVersion   string
	RequestTarget string
	Method        string
}

type Request struct {
	RequestLine *RequestLine
	state       parseState
}

func newRequest() *Request {
	return &Request{
		state: STATE_INIT,
	}
}

// Theoretically, \n could be the separator as well if the first line ends with it, but for now we only support \r\n.
var SEPARATOR = []byte("\r\n")

var VALID_METHODS = [][]byte{
	[]byte("GET"),
	[]byte("POST"),
	[]byte("PUT"),
	[]byte("DELETE"),
	[]byte("HEAD"),
	[]byte("OPTIONS"),
	[]byte("PATCH"),
	[]byte("TRACE"),
}

var HTTP_VERSION_PREFIX = []byte("HTTP/")
var VALID_VERSIONS = [][]byte{
	[]byte("1.1"),
}

var MALFORMED_REQUEST_ERROR = fmt.Errorf("malformed request")
var MALFORMED_START_LINE_ERROR = fmt.Errorf("malformed start line")
var INVALID_METHOD_ERROR = fmt.Errorf("invalid method")
var INVALID_VERSION_ERROR = fmt.Errorf("invalid version")

func isValidMethod(method []byte) bool {
	for _, m := range VALID_METHODS {
		if bytes.Equal(method, m) {
			return true
		}
	}

	return false
}

func isValidVersion(version []byte) bool {
	for _, v := range VALID_VERSIONS {
		if bytes.Equal(version, v) {
			return true
		}
	}

	return false
}

func parseRequestLine(data []byte) (requestLine *RequestLine, read int, err error) {
	// SP = single space
	// Request line format: <METHOD> <SP> <REQUEST_TARGET> <SP> HTTP/<VERSION>
	idx := bytes.Index(data, SEPARATOR)
	if idx == -1 {
		return nil, 0, nil
	}

	parts := bytes.Split(data[:idx], []byte(" "))
	if len(parts) != 3 {
		return nil, 0, MALFORMED_START_LINE_ERROR
	}

	method := parts[0]
	if !isValidMethod(method) {
		return nil, 0, INVALID_METHOD_ERROR
	}

	version := bytes.TrimPrefix(parts[2], HTTP_VERSION_PREFIX)
	if !isValidVersion(version) {
		return nil, 0, INVALID_VERSION_ERROR
	}

	rl := RequestLine{
		Method:        string(parts[0]),
		RequestTarget: string(parts[1]),
		HttpVersion:   string(version),
	}

	return &rl, idx + len(SEPARATOR), nil
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
	res := newRequest()
	buffer := make([]byte, BUFFER_SIZE)
	len := 0

	for !res.done() {
		readN, err := r.Read(buffer[len:])
		if err != nil {
			return nil, err
		}

		len += readN

		precessedN, err := res.parse(buffer[:len])
		if err != nil {
			return nil, err
		}

		copy(buffer, buffer[precessedN:len])
		len -= precessedN
	}

	return res, nil
}
