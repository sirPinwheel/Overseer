package handlers

import (
	"fmt"

	"github.com/gempir/go-twitch-irc"
)

// Package for all the handlers that
// are called when a command is sent

func testHandler(*twitch.Client, *twitch.PrivateMessage) {
	fmt.Println("TEST COMPLETE")
}

// Handlers - List of exported handlers
var Handlers map[string]func(*twitch.Client, *twitch.PrivateMessage) = map[string]func(*twitch.Client, *twitch.PrivateMessage){
	"test": testHandler,
}
