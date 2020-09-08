package cmd

import (
	"fmt"
	"os"
	"os/user"
	"strings"
	"testing"
)

func TestNewPrompt(t *testing.T) {
	wantUser, _ := user.Current()
	wantHost, _ := os.Hostname()
	wantDir, _ := os.Getwd()

	prompt := GetPrompt()

	fmt.Println("NewPrompt should return a prompt starting with username")
	gotUser := strings.Split(prompt, "@")[0]
	if gotUser != wantUser.Username {
		t.Errorf("got: %v, want: %v", gotUser, wantUser.Username)
	}

	fmt.Println("NewPrompt should include the hostname")
	gotHost := strings.Split(prompt[len(gotUser)+1:], ":")[0]
	if gotHost != wantHost {
		t.Errorf("got: %v, want: %v", gotHost, wantHost)
	}

	fmt.Println("NewPrompt should include the directory")
	gotDir := strings.Split(prompt[len(gotUser)+len(gotHost)+2:], "$")[0]
	if gotDir != wantDir {
		t.Errorf("got: %v, want: %v", gotDir, wantDir)
	}
}
