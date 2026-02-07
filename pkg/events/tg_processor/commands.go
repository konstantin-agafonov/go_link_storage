package tg_processor

import (
	"context"
	"errors"
	"go_link_storage/pkg/events"
	"go_link_storage/pkg/lib/e"
	"go_link_storage/pkg/storage"
	"log"
	"net/url"
	"strings"
)

const (
	RndCmd   = "/rnd"   // Command to get a random saved page
	HelpCmd  = "/help"  // Command to show help message
	StartCmd = "/start" // Command to start the bot
)

// doCmd processes a command or URL from a user message.
func (p *Processor) doCmd(
	text string,
	chatID int,
	username string) error {

	text = strings.TrimSpace(text)

	log.Printf("new command '%s' from '%s'", text, username)

	if isAddCmd(text) {
		return p.savePage(chatID, text, username)
	}

	switch text {
	case RndCmd:
		return p.sendRandom(chatID, username)
	case HelpCmd:
		return p.sendHelp(chatID)
	case StartCmd:
		return p.sendHello(chatID)
	default:
		return p.tg.SendMessage(chatID, msgUnknownCommand)
	}
}

// savePage saves a page URL to storage for the given user.
func (p *Processor) savePage(
	chatID int,
	pageURL string,
	username string) (err error) {

	defer func() {
		err = e.WrapIfErr("cannot process command: save page", err)
	}()

	sendMsg := NewMessageSender(chatID, p.tg)

	page := &storage.Page{
		URL:      pageURL,
		UserName: username,
	}

	exists, err := p.storage.Exists(context.Background(), page)
	if err != nil {
		return err
	}
	if exists {
		return sendMsg(msgAlreadyExists)
	}

	if err := p.storage.Save(context.Background(), page); err != nil {
		return err
	}

	if err := sendMsg(msgSaved); err != nil {
		return err
	}

	return nil
}

// sendRandom retrieves and sends a random saved page to the user, then removes it.
func (p *Processor) sendRandom(
	chatID int,
	username string) (err error) {

	defer func() { err = e.WrapIfErr("cannot do command: send random", err) }()

	sendMsg := NewMessageSender(chatID, p.tg)

	page, err := p.storage.PickRandom(context.Background(), username)
	if err != nil && !errors.Is(err, storage.ErrNoSavedPages) {
		return err
	}
	if errors.Is(err, storage.ErrNoSavedPages) {
		return sendMsg(msgNoSavedPages)
	}

	if err := sendMsg(page.URL); err != nil {
		return err
	}

	return p.storage.Remove(context.Background(), page)
}

// sendHelp sends the help message to the user.
func (p *Processor) sendHelp(chatID int) error {
	return p.tg.SendMessage(chatID, msgHelp)
}

// sendHello sends the welcome message to the user.
func (p *Processor) sendHello(chatID int) error {
	return p.tg.SendMessage(chatID, msgHello)
}

// NewMessageSender creates a closure function for sending messages to a specific chat.
func NewMessageSender(
	chatID int,
	tg events.Client) func(string) error {

	return func(msg string) error {
		return tg.SendMessage(chatID, msg)
	}
}

// isAddCmd checks if the text is a command to add a page (i.e., a URL).
func isAddCmd(text string) bool {
	return isURL(text)
}

// isURL validates if the given text is a valid URL.
func isURL(text string) bool {
	u, err := url.Parse(text)

	return err == nil && u.Host != ""
}
