// internal/bridge/errors.go
package bridge

import (
	"errors"
	"fmt"
)

// Predefined errors
var (
	ErrPythonNotFound    = errors.New("python runtime not found")
	ErrSubprocessCrashed = errors.New("python subprocess crashed")
	ErrTimeout           = errors.New("request timeout")
	ErrInvalidResponse   = errors.New("invalid JSON-RPC response")
	ErrMethodNotFound    = errors.New("method not found")
	ErrInvalidParams     = errors.New("invalid parameters")
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
	return &BridgeError{
		Op:  op,
		Err: err,
	}
}
