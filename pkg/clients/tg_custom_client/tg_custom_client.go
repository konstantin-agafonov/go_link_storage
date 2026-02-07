// Package tg_custom_client provides a client for interacting with the Telegram Bot API.
package tg_custom_client

import (
	"encoding/json"
	"go_link_storage/pkg/lib/e"
	"io"
	"net/http"
	"net/url"
	"path"
	"strconv"
)

// Client provides methods for interacting with the Telegram Bot API.
type Client struct {
	host     string      // Telegram API host
	basePath string      // Base path for API requests (includes bot token)
	client   http.Client // HTTP client for making requests
}

const (
	getUpdatesMethod  = "getUpdates"  // API method for getting updates
	sendMessageMethod = "sendMessage" // API method for sending messages
)

// New creates a new Telegram client with the given host and bot token.
func New(host string, token string) *Client {
	return &Client{
		host:     host,
		basePath: newBasePath(token),
		client:   http.Client{},
	}
}

// newBasePath constructs the base path for API requests using the bot token.
func newBasePath(token string) string {
	return "bot" + token
}

// Updates fetches updates from the Telegram Bot API.
// offset specifies the update ID to start from, limit specifies the maximum number of updates.
func (c *Client) Updates(offset int, limit int) (updates []Update, err error) {
	defer func() { err = e.WrapIfErr("can't get updates", err) }()

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

	return res.Result, nil
}

// SendMessage sends a text message to the specified chat.
func (c *Client) SendMessage(chatID int, text string) error {
	q := url.Values{}
	q.Add("chat_id", strconv.Itoa(chatID))
	q.Add("text", text)

	_, err := c.doRequest(sendMessageMethod, q)
	if err != nil {
		return e.Wrap("can't send message", err)
	}

	return nil
}

// doRequest performs an HTTP GET request to the Telegram Bot API.
// method specifies the API method, query contains the request parameters.
func (c *Client) doRequest(method string, query url.Values) (data []byte, err error) {
	const errMsg = "couldn't do request"

	defer func() { err = e.WrapIfErr(errMsg, err) }()

	u := url.URL{
		Scheme: "https",
		Host:   c.host,
		Path:   path.Join(c.basePath, method),
	}

	req, err := http.NewRequest(http.MethodGet, u.String(), nil)
	if err != nil {
		return nil, err
	}

	req.URL.RawQuery = query.Encode()

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}

	defer func() { _ = resp.Body.Close() }()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return body, nil
}
