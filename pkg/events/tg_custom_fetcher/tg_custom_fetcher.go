package tg_custom_fetcher

import (
	"errors"
	"go_link_storage/pkg/clients/tg_custom_client"
	"go_link_storage/pkg/events"
	"go_link_storage/pkg/events/tg_processor"
	"go_link_storage/pkg/lib/e"
	"log"
	"time"
)

type Fetcher struct {
	tg     *tg_custom_client.Client // Telegram API client
	offset int                      // Offset for fetching updates
}

var (
	// ErrUnknownEventType is returned when an event type cannot be determined.
	ErrUnknownEventType = errors.New("unknown event type")
	// ErrUnknownMetaType is returned when event metadata has an unexpected type.
	ErrUnknownMetaType = errors.New("unknown meta type")
)

// New creates a new Telegram event processor with the given client.
func New(client *tg_custom_client.Client) *Fetcher {
	return &Fetcher{
		tg: client,
	}
}

func (f *Fetcher) Start(handleEventsCallback func(events []events.Event) error, batchSize int) {
	for {
		cEvents, err := f.Fetch(batchSize)
		if err != nil {
			log.Printf("[ERR] consumer: %s", err)

			continue
		}

		if len(cEvents) == 0 {
			time.Sleep(1 * time.Second)

			continue
		}

		if err := handleEventsCallback(cEvents); err != nil {
			log.Println(err)

			continue
		}
	}
}

// Fetch retrieves Telegram updates and converts them to events.
// It updates the internal offset to track processed updates.
func (f *Fetcher) Fetch(limit int) ([]events.Event, error) {
	updates, err := f.tg.Updates(f.offset, limit)
	if err != nil {
		return nil, e.Wrap("cannot get events", err)
	}

	if len(updates) == 0 {
		return nil, nil
	}

	res := make([]events.Event, 0, len(updates))

	for _, u := range updates {
		res = append(res, event(u))
	}

	f.offset = updates[len(updates)-1].ID + 1

	return res, nil
}

// event converts a Telegram update to an events.Event.
func event(upd tg_custom_client.Update) events.Event {
	updType := fetchType(upd)

	res := events.Event{
		Type: fetchType(upd),
		Text: fetchText(upd),
	}

	if updType == events.Message {
		res.Meta = tg_processor.Meta{
			ChatID:   upd.Message.Chat.ID,
			Username: upd.Message.From.Username,
		}
	}

	return res
}

// fetchText extracts the text content from a Telegram update.
func fetchText(upd tg_custom_client.Update) string {
	if upd.Message == nil {
		return ""
	}

	return upd.Message.Text
}

// fetchType determines the event type from a Telegram update.
func fetchType(upd tg_custom_client.Update) events.Type {
	if upd.Message == nil {
		return events.Unknown
	}

	return events.Message
}
