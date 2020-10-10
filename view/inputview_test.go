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
		{"line 1", args{10,20,3},10,0,0},
		{"line 2", args{21,20,3},1,1,0},
		{"line 4", args{81,20,3},1,2,2},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotCx, gotCy, gotOy := caclulateCursor(tt.args.l, tt.args.mx, tt.args.my)
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
