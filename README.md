# go-telegram-bot

A simple [Telegram bot API](https://core.telegram.org/bots/api) wrapper written in Go.

This package was created with the purpose of learning Golang. It may get new updates to have more features and cover other parts of Telegram API.

Currently updates are fetched using [long polling](https://en.wikipedia.org/wiki/Push_technology#Long_polling). A job queue and worker pool is implemented using goroutines and channels to handle updates concurrently. (for more information about this implementation read [Handling 1 Million Requests per Minute with Go](http://marcio.io/2015/07/handling-1-million-requests-per-minute-with-golang/))

## Creating a bot

First you need to create a new bot and obtain a token by contacting [@BotFather](https://t.me/BotFather). 

Then you can check if it's working using this code:

```go
package main

import (
    "fmt"
    tgbot "github.com/alizare1/go-telegram-bot"
)

func main() {
    bot := tgbot.NewBot("TOKEN")

	username, _ := bot.GetMe()
	fmt.Println(username)
}
```

If everything is working you should see your bot's username in the output.

Now you can define handlers for your bot to respond to messages:

```go
package main

import tgbot "github.com/alizare1/go-telegram-bot"

func start(bot *tgbot.Bot, msg *tgbot.Message) {
	bot.SendMessage(msg.Chat.Id, "Hello World from Go!")
}

func main() {
    bot := tgbot.NewBot("TOKEN")

	bot.AddCommandHandler("start", start)
    bot.StartPolling()
}
```

If the bot receives `/start` command, it will respond to the user with `Hello World from Go!`