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
	RequestLine RequestLine
}

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

func RequestFromReader(r io.Reader) (*Request, error) {
	text, err := io.ReadAll(r)
	if err != nil {
		return nil, err
	}

	before, _, found := bytes.Cut(text, SEPARATOR)
	if !found {
		return nil, MALFORMED_REQUEST_ERROR
	}

	parts := bytes.Split(before, []byte(" "))
	if len(parts) != 3 {
		return nil, MALFORMED_START_LINE_ERROR
	}

	method := parts[0]
	if !isValidMethod(method) {
		return nil, INVALID_METHOD_ERROR
	}

	version := bytes.TrimPrefix(parts[2], HTTP_VERSION_PREFIX)
	if !isValidVersion(version) {
		return nil, INVALID_VERSION_ERROR
	}

	rl := RequestLine{
		Method:        string(parts[0]),
		RequestTarget: string(parts[1]),
		HttpVersion:   string(version),
	}

	res := &Request{
		RequestLine: rl,
	}

	return res, nil
}
