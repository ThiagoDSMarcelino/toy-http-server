package requestline

import (
	"bytes"
	"fmt"
	constants "toy-http-server/internal"
)

type RequestLine struct {
	HttpVersion   string
	RequestTarget string
	Method        string
}

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

// Parse reads and validates the HTTP request line from the given byte slice.
//
// It expects the standard format:
//
//	<METHOD> <SP> <REQUEST_TARGET> <SP> HTTP/<VERSION>
//
// It returns the parsed RequestLine, the number of bytes read, and an error if the format is invalid.
func Parse(data []byte) (requestLine *RequestLine, read int, err error) {
	idx := bytes.Index(data, constants.LINE_SEPARATOR)
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

	return &rl, idx + len(constants.LINE_SEPARATOR), nil
}
