package main

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"os"
)

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

func main() {

	v, err := os.Open("./messages.txt")
	if err != nil {
		log.Fatal(err)
	}

	for v := range getLinesChannel(v) {
		fmt.Printf("read: %s\n", v)
	}

}
