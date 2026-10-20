package request

import (
	"fmt"
	"io"
	"strings"
)

type RequestLine struct {
	HttpVersion   string
	RequestTarget string
	Method        string
}

type Request struct {
	RequestLine RequestLine
}

var ERROR_BAD_START_LINE = fmt.Errorf("bad request line")
var ERROR_BAD_HTTP_STATUS = fmt.Errorf("bad http status")
var SEPARATOR = "\r\n"

type parserState string

const (
	StateInit parserState = "init"
	StateDone parserState = "done"
)

func parseRequestLine(b string) (*RequestLine, string, error) {
	idx := strings.Index(b, SEPARATOR)

	if idx == -1 {
		return nil, "", ERROR_BAD_START_LINE
	}

	strartLine := b[:idx]
	restMsg := b[idx+len(SEPARATOR):]

	parts := strings.Split(strartLine, " ")

	if len(parts) != 3 {
		return nil, "", ERROR_BAD_START_LINE
	}

	httpsVersionParts := strings.Split(parts[2], "/")
	if len(httpsVersionParts) != 2 || httpsVersionParts[0] != "HTTP" || httpsVersionParts[1] != "1.1" {
		return nil, "", ERROR_BAD_HTTP_STATUS
	}

	rl := &RequestLine{
		Method:        parts[0],
		RequestTarget: parts[1],
		HttpVersion:   httpsVersionParts[1],
	}

	return rl,
		restMsg,
		nil
}

func RequestFromReader(reader io.Reader) (*Request, error) {
	buf, err := io.ReadAll(reader)
	if err != nil {
		return nil, err
	}
	line, _, err := parseRequestLine(string(buf))
	if err != nil {
		return nil, err
	}

	req := &Request{
		RequestLine: *line,
	}

	return req, nil
}
