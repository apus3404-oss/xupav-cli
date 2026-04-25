// internal/bridge/protocol.go
package bridge

import (
	"sync/atomic"
)

var requestID uint64

// JSONRPCRequest represents a JSON-RPC 2.0 request
type JSONRPCRequest struct {
	JSONRPC string                 `json:"jsonrpc"`
	ID      uint64                 `json:"id"`
	Method  string                 `json:"method"`
	Params  map[string]interface{} `json:"params"`
}

// JSONRPCResponse represents a JSON-RPC 2.0 response
type JSONRPCResponse struct {
	JSONRPC string                 `json:"jsonrpc"`
	ID      uint64                 `json:"id"`
	Result  map[string]interface{} `json:"result,omitempty"`
	Error   *JSONRPCError          `json:"error,omitempty"`
}

// JSONRPCError represents a JSON-RPC error
type JSONRPCError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

// CreateRequest creates a new JSON-RPC request with auto-incrementing ID
func CreateRequest(method string, params map[string]interface{}) *JSONRPCRequest {
	id := atomic.AddUint64(&requestID, 1)

	if params == nil {
		params = make(map[string]interface{})
	}

	return &JSONRPCRequest{
		JSONRPC: "2.0",
		ID:      id,
		Method:  method,
		Params:  params,
	}
}

// IsError checks if the response contains an error
func (r *JSONRPCResponse) IsError() bool {
	return r.Error != nil
}

// GetError returns the error from the response
func (r *JSONRPCResponse) GetError() error {
	if r.Error == nil {
		return nil
	}

	switch r.Error.Code {
	case -32601:
		return ErrMethodNotFound
	case -32602:
		return ErrInvalidParams
	case -32600:
		return ErrInvalidResponse
	default:
		return NewBridgeError("rpc", &BridgeError{
			Op:  "response",
			Err: ErrInvalidResponse,
		})
	}
}
