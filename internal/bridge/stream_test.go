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

	handler.Close()
}

func TestStreamHandler_MultipleClose(t *testing.T) {
	handler := NewStreamHandler()

	handler.AddChunk("test")
	handler.Close()

	// Second close should not panic
	handler.Close()

	// Adding after close should not panic
	handler.AddChunk("ignored")
}

func TestStreamHandler_Done(t *testing.T) {
	handler := NewStreamHandler()

	// Done should not be closed initially
	select {
	case <-handler.Done():
		t.Error("expected Done to not be closed")
	default:
		// Expected
	}

	handler.Close()

	// Done should be closed after Close
	select {
	case <-handler.Done():
		// Expected
	case <-time.After(100 * time.Millisecond):
		t.Error("expected Done to be closed")
	}
}
