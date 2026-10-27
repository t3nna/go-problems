package main

import "fmt"

type Opts struct {
	ip      string
	domain  string
	maxConn int
	withTLS bool
}

type OptFunc func(*Opts)

func withTls(o *Opts) {
	o.withTLS = true
}

func setMaxConn(n int) OptFunc {
	return func(o *Opts) {
		o.maxConn = n
	}
}

type Server struct {
	Opts
}

func defaultOpts() Opts {
	return Opts{
		ip:      "0.0.0.0",
		domain:  "",
		maxConn: 10,
		withTLS: false,
	}
}

func newServer(opts ...OptFunc) *Server {
	o := defaultOpts()

	for _, fn := range opts {
		fn(&o)
	}

	return &Server{
		Opts: o,
	}
}

func main() {
	s := newServer(withTls, setMaxConn(154))
	fmt.Printf("%+v", s)

}
