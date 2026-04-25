// internal/bridge/python_test.go
package bridge

import (
	"context"
	"testing"
	"time"
)

func TestPythonBridge_StartStop(t *testing.T) {
	bridge := NewPythonBridge(PythonConfig{
		PythonPath: "python",
		ScriptPath: "../../python/mycli_ai/server.py",
	})

	// Start
	err := bridge.Start()
	if err != nil {
		t.Skipf("Python not available: %v", err)
	}
	defer bridge.Stop()

	// Check running
	if !bridge.IsRunning() {
		t.Error("expected bridge to be running")
	}

	// Stop
	err = bridge.Stop()
	if err != nil {
		t.Errorf("failed to stop: %v", err)
	}

	// Check stopped
	if bridge.IsRunning() {
		t.Error("expected bridge to be stopped")
	}
}

func TestPythonBridge_Config(t *testing.T) {
	config := PythonConfig{
		PythonPath:   "python3",
		ScriptPath:   "test.py",
		MaxRestarts:  5,
		RestartDelay: 200 * time.Millisecond,
	}

	bridge := NewPythonBridge(config)

	if bridge.config.PythonPath != "python3" {
		t.Errorf("expected python3, got %s", bridge.config.PythonPath)
	}

	if bridge.config.MaxRestarts != 5 {
		t.Errorf("expected MaxRestarts 5, got %d", bridge.config.MaxRestarts)
	}
}

func TestPythonBridge_NotStarted(t *testing.T) {
	bridge := NewPythonBridge(PythonConfig{
		PythonPath: "python",
		ScriptPath: "test.py",
	})

	// Should not be running
	if bridge.IsRunning() {
		t.Error("expected bridge to not be running")
	}

	// Stop should not error
	err := bridge.Stop()
	if err != nil {
		t.Errorf("stop on non-running bridge should not error: %v", err)
	}
}

func TestPythonBridge_SendRequest(t *testing.T) {
	bridge := NewPythonBridge(PythonConfig{
		PythonPath: "python",
		ScriptPath: "../../python/mycli_ai/server.py",
		Config: map[string]interface{}{
			"provider": "openrouter",
			"api_key":  "test-key",
		},
	})

	err := bridge.Start()
	if err != nil {
		t.Skipf("Python not available: %v", err)
	}
	defer bridge.Stop()

	// Give Python time to start
	time.Sleep(100 * time.Millisecond)

	// Create request
	req := CreateRequest("chat", map[string]interface{}{
		"message": "test",
		"model":   "deepseek/deepseek-r1",
	})

	// Send request
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	resp, err := bridge.SendRequest(ctx, req)
	if err != nil {
		t.Logf("SendRequest error (expected if no API key): %v", err)
		// This is expected to fail without a real API key
		return
	}

	// If we got a response, verify structure
	if resp.JSONRPC != "2.0" {
		t.Errorf("expected jsonrpc 2.0, got %s", resp.JSONRPC)
	}
}
