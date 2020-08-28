package cmd

import (
	"fmt"
	"os"
	"testing"
)

func TestSetInitialWorkDir(t *testing.T) {
	fmt.Println("setInitialWorkDir should set the global var workdir")
	want := "/home"
	os.Chdir(want)
	setInitialWorkDir()
	got := workDir
	if want != got {
		t.Errorf("got: %v, want: %v", got, want)
	}

	fmt.Println("GetWorkDir should return initial workDir on first run")
	want = "/home"
	os.Chdir(want)
	workDir = ""
	got = GetWorkDir()
	if want != got {
		t.Errorf("got: %v, want: %v", got, want)
	}

	fmt.Println("GetWorkDir should return new workDir on second run")
	want = "/"
	workDir = "/"
	got = GetWorkDir()
	if want != got {
		t.Errorf("got: %v, want: %v", got, want)
	}
}
