package headers

import (
	"bytes"
	"fmt"
	"regexp"
	"strconv"
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

func (h Headers) Contains(key string) bool {
	key = strings.ToLower(key)
	_, exists := h.data[key]
	return exists
}

func (h Headers) Get(key string) (string, bool) {
	key = strings.ToLower(key)
	value, exists := h.data[key]
	return value, exists
}

func (h Headers) GetInt(key string, defaultValue int) (int, bool) {
	valueStr, exists := h.Get(key)
	if !exists {
		return defaultValue, false
	}

	valueInt, err := strconv.Atoi(valueStr)
	if err != nil {
		return defaultValue, false
	}

	return valueInt, true
}

func (h Headers) Set(key, value string) {
	key = strings.ToLower(key)
	previous, exists := h.data[key]
	if exists {
		value = fmt.Sprintf("%s, %s", previous, value)
	}

	h.data[key] = value
}

func (h Headers) ForEach(f func(name string, value string)) {
	for name, value := range h.data {
		f(name, value)
	}
}

var MALFORMED_HEADER_ERROR = fmt.Errorf("malformed header")

var HEADER_SEPARATOR = []byte(":")
var SPACE = []byte(" ")

var VALID_KEY_CHARS_REGEX = regexp.MustCompile("^[a-zA-Z0-9!#$%&'*+.^_`|~-]+$")

func isValidKey(key []byte) bool {
	return !bytes.HasPrefix(key, SPACE) && !bytes.HasSuffix(key, SPACE) && len(key) != 0 && VALID_KEY_CHARS_REGEX.Match(key)
}

func (h Headers) Parse(data *[]byte) (n int, done bool, err error) {
	read := 0

	for {
		idx := bytes.Index((*data)[read:], constants.LINE_SEPARATOR)
		if idx == -1 {
			return read, false, nil
		}

		field := (*data)[read : read+idx]
		read += idx + len(constants.LINE_SEPARATOR)

		if len(field) == 0 {
			return read, true, nil
		}

		parts := bytes.SplitN(field, HEADER_SEPARATOR, 2)
		if len(parts) != 2 {
			return 0, false, MALFORMED_HEADER_ERROR
		}

		if !isValidKey(parts[0]) {
			return 0, false, MALFORMED_HEADER_ERROR
		}

		value := bytes.TrimSpace(parts[1])
		h.Set(string(parts[0]), string(value))
	}
}
