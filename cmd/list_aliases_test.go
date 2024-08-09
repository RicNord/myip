package cmd

import (
	"testing"

	"github.com/spf13/viper"
)

func TestGetAliases(t *testing.T) {
	viper.Set("aliases", []Alias{
		{Alias: "home", IP: "127.0.0.1"},
		{Alias: "office", IP: "192.168.0.1"},
	})

	aliases, err := getAliases()
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if len(aliases) != 2 {
		t.Fatalf("expected 2 aliases, got %d", len(aliases))
	}

	if aliases[0].Alias != "home" || aliases[0].IP != "127.0.0.1" {
		t.Fatalf("expected alias 'home' with IP '127.0.0.1', got alias '%s' with IP '%s'", aliases[0].Alias, aliases[0].IP)
	}

	if aliases[1].Alias != "office" || aliases[1].IP != "192.168.0.1" {
		t.Fatalf("expected alias 'office' with IP '192.168.0.1', got alias '%s' with IP '%s'", aliases[1].Alias, aliases[1].IP)
	}
}
