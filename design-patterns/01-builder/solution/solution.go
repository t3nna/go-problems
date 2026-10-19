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
	config := ServerConfig{Host: "localhost", Port: 8080, Timeout: 30 * time.Second, MaxConnections: 100}
	for _, option := range opts {
		if option == nil {
			return ServerConfig{}, fmt.Errorf("server option is nil")
		}
		if err := option(&config); err != nil {
			return ServerConfig{}, err
		}
	}
	return config, nil
}

func WithHost(host string) Option {
	return func(c *ServerConfig) error {
		if host == "" {
			return fmt.Errorf("host must not be empty")
		}
		c.Host = host
		return nil
	}
}
func WithPort(port int) Option {
	return func(c *ServerConfig) error {
		if port < 1 || port > 65535 {
			return fmt.Errorf("port %d is outside 1..65535", port)
		}
		c.Port = port
		return nil
	}
}
func WithTimeout(timeout time.Duration) Option {
	return func(c *ServerConfig) error {
		if timeout <= 0 {
			return fmt.Errorf("timeout must be positive")
		}
		c.Timeout = timeout
		return nil
	}
}
func WithTLS(enabled bool) Option { return func(c *ServerConfig) error { c.TLS = enabled; return nil } }
func WithMaxConnections(limit int) Option {
	return func(c *ServerConfig) error {
		if limit <= 0 {
			return fmt.Errorf("max connections must be positive")
		}
		c.MaxConnections = limit
		return nil
	}
}
