// Package tg_processor provides Telegram-specific event processing.
// It implements both events.Fetcher and events.Processor interfaces
// for handling Telegram bot updates.
package tg_processor

import (
	"errors"
	"go_link_storage/pkg/events"
	"go_link_storage/pkg/lib/e"
	"go_link_storage/pkg/storage"
)

// Processor handles Telegram events by fetching updates and processing messages.
// It implements both events.Fetcher and events.Processor interfaces.
type Processor struct {
	tg      events.Client   // Telegram API client
	offset  int             // Offset for fetching updates
	storage storage.Storage // Storage for saving pages
}

// Meta contains metadata associated with Telegram events.
type Meta struct {
	ChatID   int    // Telegram chat ID
	Username string // Telegram username
}

var (
	// ErrUnknownEventType is returned when an event type cannot be determined.
	ErrUnknownEventType = errors.New("unknown event type")
	// ErrUnknownMetaType is returned when event metadata has an unexpected type.
	ErrUnknownMetaType = errors.New("unknown meta type")
)

// New creates a new Telegram event processor with the given client and storage.
func New(client events.Client, storage storage.Storage) *Processor {
	return &Processor{
		tg:      client,
		storage: storage,
	}
}

// Process handles an event by routing it to the appropriate handler based on event type.
func (p *Processor) Process(event events.Event) error {
	switch event.Type {
	case events.Message:
		return p.processMessage(event)
	default:
		return e.Wrap("cannot process event", ErrUnknownEventType)
	}
}

// processMessage handles message events by extracting metadata and executing commands.
func (p *Processor) processMessage(event events.Event) error {
	meta, err := meta(event)
	if err != nil {
		return e.Wrap("cannot process message", err)
	}

	if err := p.doCmd(event.Text, meta.ChatID, meta.Username); err != nil {
		return e.Wrap("cannot process message", err)
	}

	return nil
}

// meta extracts Meta from an event, returning an error if the meta type is incorrect.
func meta(event events.Event) (Meta, error) {
	res, ok := event.Meta.(Meta)
	if !ok {
		return Meta{}, e.Wrap("cannot get meta", ErrUnknownEventType)
	}

	return res, nil
}
