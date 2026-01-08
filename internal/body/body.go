package body

import "fmt"

func Parse(data []byte, n int) ([]byte, int, error) {
	fmt.Printf("Buffer: %s\n", data)

	return data, 0, nil
}
