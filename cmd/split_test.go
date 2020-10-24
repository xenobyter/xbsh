package cmd

import (
	"reflect"
	"testing"
)

func Test_splitBuffer(t *testing.T) {
	tests := []struct {
		name   string
		buffer string
		sep    []rune
		want   []string
	}{
		{"simple split", "item1|item2", groupSep, []string{"item1", "item2"}},
		{"empty string", "", groupSep, []string{}},
		{"no sep found", "oneString", groupSep, []string{"oneString"}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := splitBuffer(tt.buffer, tt.sep); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("spitBuffer() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_seperatePrompt(t *testing.T) {
	tests := []struct {
		name     string
		buffer   string
		wantPath string
		wantCmd  string
	}{
		{"Empty string", "", "", ""},
		{"Usual prompt", "user@host:/workdir$ ", "/workdir", ""},
		{"Prompt with cmd", "user@host:/workdir$ /cmd", "/workdir", "/cmd"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotPath, gotCmd := seperatePrompt(tt.buffer)
			if gotPath != tt.wantPath {
				t.Errorf("seperatePrompt() gotPath = %v, want %v", gotPath, tt.wantPath)
			}
			if gotCmd != tt.wantCmd {
				t.Errorf("seperatePrompt() gotCmd = %v, want %v", gotCmd, tt.wantCmd)
			}
		})
	}
}

func TestLastItem(t *testing.T) {
	tests := []struct {
		name            string
		buffer          string
		wantCmd         string
		wantPath        string
		wantBin bool
	}{
		{"Empty buffer", "", "", "", true},
		{"Prompt only", "user@host:/wd$ ", "", "/wd", true},
		{"Prompt and simple command", "user@host:/wd$ cmd", "cmd", "/wd", true},
		{"Command with path", "user@host:/wd$ path/cmd", "cmd", "/wd/path", false},
		{"Command with ./", "user@host:/wd$ ./path/cmd", "cmd", "/wd/path", false},
		{"Command with absolute path /", "user@host:/wd$ /path/cmd", "cmd", "/path", false},
		{"Absolute path only", "user@host:/wd$ /cmd", "cmd", "/", false},
		{"Two commands", "user@host:/wd$ /path/cmd1 cmd2", "cmd2", "/wd", false},
		{"Trim trainling \n", "user@host:/tmp$ ls\n", "ls", "/tmp", true},
		{"Handle one char item", "user@host:/tmp/tmp$ x", "x", "/tmp/tmp", true},
		{"More than one blank", "user@host:/tmp/tmp$  /fo  ls", "ls", "/tmp/tmp", false},
		{"Trailing slash", "user@host:/tmp/tmp$ dir/", "", "/tmp/tmp/dir", false},
		{"Empty second cmd", "user@host:/tmp/tmp$ cmd ", "", "/tmp/tmp", false},
		{"Empty first cmd", "user@host:/tmp/tmp$  ", "", "/tmp/tmp", true},
		{"Second cmd with absolute path", "user@host:/tmp/tmp$ cmd /abs", "abs", "/", false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotCmd, gotPath, gotBin := LastItem(tt.buffer)
			if gotCmd != tt.wantCmd {
				t.Errorf("LastItem() gotCmd = %v, want %v", gotCmd, tt.wantCmd)
			}
			if gotPath != tt.wantPath {
				t.Errorf("LastItem() gotPath = %v, want %v", gotPath, tt.wantPath)
			}
			if gotBin != tt.wantBin {
				t.Errorf("LastItem() gotCanBeNonExe = %v, want %v", gotBin, tt.wantBin)
			}
		})
	}
}
