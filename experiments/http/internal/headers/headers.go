package headers

import (
	"bytes"
	"fmt"
)

type Headers map[string]string

var rn = []byte("\r\n")

func NewHeaders() Headers {
	return make(Headers)
}

func parseHeader(fieldLine []byte) (string, string, error) {
	parts := bytes.SplitN(fieldLine, []byte(":"), 2)
	if len(parts) != 2 {
		return "", "", fmt.Errorf("malformed field line")
	}

	name := parts[0]
	value := bytes.TrimSpace(parts[1])

	fmt.Printf("'%s' - %v\n", string(name), bytes.HasSuffix(name, []byte(" ")))
	if bytes.HasPrefix(name, []byte(" ")) {
		return "", "", fmt.Errorf("malformed header name ")
	}
	return string(name), string(value), nil
}

func (h Headers) Parse(data []byte) (int, bool, error) {

	read := 0
	done := false
	for {
		idx := bytes.Index(data[read:], rn)
		if idx == -1 {
			break
		}

		if idx == 0 {
			done = true
			read += len(rn)
			break
		}
		fmt.Println(idx, read)
		fmt.Println(string(data[read : idx+read]))
		name, value, err := parseHeader(data[read : idx+read])
		if err != nil {
			return 0, false, err
		}
		read += idx + len(rn)

		h[name] = value

	}

	return read, done, nil

}
