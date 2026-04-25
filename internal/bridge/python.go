// internal/bridge/python.go
package bridge

import (
	"bufio"
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
	if b.config.Config != nil {
		configJSON, err := json.Marshal(b.config.Config)
		if err != nil {
			b.cmd.Process.Kill()
			return NewBridgeError("start", err)
		}

		if _, err := fmt.Fprintf(b.stdin, "%s\n", configJSON); err != nil {
			b.cmd.Process.Kill()
			return NewBridgeError("start", err)
		}
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
	select {
	case <-b.stopChan:
		// Already closed
	default:
		close(b.stopChan)
	}

	// Close stdin to signal Python to exit
	if b.stdin != nil {
		b.stdin.Close()
	}

	// Wait for process to exit (with timeout)
	done := make(chan error, 1)
	go func() {
		if b.cmd != nil {
			done <- b.cmd.Wait()
		}
	}()

	select {
	case <-done:
		// Process exited cleanly
	case <-time.After(5 * time.Second):
		// Force kill
		if b.cmd != nil && b.cmd.Process != nil {
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
	if b.cmd != nil {
		b.cmd.Wait()
	}

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
