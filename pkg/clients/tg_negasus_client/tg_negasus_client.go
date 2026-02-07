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

// Updates fetches updates from the Telegram Bot API.
// offset specifies the update ID to start from, limit specifies the maximum number of updates.
func (c *Client) Updates(offset int, limit int) (updates []Update, err error) {
	/*defer func() { err = e.WrapIfErr("can't get updates", err) }()

	q := url.Values{}
	q.Add("offset", strconv.Itoa(offset))
	q.Add("limit", strconv.Itoa(limit))

	data, err := c.doRequest(getUpdatesMethod, q)
	if err != nil {
		return nil, err
	}

	var res UpdatesResponse

	if err := json.Unmarshal(data, &res); err != nil {
		return nil, err
	}

	return res.Result, nil*/

	return nil, nil
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
