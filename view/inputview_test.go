package view

import (
	"fmt"
	"testing"

	"github.com/xenobyter/xbsh/cmd"
)

func TestTrimLine(t *testing.T) {
	fmt.Println("trimLine should delete \\n at the end")
	want := "test"
	got := trimLine([]string{"test\n"})
	if want != got {
		t.Errorf("got: %v, want: %v", got, want)
	}

	fmt.Println("trimLine shouldn't delete another char at the end")
	want = "run"
	got = trimLine([]string{"run"})
	if want != got {
		t.Errorf("got: %v, want: %v", got, want)
	}

	fmt.Println("trimline should cut off the prompt")
	got = trimLine([]string{cmd.GetPrompt() + "run"})
	// got = trimLine([]string{"user@host:/home$" + "\u202f" + "run"})
	if want != got {
		t.Errorf("got: %v, want: %v", got, want)
	}
}
