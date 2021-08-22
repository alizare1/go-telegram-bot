// Package telegrambot is a simple wrapper for telegram bot API written in Go
package telegrambot

import (
	"encoding/json"
	"errors"
	"fmt"
)

type Bot struct {
	token            string
	lastUpdateOffset int64
	handlers         []Handler
}

type User struct {
	Id                       int64  `json:"id"`
	IsBot                    bool   `json:"is_bot"`
	FirstName                string `json:"first_name"`
	LastName                 string `json:"last_name"`
	Username                 string `json:"username"`
	LanguageCode             string `json:"language_code"`
	CanJoinGroups            bool   `json:"can_join_groups"`
	CanReadAllGroupdMessages bool   `json:"can_read_all_group_messages"`
	SupportsInlineQueries    bool   `json:"supports_inline_queries"`
}

type Chat struct {
	Id        int64  `json:"id"`
	Type      string `json:"type"`
	Title     string `json:"title"`
	Username  string `json:"username"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
}

type Message struct {
	MessageId int32  `json:"message_id"`
	Chat      Chat   `json:"chat"`
	From      User   `json:"from"`
	Text      string `json:"text"`
}

type Update struct {
	UpdateId int64   `json:"update_id"`
	Message  Message `json:"message"`
}

func NewBot(token string) Bot {
	return Bot{token: token}
}

func (bot *Bot) AddCommandHandler(cmd string, handler func(*Bot, *Message)) {
	newHandler := &CommandHandler{cmd, handler}
	bot.handlers = append(bot.handlers, newHandler)
}

func (bot *Bot) AddTextMessageHandler(handler func(*Bot, *Message)) {
	newHandler := &TextMessageHandler{handler}
	bot.handlers = append(bot.handlers, newHandler)
}

func (bot *Bot) GetMe() (username string, err error) {
	getMeUrl := getMethodUrl("getMe", bot.token)
	tgResp, err := sendGetRequest(getMeUrl, nil)
	if err != nil {
		return "", err
	}

	if tgResp.Ok {
		var botUser User
		json.Unmarshal(tgResp.Result, &botUser)
		return botUser.Username, nil
	}

	return "", errors.New(tgResp.Description)
}

func (bot *Bot) getUpdates(timeout int) ([]Update, error) {
	reqUrl := getMethodUrl("getUpdates", bot.token)
	params := map[string]string{
		"timeout": fmt.Sprint(timeout),
		"offset":  fmt.Sprint(bot.lastUpdateOffset),
	}
	tgResp, err := sendGetRequest(reqUrl, params)
	if err != nil {
		return nil, err
	}

	if tgResp.Ok {
		var updates []Update
		json.Unmarshal(tgResp.Result, &updates)
		if len(updates) > 0 {
			bot.lastUpdateOffset = updates[len(updates)-1].UpdateId + 1
		}
		return updates, nil
	}

	return nil, errors.New(tgResp.Description)
}

func (bot *Bot) SendMessage(chatId int64, text string) (message Message, err error) {
	sendMessageUrl := getMethodUrl("sendMessage", bot.token)
	body := map[string]interface{}{
		"chat_id": chatId,
		"text":    text,
	}
	tgResp, err := sendPostRequest(sendMessageUrl, body)
	if err != nil {
		return message, err
	}

	if tgResp.Ok {
		json.Unmarshal(tgResp.Result, &message)
		return message, nil
	}
	return message, errors.New(tgResp.Description)
}

func (bot *Bot) ForwardMessage(chatId int64, msg *Message) (message Message, err error) {
	url := getMethodUrl("forwardMessage", bot.token)
	params := map[string]interface{}{
		"chat_id":      fmt.Sprint(chatId),
		"from_chat_id": fmt.Sprint(msg.From.Id),
		"message_id":   fmt.Sprint(msg.MessageId),
	}

	tgResp, err := sendPostRequest(url, params)
	if err != nil {
		return message, err
	}

	if tgResp.Ok {
		json.Unmarshal(tgResp.Result, &message)
		return message, nil
	}
	return message, errors.New(tgResp.Description)
}

// handleUpdate finds the first matching handler for the new update and calls its handle() method
func (bot *Bot) handleUpdate(update Update) {
	for _, handler := range bot.handlers {
		if handler.matches(&update.Message) {
			handler.handle(bot, &update.Message)
			return
		}
	}
}

func (bot *Bot) handleUpdates(updates []Update) {
	for _, u := range updates {
		jobQueue <- job{Update: u, Bot: bot}
	}
}

// StartPolling does long polling with specified timeout.
// timeout is a positive integer and defaults to 0 eg. short polling. short polling is only used for testing
func (bot *Bot) StartPolling(timeoutOption ...int) {
	var timeout int
	if len(timeoutOption) > 0 {
		timeout = timeoutOption[0]
	}

	initDispatcher()

	for {
		updates, _ := bot.getUpdates(timeout)
		bot.handleUpdates(updates)
	}
}
