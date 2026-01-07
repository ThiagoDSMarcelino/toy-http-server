package headers

import (
	"bytes"
	"fmt"
	"regexp"
	"strings"
	constants "toy-http-server/internal"
)

type Headers struct {
	data map[string]string
}

func NewHeaders() Headers {
	return Headers{
		data: make(map[string]string),
	}
}

func (h Headers) Get(key string) (string, bool) {
	key = strings.ToLower(key)
	value, exists := h.data[key]
	return value, exists
}

func (h Headers) Set(key, value string) {
	key = strings.ToLower(key)
	previous, exists := h.data[key]
	if exists {
		value = fmt.Sprintf("%s, %s", previous, value)
	}

	h.data[key] = value
}

var MALFORMED_HEADER_ERROR = fmt.Errorf("malformed header")

var HEADER_SEPARATOR = []byte(":")
var SPACE = []byte(" ")

var VALID_KEY_CHARS_REGEX = regexp.MustCompile("^[a-zA-Z0-9!#$%&'*+.^_`|~-]+$")

func isValidKey(key []byte) bool {
	return !bytes.HasPrefix(key, SPACE) && !bytes.HasSuffix(key, SPACE) && len(key) != 0 && VALID_KEY_CHARS_REGEX.Match(key)
}

func (h Headers) Parse(data []byte) (n int, done bool, err error) {
	idx := bytes.Index(data, constants.LINE_SEPARATOR)
	if idx == -1 {
		return 0, false, nil
	}

	line := data[:idx]
	read := idx + len(constants.LINE_SEPARATOR)

	if len(line) == 0 {
		return read, true, nil
	}

	parts := bytes.SplitN(line, HEADER_SEPARATOR, 2)
	if len(parts) != 2 {
		return 0, false, MALFORMED_HEADER_ERROR
	}

	if !isValidKey(parts[0]) {
		return 0, false, MALFORMED_HEADER_ERROR
	}

	value := bytes.TrimSpace(parts[1])
	h.Set(string(parts[0]), string(value))

	return read, false, nil
}
