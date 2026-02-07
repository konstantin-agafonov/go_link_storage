package tg_custom_client

// BaseResponse represents the base structure of Telegram API responses.
type BaseResponse struct {
	Ok bool `json:"ok"` // Indicates if the API request was successful
}

// UpdatesResponse represents the response from the getUpdates API method.
type UpdatesResponse struct {
	BaseResponse
	Result []Update `json:"result"` // List of updates
}

// From represents the sender information in a Telegram message.
type From struct {
	Username string `json:"username"` // Telegram username
}

// Chat represents a Telegram chat.
type Chat struct {
	ID int `json:"id"` // Chat ID
}

// IncomingMessage represents an incoming Telegram message.
type IncomingMessage struct {
	Text string `json:"text"` // Message text content
	From From   `json:"from"` // Sender information
	Chat Chat   `json:"chat"` // Chat information
}

// Update represents a Telegram update from the Bot API.
type Update struct {
	ID      int              `json:"update_id"` // Unique update identifier
	Message *IncomingMessage `json:"message"`   // Incoming message (nil if not a message update)
}
