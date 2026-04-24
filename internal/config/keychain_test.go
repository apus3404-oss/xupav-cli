// internal/config/keychain_test.go
package config

import (
	"testing"
)

func TestKeychainSetGet(t *testing.T) {
	service := "mycli-test"
	account := "openrouter"
	secret := "test-api-key-12345"

	// Set key
	err := SetAPIKey(service, account, secret)
	if err != nil {
		t.Skipf("keychain not available: %v", err)
	}

	// Get key
	retrieved, err := GetAPIKey(service, account)
	if err != nil {
		t.Fatalf("failed to get API key: %v", err)
	}

	if retrieved != secret {
		t.Errorf("expected %s, got %s", secret, retrieved)
	}

	// Clean up
	DeleteAPIKey(service, account)
}

func TestKeychainDelete(t *testing.T) {
	service := "mycli-test"
	account := "test-delete"
	secret := "test-key"

	// Set key
	err := SetAPIKey(service, account, secret)
	if err != nil {
		t.Skipf("keychain not available: %v", err)
	}

	// Delete key
	err = DeleteAPIKey(service, account)
	if err != nil {
		t.Fatalf("failed to delete API key: %v", err)
	}

	// Verify deleted
	_, err = GetAPIKey(service, account)
	if err == nil {
		t.Error("expected error for deleted key, got nil")
	}
}
