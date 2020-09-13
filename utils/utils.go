package utils

//import "github.com/gempir/go-twitch-irc"
import (
	"github.com/sirpinwheel/overseer/settings"
)

// GrantPoint - Give every current viewer a point
func GrantPoint(users *[]string) {
	if IsLive(settings.CHANNEL) {

	}
}

// IsLive - Checks if passed channel is currently streaming
func IsLive(channel string) bool {
	return true
}
