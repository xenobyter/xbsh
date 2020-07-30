package view

import (
	"fmt"
	"testing"
)

func TestTrimCmdString( t *testing.T) {
	fmt.Println("trimCmdString should delete \\n at the end")
	want := "test"
	got := trimCmdString("test\n")
	if want != got {
		t.Errorf("got: %v, want: %v", got, want)
	}
	fmt.Println("trimCmdString shouldn't delete another char at the end")
	want = "run"
	got = trimCmdString("run")
	if want != got {
		t.Errorf("got: %v, want: %v", got, want)
	}
}