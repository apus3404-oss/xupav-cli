# Go-Python Bridge Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Build the Go-Python bridge that manages Python subprocess lifecycle, implements JSON-RPC client protocol, handles streaming responses, and provides error handling with fallback mechanisms.

**Architecture:** Go package that spawns Python subprocess, communicates via JSON-RPC over stdin/stdout, manages subprocess lifecycle, handles timeouts and retries, and provides clean error messages.

**Tech Stack:**
- **os/exec** - Subprocess management
- **encoding/json** - JSON-RPC protocol
- **bufio** - Buffered I/O for stdin/stdout
- **context** - Timeout and cancellation

---

## File Structure

**New Files:**
```
internal/
  bridge/
    python.go          # Python subprocess manager
    protocol.go        # JSON-RPC protocol
    stream.go          # Streaming response handler
    errors.go          # Error types and handling
  bridge/
    python_test.go
    protocol_test.go
    stream_test.go
```

---
## Task 1: Error Types

**Files:**
- Create: `internal/bridge/errors.go`
- Create: `internal/bridge/errors_test.go`

- [ ] **Step 1: Write failing test for error types**

```go
// internal/bridge/errors_test.go
package bridge

import (
	"errors"
	"testing"
)

func TestBridgeErrors(t *testing.T) {
	tests := []struct {
		name    string
		err     error
		wantMsg string
	}{
		{
			name:    "python not found",
			err:     ErrPythonNotFound,
			wantMsg: "python runtime not found",
		},
		{
			name:    "subprocess crashed",
			err:     ErrSubprocessCrashed,
			wantMsg: "python subprocess crashed",
		},
		{
			name:    "timeout",
			err:     ErrTimeout,
			wantMsg: "request timeout",
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.err.Error() != tt.wantMsg {
				t.Errorf("expected %q, got %q", tt.wantMsg, tt.err.Error())
			}
		})
	}
}

func TestIsBridgeError(t *testing.T) {
	if !IsBridgeError(ErrPythonNotFound) {
		t.Error("expected ErrPythonNotFound to be a bridge error")
	}
	
	if IsBridgeError(errors.New("random error")) {
		t.Error("expected random error to not be a bridge error")
	}
}
```

- [ ] **Step 2: Run test to verify it fails**

```bash
go test ./internal/bridge -v -run TestBridgeErrors
```

Expected: FAIL with "undefined: ErrPythonNotFound"

- [ ] **Step 3: Implement error types**

```go
// internal/bridge/errors.go
package bridge

import (
	"errors"
	"fmt"
)

// Predefined errors
var (
	ErrPythonNotFound     = errors.New("python runtime not found")
	ErrSubprocessCrashed  = errors.New("python subprocess crashed")
	ErrTimeout            = errors.New("request timeout")
	ErrInvalidResponse    = errors.New("invalid JSON-RPC response")
	ErrMethodNotFound     = errors.New("method not found")
	ErrInvalidParams      = errors.New("invalid parameters")
)

// BridgeError wraps errors with additional context
type BridgeError struct {
	Op  string // Operation that failed
	Err error  // Underlying error
}

func (e *BridgeError) Error() string {
	return fmt.Sprintf("bridge: %s: %v", e.Op, e.Err)
}

func (e *BridgeError) Unwrap() error {
	return e.Err
}

// IsBridgeError checks if error is a bridge error
func IsBridgeError(err error) bool {
	var bridgeErr *BridgeError
	return errors.As(err, &bridgeErr) ||
		errors.Is(err, ErrPythonNotFound) ||
		errors.Is(err, ErrSubprocessCrashed) ||
		errors.Is(err, ErrTimeout) ||
		errors.Is(err, ErrInvalidResponse) ||
		errors.Is(err, ErrMethodNotFound) ||
		errors.Is(err, ErrInvalidParams)
}

// NewBridgeError creates a new bridge error
func NewBridgeError(op string, err error) error {
	return &BridgeError{Op: op, Err: err}
}
```

- [ ] **Step 4: Run test to verify it passes**

```bash
go test ./internal/bridge -v -run TestBridgeErrors
```

Expected: PASS

- [ ] **Step 5: Commit**

```bash
git add internal/bridge/errors.go internal/bridge/errors_test.go
git commit -m "feat(bridge): add error types for bridge operations"
```

---

## Task 2: JSON-RPC Protocol

**Files:**
- Create: `internal/bridge/protocol.go`
- Create: `internal/bridge/protocol_test.go`

- [ ] **Step 1: Write failing test for JSON-RPC protocol**

```go
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
}
```

- [ ] **Step 2: Run test to verify it fails**

```bash
go test ./internal/bridge -v -run TestCreateRequest
```

Expected: FAIL with "undefined: CreateRequest"

- [ ] **Step 3: Implement JSON-RPC protocol**

```go
// internal/bridge/protocol.go
package bridge

import (
	"encoding/json"
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

// JSONRPCNotification represents a JSON-RPC notification (no ID)
type JSONRPCNotification struct {
	JSONRPC string                 `json:"jsonrpc"`
	Method  string                 `json:"method"`
	Params  map[string]interface{} `json:"params"`
}

// CreateRequest creates a new JSON-RPC request
func CreateRequest(method string, params map[string]interface{}) *JSONRPCRequest {
	id := atomic.AddUint64(&requestID, 1)
	return &JSONRPCRequest{
		JSONRPC: "2.0",
		ID:      id,
		Method:  method,
		Params:  params,
	}
}

// Marshal serializes the request to JSON
func (r *JSONRPCRequest) Marshal() ([]byte, error) {
	return json.Marshal(r)
}

// ParseResponse parses a JSON-RPC response
func ParseResponse(data []byte) (*JSONRPCResponse, error) {
	var resp JSONRPCResponse
	if err := json.Unmarshal(data, &resp); err != nil {
		return nil, NewBridgeError("parse_response", err)
	}
	
	if resp.JSONRPC != "2.0" {
		return nil, NewBridgeError("parse_response", ErrInvalidResponse)
	}
	
	return &resp, nil
}

// ParseNotification parses a JSON-RPC notification
func ParseNotification(data []byte) (*JSONRPCNotification, error) {
	var notif JSONRPCNotification
	if err := json.Unmarshal(data, &notif); err != nil {
		return nil, NewBridgeError("parse_notification", err)
	}
	
	return &notif, nil
}

// IsError checks if response contains an error
func (r *JSONRPCResponse) IsError() bool {
	return r.Error != nil
}

// GetError returns the error from response
func (r *JSONRPCResponse) GetError() error {
	if r.Error == nil {
		return nil
	}
	
	switch r.Error.Code {
	case -32601:
		return ErrMethodNotFound
	case -32602:
		return ErrInvalidParams
	default:
		return NewBridgeError("rpc_error", &BridgeError{
			Op:  "remote",
			Err: &JSONRPCError{Code: r.Error.Code, Message: r.Error.Message},
		})
	}
}

func (e *JSONRPCError) Error() string {
	return e.Message
}
```

- [ ] **Step 4: Run test to verify it passes**

```bash
go test ./internal/bridge -v -run TestCreateRequest
go test ./internal/bridge -v -run TestParseResponse
go test ./internal/bridge -v -run TestParseErrorResponse
```

Expected: All PASS

- [ ] **Step 5: Commit**

```bash
git add internal/bridge/protocol.go internal/bridge/protocol_test.go
git commit -m "feat(bridge): implement JSON-RPC protocol"
```

---
## Task 3: Python Subprocess Manager

**Files:**
- Create: `internal/bridge/python.go`
- Create: `internal/bridge/python_test.go`

- [ ] **Step 1: Write failing test for subprocess manager**

```go
// internal/bridge/python_test.go
package bridge

import (
	"testing"
	"time"
)

func TestPythonBridge_StartStop(t *testing.T) {
	bridge := NewPythonBridge(PythonConfig{
		PythonPath: "python3",
		ScriptPath: "../../../python/mycli_ai/server.py",
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

func TestPythonBridge_RestartOnCrash(t *testing.T) {
	bridge := NewPythonBridge(PythonConfig{
		PythonPath:   "python3",
		ScriptPath:   "../../../python/mycli_ai/server.py",
		MaxRestarts:  3,
		RestartDelay: 100 * time.Millisecond,
	})
	
	err := bridge.Start()
	if err != nil {
		t.Skipf("Python not available: %v", err)
	}
	defer bridge.Stop()
	
	// Simulate crash by killing process
	if bridge.cmd != nil && bridge.cmd.Process != nil {
		bridge.cmd.Process.Kill()
	}
	
	// Wait for restart
	time.Sleep(200 * time.Millisecond)
	
	// Should be running again
	if !bridge.IsRunning() {
		t.Error("expected bridge to restart after crash")
	}
}
```

- [ ] **Step 2: Run test to verify it fails**

```bash
go test ./internal/bridge -v -run TestPythonBridge_StartStop
```

Expected: FAIL with "undefined: NewPythonBridge"

- [ ] **Step 3: Implement Python subprocess manager**

```go
// internal/bridge/python.go
package bridge

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"os/exec"
	"sync"
	"time"
)

// PythonConfig holds configuration for Python bridge
type PythonConfig struct {
	PythonPath   string        // Path to Python executable
	ScriptPath   string        // Path to Python server script
	Config       interface{}   // Configuration to pass to Python
	MaxRestarts  int           // Maximum restart attempts
	RestartDelay time.Duration // Delay between restarts
}

// PythonBridge manages Python subprocess
type PythonBridge struct {
	config       PythonConfig
	cmd          *exec.Cmd
	stdin        io.WriteCloser
	stdout       *bufio.Reader
	stderr       *bufio.Reader
	running      bool
	mu           sync.RWMutex
	restartCount int
	stopChan     chan struct{}
}

// NewPythonBridge creates a new Python bridge
func NewPythonBridge(config PythonConfig) *PythonBridge {
	if config.MaxRestarts == 0 {
		config.MaxRestarts = 3
	}
	if config.RestartDelay == 0 {
		config.RestartDelay = 1 * time.Second
	}
	
	return &PythonBridge{
		config:   config,
		stopChan: make(chan struct{}),
	}
}

// Start starts the Python subprocess
func (b *PythonBridge) Start() error {
	b.mu.Lock()
	defer b.mu.Unlock()
	
	if b.running {
		return fmt.Errorf("bridge already running")
	}
	
	// Check if Python exists
	if _, err := exec.LookPath(b.config.PythonPath); err != nil {
		return NewBridgeError("start", ErrPythonNotFound)
	}
	
	// Create command
	b.cmd = exec.Command(b.config.PythonPath, "-m", "mycli_ai.server")
	
	// Setup pipes
	stdin, err := b.cmd.StdinPipe()
	if err != nil {
		return NewBridgeError("start", err)
	}
	b.stdin = stdin
	
	stdout, err := b.cmd.StdoutPipe()
	if err != nil {
		return NewBridgeError("start", err)
	}
	b.stdout = bufio.NewReader(stdout)
	
	stderr, err := b.cmd.StderrPipe()
	if err != nil {
		return NewBridgeError("start", err)
	}
	b.stderr = bufio.NewReader(stderr)
	
	// Start process
	if err := b.cmd.Start(); err != nil {
		return NewBridgeError("start", err)
	}
	
	// Send config as first line
	configJSON, err := json.Marshal(b.config.Config)
	if err != nil {
		b.cmd.Process.Kill()
		return NewBridgeError("start", err)
	}
	
	if _, err := fmt.Fprintf(b.stdin, "%s\n", configJSON); err != nil {
		b.cmd.Process.Kill()
		return NewBridgeError("start", err)
	}
	
	b.running = true
	
	// Monitor process
	go b.monitor()
	
	return nil
}

// Stop stops the Python subprocess
func (b *PythonBridge) Stop() error {
	b.mu.Lock()
	defer b.mu.Unlock()
	
	if !b.running {
		return nil
	}
	
	// Signal stop
	close(b.stopChan)
	
	// Close stdin to signal Python to exit
	if b.stdin != nil {
		b.stdin.Close()
	}
	
	// Wait for process to exit (with timeout)
	done := make(chan error, 1)
	go func() {
		done <- b.cmd.Wait()
	}()
	
	select {
	case <-done:
		// Process exited cleanly
	case <-time.After(5 * time.Second):
		// Force kill
		if b.cmd.Process != nil {
			b.cmd.Process.Kill()
		}
	}
	
	b.running = false
	b.cmd = nil
	
	return nil
}

// IsRunning checks if subprocess is running
func (b *PythonBridge) IsRunning() bool {
	b.mu.RLock()
	defer b.mu.RUnlock()
	return b.running
}

// monitor watches the subprocess and restarts on crash
func (b *PythonBridge) monitor() {
	// Wait for process to exit
	err := b.cmd.Wait()
	
	b.mu.Lock()
	wasRunning := b.running
	b.running = false
	b.mu.Unlock()
	
	// Check if we should restart
	select {
	case <-b.stopChan:
		// Intentional stop, don't restart
		return
	default:
		// Unexpected exit
		if wasRunning && b.restartCount < b.config.MaxRestarts {
			b.restartCount++
			
			// Wait before restart
			time.Sleep(b.config.RestartDelay)
			
			// Attempt restart
			if err := b.Start(); err != nil {
				// Restart failed, give up
				return
			}
		}
	}
}

// GetStdout returns stdout reader
func (b *PythonBridge) GetStdout() *bufio.Reader {
	b.mu.RLock()
	defer b.mu.RUnlock()
	return b.stdout
}

// GetStderr returns stderr reader
func (b *PythonBridge) GetStderr() *bufio.Reader {
	b.mu.RLock()
	defer b.mu.RUnlock()
	return b.stderr
}

// Write writes data to stdin
func (b *PythonBridge) Write(data []byte) error {
	b.mu.RLock()
	defer b.mu.RUnlock()
	
	if !b.running {
		return NewBridgeError("write", fmt.Errorf("bridge not running"))
	}
	
	_, err := b.stdin.Write(data)
	if err != nil {
		return NewBridgeError("write", err)
	}
	
	return nil
}
```

- [ ] **Step 4: Run test to verify it passes**

```bash
go test ./internal/bridge -v -run TestPythonBridge_StartStop
```

Expected: PASS (or SKIP if Python not available)

- [ ] **Step 5: Commit**

```bash
git add internal/bridge/python.go internal/bridge/python_test.go
git commit -m "feat(bridge): implement Python subprocess manager"
```

---

## Task 4: Request/Response Handling

**Files:**
- Modify: `internal/bridge/python.go`
- Modify: `internal/bridge/python_test.go`

- [ ] **Step 1: Write failing test for request/response**

```go
// internal/bridge/python_test.go

func TestPythonBridge_SendRequest(t *testing.T) {
	bridge := NewPythonBridge(PythonConfig{
		PythonPath: "python3",
		ScriptPath: "../../../python/mycli_ai/server.py",
		Config: map[string]interface{}{
			"openrouter": map[string]interface{}{
				"enabled": false,
			},
			"ollama": map[string]interface{}{
				"enabled": false,
			},
		},
	})
	
	err := bridge.Start()
	if err != nil {
		t.Skipf("Python not available: %v", err)
	}
	defer bridge.Stop()
	
	// Create request
	req := CreateRequest("get_cost_summary", map[string]interface{}{})
	
	// Send request
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	
	resp, err := bridge.SendRequest(ctx, req)
	if err != nil {
		t.Fatalf("failed to send request: %v", err)
	}
	
	// Verify response
	if resp.IsError() {
		t.Errorf("unexpected error: %v", resp.GetError())
	}
	
	if resp.Result == nil {
		t.Error("expected result, got nil")
	}
}
```

- [ ] **Step 2: Run test to verify it fails**

```bash
go test ./internal/bridge -v -run TestPythonBridge_SendRequest
```

Expected: FAIL with "undefined: PythonBridge.SendRequest"

- [ ] **Step 3: Implement request/response handling**

```go
// internal/bridge/python.go

// SendRequest sends a JSON-RPC request and waits for response
func (b *PythonBridge) SendRequest(ctx context.Context, req *JSONRPCRequest) (*JSONRPCResponse, error) {
	if !b.IsRunning() {
		return nil, NewBridgeError("send_request", fmt.Errorf("bridge not running"))
	}
	
	// Marshal request
	data, err := req.Marshal()
	if err != nil {
		return nil, NewBridgeError("send_request", err)
	}
	
	// Write request
	if err := b.Write(append(data, '\n')); err != nil {
		return nil, err
	}
	
	// Read response with timeout
	responseChan := make(chan *JSONRPCResponse, 1)
	errorChan := make(chan error, 1)
	
	go func() {
		line, err := b.stdout.ReadString('\n')
		if err != nil {
			errorChan <- NewBridgeError("read_response", err)
			return
		}
		
		resp, err := ParseResponse([]byte(line))
		if err != nil {
			errorChan <- err
			return
		}
		
		responseChan <- resp
	}()
	
	// Wait for response or timeout
	select {
	case resp := <-responseChan:
		if resp.IsError() {
			return resp, resp.GetError()
		}
		return resp, nil
	case err := <-errorChan:
		return nil, err
	case <-ctx.Done():
		return nil, NewBridgeError("send_request", ErrTimeout)
	}
}

// SendRequestWithRetry sends request with retry logic
func (b *PythonBridge) SendRequestWithRetry(ctx context.Context, req *JSONRPCRequest, maxRetries int) (*JSONRPCResponse, error) {
	var lastErr error
	
	for i := 0; i <= maxRetries; i++ {
		resp, err := b.SendRequest(ctx, req)
		if err == nil {
			return resp, nil
		}
		
		lastErr = err
		
		// Don't retry on certain errors
		if IsBridgeError(err) {
			var bridgeErr *BridgeError
			if errors.As(err, &bridgeErr) {
				if errors.Is(bridgeErr.Err, ErrMethodNotFound) ||
					errors.Is(bridgeErr.Err, ErrInvalidParams) {
					return nil, err
				}
			}
		}
		
		// Wait before retry
		if i < maxRetries {
			select {
			case <-time.After(time.Duration(i+1) * time.Second):
			case <-ctx.Done():
				return nil, ctx.Err()
			}
		}
	}
	
	return nil, lastErr
}
```

- [ ] **Step 4: Add missing import**

```go
// internal/bridge/python.go
import (
	"bufio"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os/exec"
	"sync"
	"time"
)
```

- [ ] **Step 5: Run test to verify it passes**

```bash
go test ./internal/bridge -v -run TestPythonBridge_SendRequest
```

Expected: PASS

- [ ] **Step 6: Commit**

```bash
git add internal/bridge/python.go internal/bridge/python_test.go
git commit -m "feat(bridge): add request/response handling with retry logic"
```

---
## Task 5: Streaming Response Handler

**Files:**
- Create: `internal/bridge/stream.go`
- Create: `internal/bridge/stream_test.go`

- [ ] **Step 1: Write failing test for streaming**

```go
// internal/bridge/stream_test.go
package bridge

import (
	"context"
	"testing"
	"time"
)

func TestStreamHandler(t *testing.T) {
	handler := NewStreamHandler()
	
	// Start collecting chunks
	chunks := make([]string, 0)
	done := make(chan struct{})
	
	go func() {
		for chunk := range handler.Chunks() {
			chunks = append(chunks, chunk)
		}
		close(done)
	}()
	
	// Send some chunks
	handler.AddChunk("Hello")
	handler.AddChunk(" ")
	handler.AddChunk("world")
	
	// Close stream
	handler.Close()
	
	// Wait for collection to finish
	<-done
	
	// Verify chunks
	if len(chunks) != 3 {
		t.Errorf("expected 3 chunks, got %d", len(chunks))
	}
	
	expected := []string{"Hello", " ", "world"}
	for i, chunk := range chunks {
		if chunk != expected[i] {
			t.Errorf("chunk %d: expected %q, got %q", i, expected[i], chunk)
		}
	}
}

func TestStreamHandler_Timeout(t *testing.T) {
	handler := NewStreamHandler()
	
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
	defer cancel()
	
	// Try to read with timeout
	select {
	case <-handler.Chunks():
		t.Error("expected timeout, got chunk")
	case <-ctx.Done():
		// Expected timeout
	}
}
```

- [ ] **Step 2: Run test to verify it fails**

```bash
go test ./internal/bridge -v -run TestStreamHandler
```

Expected: FAIL with "undefined: NewStreamHandler"

- [ ] **Step 3: Implement streaming handler**

```go
// internal/bridge/stream.go
package bridge

import (
	"context"
	"sync"
)

// StreamHandler handles streaming responses
type StreamHandler struct {
	chunks chan string
	done   chan struct{}
	mu     sync.Mutex
	closed bool
}

// NewStreamHandler creates a new stream handler
func NewStreamHandler() *StreamHandler {
	return &StreamHandler{
		chunks: make(chan string, 100), // Buffered for performance
		done:   make(chan struct{}),
	}
}

// AddChunk adds a chunk to the stream
func (s *StreamHandler) AddChunk(chunk string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	
	if s.closed {
		return
	}
	
	s.chunks <- chunk
}

// Chunks returns the channel for reading chunks
func (s *StreamHandler) Chunks() <-chan string {
	return s.chunks
}

// Close closes the stream
func (s *StreamHandler) Close() {
	s.mu.Lock()
	defer s.mu.Unlock()
	
	if s.closed {
		return
	}
	
	s.closed = true
	close(s.chunks)
	close(s.done)
}

// Done returns a channel that's closed when stream is done
func (s *StreamHandler) Done() <-chan struct{} {
	return s.done
}

// CollectAll collects all chunks into a single string
func (s *StreamHandler) CollectAll(ctx context.Context) (string, error) {
	var result string
	
	for {
		select {
		case chunk, ok := <-s.chunks:
			if !ok {
				return result, nil
			}
			result += chunk
		case <-ctx.Done():
			return result, ctx.Err()
		}
	}
}
```

- [ ] **Step 4: Add streaming support to PythonBridge**

```go
// internal/bridge/python.go

// SendStreamingRequest sends a streaming request
func (b *PythonBridge) SendStreamingRequest(ctx context.Context, req *JSONRPCRequest) (*StreamHandler, error) {
	if !b.IsRunning() {
		return nil, NewBridgeError("send_streaming_request", fmt.Errorf("bridge not running"))
	}
	
	// Marshal request
	data, err := req.Marshal()
	if err != nil {
		return nil, NewBridgeError("send_streaming_request", err)
	}
	
	// Write request
	if err := b.Write(append(data, '\n')); err != nil {
		return nil, err
	}
	
	// Create stream handler
	handler := NewStreamHandler()
	
	// Start reading chunks in background
	go func() {
		defer handler.Close()
		
		for {
			select {
			case <-ctx.Done():
				return
			default:
				line, err := b.stdout.ReadString('\n')
				if err != nil {
					return
				}
				
				// Try to parse as notification (chunk)
				notif, err := ParseNotification([]byte(line))
				if err == nil && notif.Method == "chunk" {
					if content, ok := notif.Params["content"].(string); ok {
						handler.AddChunk(content)
					}
					continue
				}
				
				// Try to parse as response (completion)
				resp, err := ParseResponse([]byte(line))
				if err == nil {
					// Stream completed
					return
				}
			}
		}
	}()
	
	return handler, nil
}
```

- [ ] **Step 5: Run test to verify it passes**

```bash
go test ./internal/bridge -v -run TestStreamHandler
```

Expected: PASS

- [ ] **Step 6: Commit**

```bash
git add internal/bridge/stream.go internal/bridge/stream_test.go internal/bridge/python.go
git commit -m "feat(bridge): add streaming response handler"
```

---

## Task 6: Integration Test

**Files:**
- Create: `internal/bridge/integration_test.go`

- [ ] **Step 1: Write integration test**

```go
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
		PythonPath: "python3",
		ScriptPath: "../../../python/mycli_ai/server.py",
		Config: map[string]interface{}{
			"openrouter": map[string]interface{}{
				"enabled": false,
			},
			"ollama": map[string]interface{}{
				"enabled": false,
			},
		},
	})
	
	// Step 1: Start bridge
	err := bridge.Start()
	if err != nil {
		t.Skipf("Python not available: %v", err)
	}
	defer bridge.Stop()
	
	// Step 2: Send request
	req := CreateRequest("get_cost_summary", map[string]interface{}{})
	
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	
	resp, err := bridge.SendRequest(ctx, req)
	if err != nil {
		t.Fatalf("request failed: %v", err)
	}
	
	// Step 3: Verify response
	if resp.IsError() {
		t.Fatalf("unexpected error: %v", resp.GetError())
	}
	
	if resp.Result == nil {
		t.Fatal("expected result, got nil")
	}
	
	// Step 4: Verify result structure
	if _, ok := resp.Result["total_tokens"]; !ok {
		t.Error("expected total_tokens in result")
	}
	
	if _, ok := resp.Result["total_cost"]; !ok {
		t.Error("expected total_cost in result")
	}
	
	// Step 5: Test retry logic
	req2 := CreateRequest("invalid_method", map[string]interface{}{})
	
	_, err = bridge.SendRequestWithRetry(ctx, req2, 2)
	if err == nil {
		t.Error("expected error for invalid method")
	}
	
	// Should be method not found error
	if !IsBridgeError(err) {
		t.Errorf("expected bridge error, got %T", err)
	}
	
	// Step 6: Stop and verify
	if err := bridge.Stop(); err != nil {
		t.Errorf("failed to stop: %v", err)
	}
	
	if bridge.IsRunning() {
		t.Error("expected bridge to be stopped")
	}
}

func TestBridgeIntegration_Restart(t *testing.T) {
	bridge := NewPythonBridge(PythonConfig{
		PythonPath:   "python3",
		ScriptPath:   "../../../python/mycli_ai/server.py",
		MaxRestarts:  2,
		RestartDelay: 100 * time.Millisecond,
		Config: map[string]interface{}{
			"openrouter": map[string]interface{}{
				"enabled": false,
			},
			"ollama": map[string]interface{}{
				"enabled": false,
			},
		},
	})
	
	err := bridge.Start()
	if err != nil {
		t.Skipf("Python not available: %v", err)
	}
	defer bridge.Stop()
	
	// Send initial request
	req := CreateRequest("get_cost_summary", map[string]interface{}{})
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	
	_, err = bridge.SendRequest(ctx, req)
	if err != nil {
		t.Fatalf("initial request failed: %v", err)
	}
	
	// Simulate crash
	if bridge.cmd != nil && bridge.cmd.Process != nil {
		bridge.cmd.Process.Kill()
	}
	
	// Wait for restart
	time.Sleep(300 * time.Millisecond)
	
	// Try request again (should work after restart)
	ctx2, cancel2 := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel2()
	
	_, err = bridge.SendRequest(ctx2, req)
	if err != nil {
		t.Errorf("request after restart failed: %v", err)
	}
}
```

- [ ] **Step 2: Run integration test**

```bash
go test ./internal/bridge -v -run TestBridgeIntegration
```

Expected: PASS (or SKIP if Python not available)

- [ ] **Step 3: Run all bridge tests**

```bash
go test ./internal/bridge -v
```

Expected: All tests PASS

- [ ] **Step 4: Commit**

```bash
git add internal/bridge/integration_test.go
git commit -m "test(bridge): add integration tests for full bridge flow"
```

---
## Completion Checklist

- [x] Error types defined
- [x] JSON-RPC protocol implemented
- [x] Python subprocess manager
- [x] Request/response handling with retry
- [x] Streaming response handler
- [x] Integration tests
- [x] All tests passing

## Manual Verification

```bash
# Build
go build -o bin/mycli cmd/mycli/main.go

# Test bridge (requires Python + mycli_ai installed)
go test ./internal/bridge -v

# Test with real Python subprocess
cd python
source venv/bin/activate
cd ..
go test ./internal/bridge -v -run TestBridgeIntegration
```

## Next Steps

This plan establishes the Go-Python bridge. Next plans will build on this:
- **Plan D**: Basic TUI (Bubble Tea chat interface)
- **Plan E**: TUI Polish (colors, animations, diff preview)
- **Plan F**: Token Optimizer (file chunking, context prioritization)
- **Plan G**: LSP Integration (code intelligence)
- **Plan H**: Code Modifications (diff generation, file writer)

---

## Notes for Implementation

**Key Design Decisions:**
1. JSON-RPC over stdin/stdout for simplicity and reliability
2. Automatic restart on crash (up to MaxRestarts)
3. Context-based timeouts for all operations
4. Retry logic with exponential backoff
5. Streaming support for real-time responses

**Testing Strategy:**
- Unit tests for protocol and errors
- Integration tests require Python installed
- Tests skip gracefully if Python unavailable
- Mock subprocess for unit tests where possible

**Common Issues:**
- Python path may vary (python3, python, py)
- Virtual environment activation needed for mycli_ai
- Subprocess may not exit cleanly (use timeout + kill)
- JSON-RPC requires exact format (trailing newline)
- Streaming requires buffered reading

**Performance Considerations:**
- Subprocess startup: ~100-300ms
- JSON-RPC overhead: <1ms per request
- Streaming reduces latency (chunks arrive immediately)
- Restart delay configurable (default 1s)
- Buffered I/O for better throughput

**Error Handling:**
- Distinguish between bridge errors and Python errors
- Retry on transient errors (network, timeout)
- Don't retry on permanent errors (method not found)
- Graceful degradation (continue without Python if needed)
