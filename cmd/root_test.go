package cmd

import (
	"testing"

	"github.com/spf13/viper"
)

func TestGetAliasOrIp(t *testing.T) {
	viper.Set("aliases", []Alias{
		{Alias: "home", IP: "127.0.0.1"},
	})

	ip, err := getAliasOrIp("127.0.0.1")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if ip != "home" {
		t.Fatalf("expected 'home', got %s", ip)
	}

	ip, err = getAliasOrIp("192.168.0.1")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if ip != "192.168.0.1" {
		t.Fatalf("expected '192.168.0.1', got %s", ip)
	}
}
