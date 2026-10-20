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

func parseRequestLine(buf []byte, req *Request) (int, error) {
	idx := bytes.Index(buf, []byte{'\r', '\n'})
	if idx == -1 {
		return 0, fmt.Errorf("error")
	}
	line := buf[:idx]

	reqLineArr := bytes.Split(line, []byte{' '})
	if len(reqLineArr) != 3 {
		return 0, fmt.Errorf("error")
	}

	req.RequestLine.Method = string(reqLineArr[0])
	req.RequestLine.RequestTarget = string(reqLineArr[1])
	req.RequestLine.HttpVersion = string(bytes.Split(reqLineArr[2], []byte{'/'})[1])
	return idx, nil

}

func RequestFromReader(reader io.Reader) (*Request, error) {
	buf, err := io.ReadAll(reader)
	if err != nil {
		return nil, err
	}
	req := &Request{
		RequestLine: RequestLine{},
	}

	_, err = parseRequestLine(buf, req)
	if err != nil {
		return nil, err
	}

	return req, nil
}
