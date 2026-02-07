// Package consumer defines the interface for event consumers.
package consumer

// Consumer defines the interface for consuming and processing events.
type Consumer interface {
	// Start begins consuming and processing events.
	// It should run continuously until an error occurs or the consumer is stopped.
	Start() error
}
