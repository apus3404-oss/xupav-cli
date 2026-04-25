// internal/bridge/protocol_test.go
package bridge

import (
	"encoding/json"
	"testing"
)

func TestCreateRequest(t *testing.T) {
	req := CreateRequest("chat", map[string]interface{}{
		"message": "test",
		"model":   "deepseek/deepseek-r1",
	})

	if req.JSONRPC != "2.0" {
		t.Errorf("expected jsonrpc 2.0, got %s", req.JSONRPC)
	}

	if req.Method != "chat" {
		t.Errorf("expected method chat, got %s", req.Method)
	}

	if req.ID == 0 {
		t.Error("expected non-zero ID")
	}

	if req.Params["message"] != "test" {
		t.Errorf("expected message test, got %v", req.Params["message"])
	}
}

func TestParseResponse(t *testing.T) {
	responseJSON := `{
		"jsonrpc": "2.0",
		"id": 1,
		"result": {
			"content": "test response",
			"tokens": 100,
			"cost": 0.001
		}
	}`

	var resp JSONRPCResponse
	if err := json.Unmarshal([]byte(responseJSON), &resp); err != nil {
		t.Fatalf("failed to unmarshal: %v", err)
	}

	if resp.JSONRPC != "2.0" {
		t.Errorf("expected jsonrpc 2.0, got %s", resp.JSONRPC)
	}

	if resp.Error != nil {
		t.Errorf("expected no error, got %v", resp.Error)
	}

	if resp.Result["content"] != "test response" {
		t.Errorf("expected content 'test response', got %v", resp.Result["content"])
	}
}

func TestParseErrorResponse(t *testing.T) {
	errorJSON := `{
		"jsonrpc": "2.0",
		"id": 1,
		"error": {
			"code": -32601,
			"message": "Method not found"
		}
	}`

	var resp JSONRPCResponse
	if err := json.Unmarshal([]byte(errorJSON), &resp); err != nil {
		t.Fatalf("failed to unmarshal: %v", err)
	}

	if resp.Error == nil {
		t.Fatal("expected error, got nil")
	}

	if resp.Error.Code != -32601 {
		t.Errorf("expected code -32601, got %d", resp.Error.Code)
	}

	if resp.Error.Message != "Method not found" {
		t.Errorf("expected message 'Method not found', got %s", resp.Error.Message)
	}
}

func TestRequestIDIncrement(t *testing.T) {
	req1 := CreateRequest("test", nil)
	req2 := CreateRequest("test", nil)

	if req2.ID <= req1.ID {
		t.Errorf("expected req2.ID (%d) > req1.ID (%d)", req2.ID, req1.ID)
	}
}
