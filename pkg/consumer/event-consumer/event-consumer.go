// Package event_consumer provides an implementation of the consumer.Consumer interface.
// It fetches events in batches and processes them continuously.
package event_consumer

import (
	"go_link_storage/pkg/events"
	"log"
)

// Consumer implements the consumer.Consumer interface.
// It fetches events from a Fetcher and processes them using a Processor.
type Consumer struct {
	fetcher   events.Fetcher   // Source of events to fetch
	processor events.Processor // Processor for handling events
	batchSize int              // Number of events to fetch per batch
}

// New creates a new event consumer with the given fetcher, processor, and batch size.
func New(fetcher events.Fetcher,
	processor events.Processor,
	batchSize int) Consumer {

	return Consumer{
		fetcher:   fetcher,
		processor: processor,
		batchSize: batchSize,
	}
}

// Start begins consuming events in a continuous loop.
// It fetches events in batches and processes them, sleeping when no events are available.
func (c Consumer) Start() {
	c.fetcher.Start(c.handleEvents, c.batchSize)
}

// handleEvents processes a batch of events sequentially.
func (c *Consumer) handleEvents(events []events.Event) error {
	for _, event := range events {
		log.Printf("new event: %s", event.Text)

		if err := c.processor.Process(event); err != nil {
			log.Printf("couldn't handle event: %s", err.Error())

			continue
		}
	}

	return nil
}
