package request

import (
	"bytes"
	"fmt"
	"io"
)

type parserState string

const (
	StateInit  parserState = "init"
	StateDone  parserState = "done"
	StateError parserState = "error"
)

type RequestLine struct {
	HttpVersion   string
	RequestTarget string
	Method        string
}

type Request struct {
	RequestLine RequestLine
	state       parserState
}

func newRequest() *Request {
	return &Request{
		state: StateInit,
	}
}

var ERROR_BAD_START_LINE = fmt.Errorf("bad request line")
var ERROR_BAD_HTTP_STATUS = fmt.Errorf("bad http status")
var ERROR_REQUEST_IN_ERROR_STATTE = fmt.Errorf("request in error state")
var SEPARATOR = []byte("\r\n")

func parseRequestLine(b []byte) (*RequestLine, int, error) {
	idx := bytes.Index(b, SEPARATOR)

	if idx == -1 {
		return nil, 0, nil
	}

	strartLine := b[:idx]
	read := idx + len(SEPARATOR)

	parts := bytes.Split(strartLine, []byte(" "))

	if len(parts) != 3 {
		return nil, 0, ERROR_BAD_START_LINE
	}

	httpsVersionParts := bytes.Split(parts[2], []byte("/"))
	if len(httpsVersionParts) != 2 || string(httpsVersionParts[0]) != "HTTP" || string(httpsVersionParts[1]) != "1.1" {
		return nil, 0, ERROR_BAD_HTTP_STATUS
	}

	rl := &RequestLine{
		Method:        string(parts[0]),
		RequestTarget: string(parts[1]),
		HttpVersion:   string(httpsVersionParts[1]),
	}

	return rl,
		read,
		nil
}

func (r *Request) parse(data []byte) (int, error) {

	read := 0
outer:
	for {
		switch r.state {
		case StateError:
			return 0, ERROR_REQUEST_IN_ERROR_STATTE

		case StateInit:
			rl, n, err := parseRequestLine(data[read:])
			if err != nil {
				r.state = StateError
				return 0, err
			}
			if n == 0 {
				break outer
			}

			r.RequestLine = *rl
			read += n
			r.state = StateDone

		case StateDone:
			break outer
		}
	}
	return read, nil
}

func (r *Request) isDone() bool {
	return r.state == StateDone || r.state == StateError
}

func RequestFromReader(reader io.Reader) (*Request, error) {
	request := newRequest()

	// NOTE: buffer could be overrun
	buf := make([]byte, 1024)
	bufLen := 0

	for !request.isDone() {
		n, err := reader.Read(buf[bufLen:])
		if err != nil {
			return nil, err
		}

		bufLen += n

		readN, err := request.parse(buf[:bufLen+n])

		if err != nil {
			return nil, err
		}

		copy(buf, buf[readN:bufLen])
		bufLen -= readN

	}

	return request, nil
}
