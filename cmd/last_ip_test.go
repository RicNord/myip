package cmd

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/google/uuid"
	"github.com/spf13/viper"
)

func TestGetLastKnownIp(t *testing.T) {
	// Setup a temporary directory for test files
	tempDir := t.TempDir()

	// Helper function to write a file
	writeFile := func(filePath, content string) error {
		return os.WriteFile(filePath, []byte(content), 0644)
	}

	// Test case: file exists and contains a valid IP
	t.Run("File exists and contains IP", func(t *testing.T) {
		stateFile := uuid.NewString()
		filePath := filepath.Join(tempDir, stateFile)
		mockIP := "192.168.1.1"

		err := writeFile(filePath, mockIP)
		if err != nil {
			t.Fatalf("Failed to write mock IP to file: %v", err)
		}

		viper.Set("state-file", filePath)
		defer viper.Reset()

		ip, exist := getLastKnownIp()
		if !exist {
			t.Errorf("Expected IP to exist, but it doesn't")
		}
		if ip != mockIP {
			t.Errorf("Expected IP %s, but got %s", mockIP, ip)
		}
	})

	// Test case: file does not exist
	t.Run("File does not exist", func(t *testing.T) {
		nonFile := uuid.NewString()
		nonExistentFile := filepath.Join(tempDir, nonFile)
		viper.Set("state-file", nonExistentFile)
		defer viper.Reset()

		ip, exist := getLastKnownIp()
		if exist {
			t.Errorf("Expected file to not exist, but it does")
		}
		if ip != "" {
			t.Errorf("Expected empty IP, but got %s", ip)
		}
	})

	// Test case: file exists but is empty
	t.Run("File exists but is empty", func(t *testing.T) {
		emptyUuid := uuid.NewString()
		emptyFile := filepath.Join(tempDir, emptyUuid)
		err := writeFile(emptyFile, "")
		if err != nil {
			t.Fatalf("Failed to create empty file: %v", err)
		}

		viper.Set("state-file", emptyFile)
		defer viper.Reset()

		ip, exist := getLastKnownIp()
		if !exist {
			t.Errorf("Expected file to exist, but it does")
		}
		if ip != "" {
			t.Errorf("Expected empty IP, but got %s", ip)
		}
	})
}
