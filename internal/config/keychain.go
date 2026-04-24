// internal/config/keychain.go
package config

import (
	"fmt"

	"github.com/zalando/go-keyring"
)

const (
	KeychainService = "mycli"
)

// SetAPIKey stores an API key in the system keychain
func SetAPIKey(service, account, secret string) error {
	err := keyring.Set(service, account, secret)
	if err != nil {
		return fmt.Errorf("failed to store API key in keychain: %w", err)
	}
	return nil
}

// GetAPIKey retrieves an API key from the system keychain
func GetAPIKey(service, account string) (string, error) {
	secret, err := keyring.Get(service, account)
	if err != nil {
		return "", fmt.Errorf("failed to retrieve API key from keychain: %w", err)
	}
	return secret, nil
}

// DeleteAPIKey removes an API key from the system keychain
func DeleteAPIKey(service, account string) error {
	err := keyring.Delete(service, account)
	if err != nil {
		return fmt.Errorf("failed to delete API key from keychain: %w", err)
	}
	return nil
}

// GetOpenRouterKey retrieves the OpenRouter API key
func GetOpenRouterKey() (string, error) {
	return GetAPIKey(KeychainService, "openrouter")
}

// SetOpenRouterKey stores the OpenRouter API key
func SetOpenRouterKey(key string) error {
	return SetAPIKey(KeychainService, "openrouter", key)
}

// GetOllamaKey retrieves the Ollama API key (if needed)
func GetOllamaKey() (string, error) {
	return GetAPIKey(KeychainService, "ollama")
}

// SetOllamaKey stores the Ollama API key
func SetOllamaKey(key string) error {
	return SetAPIKey(KeychainService, "ollama", key)
}
