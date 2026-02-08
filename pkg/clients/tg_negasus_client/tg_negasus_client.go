package tg_negasus_client

import (
	"context"
	"go_link_storage/pkg/lib/e"
	"os"
	"os/signal"

	"github.com/go-telegram/bot"
)

// Client provides methods for interacting with the Telegram Bot API.
type Client struct {
	Bot   *bot.Bot
	Ctx   context.Context
	Token string
}

// New creates a new Telegram client with the given bot token and options.
func New(token string) *Client {
	ctx, _ := signal.NotifyContext(context.Background(), os.Interrupt)

	return &Client{
		Ctx:   ctx,
		Token: token,
	}
}

// newBasePath constructs the base path for API requests using the bot token.
func newBasePath(token string) string {
	return "bot" + token
}

// SendMessage sends a text message to the specified chat.
func (c *Client) SendMessage(chatID int, text string) error {
	_, err := c.Bot.SendMessage(c.Ctx, &bot.SendMessageParams{
		ChatID: chatID,
		Text:   text,
	})
	if err != nil {
		return e.Wrap("can't send message", err)
	}

	return nil
}
