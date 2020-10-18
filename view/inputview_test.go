package view

import (
	"testing"

	"github.com/xenobyter/xbsh/cmd"
)

func Test_trimLine(t *testing.T) {
	tests := []struct {
		name        string
		bufferLines []string
		want        string
	}{
		{"should delete \\n at the end", []string{"test\n"}, "test"},
		{"shouldn't delete another char at the end", []string{"run"}, "run"},
		{"should cut off the prompt", []string{cmd.GetPrompt() + "run"}, "run"},
		{"should trim leading blanks", []string{cmd.GetPrompt() + "  run"}, "run"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := trimLine(tt.bufferLines); got != tt.want {
				t.Errorf("trimLine() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_caclulateCursor(t *testing.T) {
	type args struct {
		l  int
		mx int
		my int
	}
	tests := []struct {
		name   string
		args   args
		wantCx int
		wantCy int
		wantOy int
	}{
		{"line 1", args{10, 20, 3}, 10, 0, 0},
		{"line 2", args{21, 20, 3}, 1, 1, 0},
		{"line 4", args{81, 20, 3}, 1, 2, 2},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotCx, gotCy, gotOy := calculateCursor(tt.args.l, tt.args.mx, tt.args.my)
			if gotCx != tt.wantCx {
				t.Errorf("caclulateCursor() gotCx = %v, want %v", gotCx, tt.wantCx)
			}
			if gotCy != tt.wantCy {
				t.Errorf("caclulateCursor() gotCy = %v, want %v", gotCy, tt.wantCy)
			}
			if gotOy != tt.wantOy {
				t.Errorf("caclulateCursor() gotOy = %v, want %v", gotOy, tt.wantOy)
			}
		})
	}
}

func Test_splitCmd(t *testing.T) {
	tests := []struct {
		name     string
		cmd      string
		wantItem string
		wantPath string
	}{
		{"find path and first word", "user@host:/tmp$ ls", "ls", "/tmp"},
		{"trim trainling \n", "user@host:/tmp$ ls\n", "ls", "/tmp"},
		{"handle absolut paths", "user@host:/tmp$ /wd/ls", "ls", "/wd"},
		{"handle ./", "user@host:/tmp$ ./ls", "ls", "/tmp"},
		{"handle empty string", "", "", ""},
		{"handle prompt only", "user@host:/tmp/tmp$ ", "", "/tmp/tmp"},
		{"handle one char item", "user@host:/tmp/tmp$ x", "x", "/tmp/tmp"},
		{"find deeper dir with ./", "user@host:/tmp/tmp$ ./dir/fo", "fo", "/tmp/tmp/dir"},
		{"find deeper dir without ./", "user@host:/tmp/tmp$ dir/fo", "fo", "/tmp/tmp/dir"},
		{"find absolute paths starting from root", "user@host:/tmp/tmp$ /fo", "fo", "/"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotItem, gotPath := splitCmd(tt.cmd)
			if gotItem != tt.wantItem {
				t.Errorf("splitCmd(%v) gotItem = %v, want %v", tt.cmd, gotItem, tt.wantItem)
			}
			if gotPath != tt.wantPath {
				t.Errorf("splitCmd(%v) gotPath = %v, want %v", tt.cmd, gotPath, tt.wantPath)
			}
		})
	}
}

func Test_findLastSep(t *testing.T) {
	type args struct {
		str string
		sep []string
	}
	tests := []struct {
		name    string
		args    args
		wantMax int
	}{
		{"find last /", args{"test/test/", []string{"/"}}, 9},
		{"find last \\", args{"test/test\\", []string{"/", "\\"}}, 9},
		{"find last \\ with following other chars", args{"test/test\\test", []string{"/", "\\"}}, 9},
		{"return -1 without sep", args{"test", []string{"/", "\\"}}, -1},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotMax := findLastSep(tt.args.str, tt.args.sep); gotMax != tt.wantMax {
				t.Errorf("findLastSep() = %v, want %v", gotMax, tt.wantMax)
			}
		})
	}
}
