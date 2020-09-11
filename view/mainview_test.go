package view

import (
	"fmt"
	"testing"
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
	got = trimLine([]string{"user@host:/home$" + "\u202f" + "run"})
	if want != got {
		t.Errorf("got: %v, want: %v", got, want)
	}
}

func TestCountLines(t *testing.T) {
	fmt.Println("countLines should return 0 on empty input")
	want := 0
	got := countLines([]byte{})
	if want != got {
		t.Errorf("got: %v, want: %v", got, want)
	}

	fmt.Println("countLines should return 0 on input without Newlines")
	want = 0
	got = countLines([]byte{80,80,80,80})
	if want != got {
		t.Errorf("got: %v, want: %v", got, want)
	}

	fmt.Println("countLines should return 1 on input with 1 newline")
	want = 1
	got = countLines([]byte{80,80,'\n',80})
	if want != got {
		t.Errorf("got: %v, want: %v", got, want)
	}

	fmt.Println("countLines should count newlines in two arguments")
	want = 3
	got = countLines([]byte{80,80,'\n',80},[]byte{80,'\n','\n',80})
	if want != got {
		t.Errorf("got: %v, want: %v", got, want)
	}
}
