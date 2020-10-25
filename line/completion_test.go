package line

import (
	"os"
	"reflect"
	"testing"
)

func TestLastItem(t *testing.T) {
	wd, _ := os.Getwd()
	tests := []struct {
		name     string
		buffer   string
		wantCmd  string
		wantPath string
		wantBin  bool
	}{
		{"Empty buffer", "", "", wd, true},
		{"Simple command", "cmd", "cmd", wd, true},
		{"Command with path", "path/cmd", "cmd", wd + "/path", false},
		{"Command with ./", "./path/cmd", "cmd", wd + "/path", false},
		{"Command with absolute path /", "/path/cmd", "cmd", "/path", false},
		{"Absolute path only", "/cmd", "cmd", "/", false},
		{"Two commands", "/path/cmd1 cmd2", "cmd2", wd, false},
		{"Trim trainling \n", "ls\n", "ls", wd, true},
		{"Handle one char item", "x", "x", wd, true},
		{"More than one blank", "  /fo  ls", "ls", wd, false},
		{"Trailing slash", "dir/", "", wd + "/dir", false},
		{"Empty second cmd", "cmd ", "", wd, false},
		{"Empty first cmd", " ", "", wd, true},
		{"Second cmd with absolute path", "cmd /abs", "abs", "/", false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotCmd, gotPath, gotBin := lastItem(tt.buffer)
			if gotCmd != tt.wantCmd {
				t.Errorf("LastItem() gotCmd = %v, want %v", gotCmd, tt.wantCmd)
			}
			if gotPath != tt.wantPath {
				t.Errorf("LastItem() gotPath = %v, want %v", gotPath, tt.wantPath)
			}
			if gotBin != tt.wantBin {
				t.Errorf("LastItem() gotBin = %v, want %v", gotBin, tt.wantBin)
			}
		})
	}
}

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
