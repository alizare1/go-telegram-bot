package telegrambot

import (
	"fmt"
	"regexp"
)

// structs implementing this interface can be used as message handlers
// matches is used to know if the handler can handle an update, and then handle can be called for that update's message
type Handler interface {
	matches(msg *Message) bool
	handle(bot *Bot, msg *Message)
}

type CommandHandler struct {
	cmd     string
	handler func(bot *Bot, msg *Message)
}

type TextMessageHandler struct {
	handler func(bot *Bot, msg *Message)
}

func (cmdHandler *CommandHandler) New(cmd string, handler func(*Bot, *Message)) {
	cmdHandler.cmd = cmd
	cmdHandler.handler = handler
}

func (cmdHandler *CommandHandler) matches(msg *Message) bool {
	pattern := fmt.Sprintf(`^\/%v(\s|$)`, cmdHandler.cmd)
	matches, _ := regexp.MatchString(pattern, msg.Text)
	return matches
}

func (cmdHandler *CommandHandler) handle(bot *Bot, msg *Message) {
	cmdHandler.handler(bot, msg)
}

func (txtHandler *TextMessageHandler) matches(msg *Message) bool {
	return msg.Text != ""
}

func (txtHandler *TextMessageHandler) handle(bot *Bot, msg *Message) {
	txtHandler.handler(bot, msg)
}
