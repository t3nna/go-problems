package main

import (
	"errors"
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
	return ServerConfig{}, errors.New("TODO: build server configuration")
}

func WithHost(host string) Option              { return func(*ServerConfig) error { return nil } }
func WithPort(port int) Option                 { return func(*ServerConfig) error { return nil } }
func WithTimeout(timeout time.Duration) Option { return func(*ServerConfig) error { return nil } }
func WithTLS(enabled bool) Option              { return func(*ServerConfig) error { return nil } }
func WithMaxConnections(limit int) Option      { return func(*ServerConfig) error { return nil } }
