package main

import (
	"fmt"
	"time"
)

type ServerConfig struct {
	Host           string
	Port           int
	Timeout        time.Duration
	TLS            bool
	MaxConnections int
}

type Option func(*ServerConfig) error

func NewServerConfig(opts ...Option) (ServerConfig, error) {
	conf := &ServerConfig{
		Host:           "localhost",
		Port:           8080,
		Timeout:        30 * time.Second,
		MaxConnections: 100,
	}

	for _, option := range opts {
		if option == nil {
			return ServerConfig{}, fmt.Errorf("option error")
		}
		err := option(conf)
		if err != nil {
			return ServerConfig{}, err
		}
	}

	return *conf, nil
}

func WithHost(host string) Option {
	return func(s *ServerConfig) error {
		if host == "" {
			return fmt.Errorf("host must be a string")
		}
		s.Host = host

		return nil
	}
}
func WithPort(port int) Option {
	return func(s *ServerConfig) error {
		if port < 1 || port > 65535 {
			return fmt.Errorf("error")
		}
		s.Port = port

		return nil
	}
}
func WithTimeout(timeout time.Duration) Option {
	return func(s *ServerConfig) error {
		if timeout <= 0 {
			return fmt.Errorf("error")
		}

		s.Timeout = timeout

		return nil
	}
}
func WithTLS(enabled bool) Option {
	return func(s *ServerConfig) error {
		s.TLS = enabled
		return nil
	}
}
func WithMaxConnections(limit int) Option {
	return func(s *ServerConfig) error {
		if limit <= 0 {
			return fmt.Errorf("error")
		}
		s.MaxConnections = limit
		return nil
	}
}
