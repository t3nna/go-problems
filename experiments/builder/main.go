package main

import "fmt"

type Opts struct {
	ip      string
	maxConn int
	tls     bool
	vpn     bool
}

type setOpts func(*Opts)

func newConn(o ...setOpts) *Opts {

	opts := &Opts{
		ip: "127.0.0.1",
	}
	for _, fn := range o {
		fn(opts)
	}

	return opts
}

func setMaxConn(n int) setOpts {
	return func(o *Opts) {
		o.maxConn = n
	}
}

func setTls(o *Opts) {
	o.tls = true
}

func main() {
	conn := newConn(setMaxConn(44))

	fmt.Println(conn)

}
