package request

import (
	"bytes"
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
