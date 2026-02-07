package tg_negasus_fetcher

import (
	"context"
	"errors"
	"go_link_storage/pkg/clients/tg_negasus_client"
	"go_link_storage/pkg/events"
	"go_link_storage/pkg/events/tg_processor"
	"log"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
)

type Fetcher struct {
	tg *tg_negasus_client.Client // Telegram API client
}

var (
	// ErrUnknownEventType is returned when an event type cannot be determined.
	ErrUnknownEventType = errors.New("unknown event type")
	// ErrUnknownMetaType is returned when event metadata has an unexpected type.
	ErrUnknownMetaType = errors.New("unknown meta type")
)

// New creates a new Telegram event processor with the given client.
func New(client *tg_negasus_client.Client) *Fetcher {
	return &Fetcher{
		tg: client,
	}
}

func (f *Fetcher) Start(handleEventsCallback func(events []events.Event) error, batchSize int) {
	opts := []bot.Option{
		bot.WithDefaultHandler(func(ctx context.Context, b *bot.Bot, update *models.Update) {
			res := make([]events.Event, 0, 1)
			res = append(res, event(update))
			if err := handleEventsCallback(res); err != nil {
				log.Println("error handling event:", err)
			}
		}),
	}

	b, err := bot.New(f.tg.Token, opts...)
	if err != nil {
		panic(err)
	}

	f.tg.Bot = b
	f.tg.Bot.Start(f.tg.Ctx)
}

// event converts a Telegram update to an events.Event.
func event(upd *models.Update) events.Event {
	updType := fetchType(upd)

	res := events.Event{
		Type: fetchType(upd),
		Text: fetchText(upd),
	}

	if updType == events.Message {
		res.Meta = tg_processor.Meta{
			ChatID:   int(upd.Message.Chat.ID),
			Username: upd.Message.From.Username,
		}
	}

	return res
}

// fetchText extracts the text content from a Telegram update.
func fetchText(upd *models.Update) string {
	if upd.Message == nil {
		return ""
	}

	return upd.Message.Text
}

// fetchType determines the event type from a Telegram update.
func fetchType(upd *models.Update) events.Type {
	if upd.Message == nil {
		return events.Unknown
	}

	return events.Message
}
