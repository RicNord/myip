package cmd

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/google/uuid"
	"github.com/spf13/viper"
)

func TestGetAliasOrIp(t *testing.T) {
	viper.Set("aliases", []Alias{
		{Alias: "home", IP: "127.0.0.1"},
	})
	defer viper.Reset()

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

func TestWriteIpToFile(t *testing.T) {
	// Setup a temporary file to simulate the state file
	tempDir := t.TempDir()
	stateFile := uuid.NewString()
	filePath := filepath.Join(tempDir, stateFile)

	viper.Set("state-file", filePath)
	defer viper.Reset()

	t.Run("Write IP to file successfully", func(t *testing.T) {
		ip := "192.168.1.1"
		err := writeIpToFile(ip)
		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}

		// Read the file to verify the content
		content, err := os.ReadFile(filePath)
		if err != nil {
			t.Fatalf("Failed to read file: %v", err)
		}
		if string(content) != ip {
			t.Errorf("Expected file content %s, got %s", ip, string(content))
		}
	})
}
