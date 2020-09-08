package cmd

import (
	"os"
	"os/user"
)

// GetPrompt returns the Prompt
func GetPrompt() string {
	user, _ := user.Current()
	dir, _ := os.Getwd()
	host, _ := os.Hostname()

	return user.Username + "@" + host + ":" + dir + "$"
}
