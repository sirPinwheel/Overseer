package main // version 0.0.3

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/gempir/go-twitch-irc"
	"github.com/sirpinwheel/overseer/handlers"
	"github.com/sirpinwheel/overseer/settings"
	"github.com/sirpinwheel/overseer/utils"
)

// BotClient - exportin connection
var BotClient *twitch.Client = twitch.NewClient(settings.BOT, settings.OAUTH)
var ticker *time.Ticker = time.NewTicker(settings.PERIOD)

// string -> function map for commands called locally in console
var consoleHandlerMap = map[string]func(string){
	"stop": func(arguments string) {
		stop()
	},

	"say": func(arguments string) {
		BotClient.Say(settings.CHANNEL, arguments)
	},
}

// string -> function map for commands called in chat by owner
var adminHandlerMap = map[string]func(*twitch.PrivateMessage){
	"stop": func(msg *twitch.PrivateMessage) {
		stop()
	},
}

// Function for starting needed goroutines
func start() {
	// Greeting
	fmt.Println("Connected to #" + settings.CHANNEL + " as " + settings.BOT)
	fmt.Println("- - - - - - - - - - - - - - - - - - - - - - -")

	// Goroutine for periodic task of giving current viewers a point
	go func() {
		for t := range ticker.C {
			_ = t
			utils.GrantPoint()
		}
	}()

	// Goroutine for handling console input
	go func() {
		reader := bufio.NewReader(os.Stdin)
		for {
			fmt.Print(">> ")

			text, _ := reader.ReadString('\n')
			text = strings.TrimRight(text, "\r\n")
			split := strings.SplitN(text, " ", 2)
			command := split[0]
			arguments := ""

			if len(split) > 1 {
				arguments = split[1]
			}

			for k, v := range consoleHandlerMap {
				if k == command {
					v(arguments)
				}
			}
		}
	}()
}

// Function for halting the bot safely
func stop() {
	boterr := BotClient.Disconnect()
	ticker.Stop()

	if boterr != nil {
		panic(boterr)
	}
}

func main() {
	// Hook / callback for general message type sent in chat
	BotClient.OnPrivateMessage(func(message twitch.PrivateMessage) {
		// Check if message is not empty
		if len(message.Message) != 0 {
			if message.User.Name == settings.CHANNEL {
				for k, v := range adminHandlerMap {
					if k == strings.TrimPrefix(message.Message, settings.PREFIX) {
						v(&message)
					}
				}
			}

			// Check if message begins with prefix (a.k.a. is a command)
			if strings.HasPrefix(message.Message, settings.PREFIX) {
				for k, v := range handlers.Handlers {
					if k == strings.TrimPrefix(message.Message, settings.PREFIX) {
						v(BotClient, &message)
					}
				}
			}
		}
	})

	// Joining channel
	BotClient.Join(settings.CHANNEL)
	BotClient.OnConnect(start)

	fmt.Println("Connecting...")

	err := BotClient.Connect()

	if err != nil {
		if !strings.Contains(err.Error(), "client called Disconnect()") {
			panic(err)
		}
	}
}
