// internal/bridge/integration_test.go
package bridge

import (
	"context"
	"testing"
	"time"
)

func TestBridgeIntegration(t *testing.T) {
	// This test requires Python and mycli_ai package installed
	bridge := NewPythonBridge(PythonConfig{
		PythonPath: "python",
		ScriptPath: "../../python/mycli_ai/server.py",
		Config: map[string]interface{}{
			"provider": "openrouter",
			"api_key":  "test-key",
		},
	})

	// Step 1: Start bridge
	err := bridge.Start()
	if err != nil {
		t.Skipf("Python not available: %v", err)
	}
	defer bridge.Stop()

	// Give Python time to start
	time.Sleep(100 * time.Millisecond)

	// Step 2: Verify bridge is running
	if !bridge.IsRunning() {
		t.Fatal("expected bridge to be running")
	}

	// Step 3: Send request
	req := CreateRequest("chat", map[string]interface{}{
		"message": "test",
		"model":   "deepseek/deepseek-r1",
	})

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	resp, err := bridge.SendRequest(ctx, req)
	if err != nil {
		// Expected to fail without real API key
		t.Logf("Request failed (expected without API key): %v", err)
	} else {
		// If we got a response, verify structure
		if resp.JSONRPC != "2.0" {
			t.Errorf("expected jsonrpc 2.0, got %s", resp.JSONRPC)
		}
	}

	// Step 4: Test invalid method
	req2 := CreateRequest("invalid_method", map[string]interface{}{})

	ctx2, cancel2 := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel2()

	_, err = bridge.SendRequest(ctx2, req2)
	if err == nil {
		t.Error("expected error for invalid method")
	}

	// Step 5: Stop and verify
	if err := bridge.Stop(); err != nil {
		t.Errorf("failed to stop: %v", err)
	}

	if bridge.IsRunning() {
		t.Error("expected bridge to be stopped")
	}
}

func TestBridgeIntegration_MultipleRequests(t *testing.T) {
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

	time.Sleep(100 * time.Millisecond)

	// Send multiple requests
	for i := 0; i < 3; i++ {
		req := CreateRequest("chat", map[string]interface{}{
			"message": "test",
			"model":   "deepseek/deepseek-r1",
		})

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		_, err := bridge.SendRequest(ctx, req)
		cancel()

		if err != nil {
			t.Logf("Request %d failed (expected): %v", i+1, err)
		}
	}
}

func TestBridgeIntegration_ErrorHandling(t *testing.T) {
	bridge := NewPythonBridge(PythonConfig{
		PythonPath: "python",
		ScriptPath: "../../python/mycli_ai/server.py",
		Config: map[string]interface{}{
			"provider": "openrouter",
			"api_key":  "test-key",
		},
	})

	// Test sending request before start
	req := CreateRequest("chat", map[string]interface{}{
		"message": "test",
		"model":   "test",
	})

	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	_, err := bridge.SendRequest(ctx, req)
	if err == nil {
		t.Error("expected error when bridge not running")
	}

	if !IsBridgeError(err) {
		t.Errorf("expected bridge error, got %T", err)
	}

	// Now start and test
	err = bridge.Start()
	if err != nil {
		t.Skipf("Python not available: %v", err)
	}
	defer bridge.Stop()

	time.Sleep(100 * time.Millisecond)

	// Test timeout
	ctx2, cancel2 := context.WithTimeout(context.Background(), 1*time.Millisecond)
	defer cancel2()

	time.Sleep(10 * time.Millisecond) // Ensure context is expired

	_, err = bridge.SendRequest(ctx2, req)
	if err == nil {
		t.Error("expected timeout error")
	}
}
