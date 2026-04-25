// internal/bridge/stream.go
package bridge

import (
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

// IsClosed returns whether the stream is closed
func (s *StreamHandler) IsClosed() bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.closed
}
