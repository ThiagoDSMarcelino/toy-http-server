package body

func Parse(body, data *[]byte, contentLength int) (read int, done bool, err error) {
	done = false

	remaining := contentLength - len(*body)
	if remaining <= 0 {
		return 0, true, nil
	}

	read = min(remaining, len(*data))

	*body = append(*body, (*data)[:read]...)

	if len(*body) >= contentLength {
		done = true
	}

	return read, done, nil
}
