package main // version 0.0.2

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/gempir/go-twitch-irc"
	"github.com/sirpinwheel/overseer/handlers"
	"github.com/sirpinwheel/overseer/settings"
)

// Exported global constants
const (
	CHANNEL string = settings.CHANNEL
	BOT     string = settings.BOT
	OAUTH   string = settings.OAUTH
	PREFIX  string = settings.PREFIX
)

// BotClient - exportin connection
var BotClient *twitch.Client = twitch.NewClient(BOT, OAUTH)

// Function for halting the bot safely
func stop() {
	err := BotClient.Disconnect()

	if err != nil {
		panic(err)
	}
}

func main() {
	// string -> function map for commands called locally in console
	consoleHandlerMap := map[string]func(string){
		"stop": func(arguments string) {
			stop()
		},

		"say": func(arguments string) {
			BotClient.Say(CHANNEL, arguments)
		},
	}

	// string -> function map for commands called in chat by owner
	adminHandlerMap := map[string]func(*twitch.PrivateMessage){
		"stop": func(msg *twitch.PrivateMessage) {
			stop()
		},
	}

	// Hook / callback for general message type sent in chat
	BotClient.OnPrivateMessage(func(message twitch.PrivateMessage) {
		// Check if message is not empty
		if len(message.Message) != 0 {
			if message.User.Name == CHANNEL {
				for k, v := range adminHandlerMap {
					if k == strings.TrimPrefix(message.Message, PREFIX) {
						v(&message)
					}
				}
			}

			// Check if message begins with prefix (a.k.a. is a command)
			if strings.HasPrefix(message.Message, PREFIX) {
				for k, v := range handlers.Handlers {
					if k == strings.TrimPrefix(message.Message, PREFIX) {
						v(BotClient, &message)
					}
				}
			}
		}
	})

	// Greeting
	fmt.Println("Connected to #" + CHANNEL + " as " + BOT)
	fmt.Println("- - - - - - - - - - - - - - - - - - - - - - -")

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

	// Joining channel
	BotClient.Join(CHANNEL)
	err := BotClient.Connect()
	if err != nil {
		if !strings.Contains(err.Error(), "client called Disconnect()") {
			panic(err)
		}
	}
}
