package line

import (
	"os"
	"os/user"
	"testing"
)

func TestGetPrompt(t *testing.T) {
	wantUser, _ := user.Current()
	wantHost, _ := os.Hostname()
	wantDir, _ := os.Getwd()
	wantPrompt := wantUser.Username + "@" + wantHost + ":" + wantDir + "$ "

	tests := []struct {
		name string
		want string
	}{
		{"Simple Prompt", wantPrompt},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Prompt(); got != tt.want {
				t.Errorf("GetPrompt() = %v, want %v", got, tt.want)
			}
		})
	}
}
