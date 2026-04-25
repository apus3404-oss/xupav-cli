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

func TestBridgeErrorWrapping(t *testing.T) {
	baseErr := errors.New("base error")
	bridgeErr := NewBridgeError("test_operation", baseErr)

	if !errors.Is(bridgeErr, baseErr) {
		t.Error("expected bridge error to wrap base error")
	}

	var be *BridgeError
	if !errors.As(bridgeErr, &be) {
		t.Error("expected error to be a BridgeError")
	}

	if be.Op != "test_operation" {
		t.Errorf("expected op %q, got %q", "test_operation", be.Op)
	}
}
