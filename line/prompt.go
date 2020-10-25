package line

import (
	"os"
	"os/user"
)

// Prompt returns the Prompt
func Prompt() string {
	user, _ := user.Current()
	dir, _ := os.Getwd()
	host, _ := os.Hostname()

	return user.Username + "@" + host + ":" + dir + "$ "
}
