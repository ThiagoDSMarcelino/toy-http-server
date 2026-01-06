package request

import (
	"bytes"
	"fmt"
	"io"
)

type RequestLine struct {
	HttpVersion   string
	RequestTarget string
	Method        string
}

type Request struct {
	RequestLine *RequestLine
}

// Theoretically, \n could be the separator as well if the first line ends with it, but for now we only support \r\n.
var SEPARATOR = []byte("\r\n")

var HTTP_VERSION_PREFIX = []byte("HTTP/")
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

func parseRequestLine(buffer []byte) (requestLine *RequestLine, rest []byte, err error) {
	// SP = single space
	// Request line format: <METHOD> <SP> <REQUEST_TARGET> <SP> HTTP/<VERSION>

	before, after, found := bytes.Cut(buffer, SEPARATOR)
	if !found {
		return nil, after, MALFORMED_REQUEST_ERROR
	}

	parts := bytes.Split(before, []byte(" "))
	if len(parts) != 3 {
		return nil, after, MALFORMED_START_LINE_ERROR
	}

	method := parts[0]
	if !isValidMethod(method) {
		return nil, after, INVALID_METHOD_ERROR
	}

	version := bytes.TrimPrefix(parts[2], HTTP_VERSION_PREFIX)
	if !isValidVersion(version) {
		return nil, after, INVALID_VERSION_ERROR
	}

	rl := RequestLine{
		Method:        string(parts[0]),
		RequestTarget: string(parts[1]),
		HttpVersion:   string(version),
	}

	return &rl, buffer[len(before)+len(SEPARATOR):], nil
}

func RequestFromReader(r io.Reader) (*Request, error) {
	buffer, err := io.ReadAll(r)
	if err != nil {
		return nil, err
	}

	rl, _, err := parseRequestLine(buffer)
	if err != nil {
		return nil, err
	}

	res := &Request{
		RequestLine: rl,
	}

	return res, nil
}
