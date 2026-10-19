package main

import (
	"testing"
	"time"
)

func TestNewServerConfigUsesDefaultsAndOverrides(t *testing.T) {
	config, err := NewServerConfig(WithHost("api.example.com"), WithPort(8443), WithTLS(true), WithTimeout(5*time.Second), WithMaxConnections(50))
	if err != nil {
		t.Fatalf("NewServerConfig() error = %v", err)
	}
	if config.Host != "api.example.com" || config.Port != 8443 || !config.TLS || config.Timeout != 5*time.Second || config.MaxConnections != 50 {
		t.Fatalf("unexpected config: %+v", config)
	}
}
func TestNewServerConfigHasUsefulDefaults(t *testing.T) {
	config, err := NewServerConfig()
	if err != nil {
		t.Fatal(err)
	}
	if config.Host != "localhost" || config.Port != 8080 || config.Timeout != 30*time.Second || config.MaxConnections != 100 {
		t.Fatalf("unexpected defaults: %+v", config)
	}
}
func TestNewServerConfigRejectsInvalidOptions(t *testing.T) {
	for _, option := range []Option{WithHost(""), WithPort(0), WithPort(65536), WithTimeout(0), WithMaxConnections(0)} {
		if _, err := NewServerConfig(option); err == nil {
			t.Fatal("expected invalid option error")
		}
	}
}
