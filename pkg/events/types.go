// Package events defines the core event types and interfaces for the event system.
package events

// Fetcher defines the interface for fetching events from a source.
type Fetcher interface {
	Start(handleEventsCallback func(events []Event) error, batchSize int)
}

// Processor defines the interface for processing events.
type Processor interface {
	// Process handles a single event.
	Process(evt Event) error
}

type Client interface {
	SendMessage(chatID int, text string) error
}

// Type represents the type of an event.
type Type int

const (
	Unknown Type = iota // Unknown event type
	Message             // Message event type
)

// Event represents a single event in the system.
type Event struct {
	Type Type   // The type of the event
	Text string // The text content of the event
	Meta any    // Additional metadata associated with the event
}
