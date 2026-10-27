package main

import (
	"bytes"
	"fmt"
	"http/primagen/internal/request"
	"io"
	"log"
	"net"
)

func main() {

	listener, err := net.Listen("tcp", ":42069")
	if err != nil {
		log.Fatalf("error", "error", err)
	}

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Fatalf("error", "error", err)
		}

		req, err := request.RequestFromReader(conn)
		if err != nil {
			log.Fatalf("error", "error", err)
		}

		fmt.Printf("Request line:\n- Method: %s\n- Target: %s\n- Version: %s\n", req.RequestLine.Method, req.RequestLine.RequestTarget, req.RequestLine.HttpVersion)
	}

}

func getLinesChannel(f io.ReadCloser) <-chan string {
	out := make(chan string, 1)

	go func() {
		defer close(out)
		defer f.Close()

		str := ""
		for {
			data := make([]byte, 8)
			n, err := f.Read(data)
			if err != nil {
				break
			}
			data = data[:n]
			if i := bytes.IndexByte(data, '\n'); i != -1 {
				str += string(data[:i])
				data = data[i+1:]
				// fmt.Printf("read: %s\n", str)
				out <- str
				str = ""
			}
			str += string(data)

		}

	}()

	return out
}
