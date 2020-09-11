package cmd

import (
	"errors"
	"fmt"
	"testing"
)

func TestParseCmd( t *testing.T) {
	fmt.Println("parseCmd should return empty args")
	wantCmd := "ls"
	var wantArgs []string 
	gotCmd, gotArgs, err := parseCmd("ls")
	if gotCmd != wantCmd {
		t.Errorf("got: %v, want: %v", gotCmd, wantCmd)
	}
	if len(gotArgs) != 0 {
		t.Errorf("got: %v, want: %v", gotArgs, wantArgs)
	}
	
	fmt.Println("parseCmd should handle a single argument")
	wantArgs = append(wantArgs, "/")
	gotCmd, gotArgs, err = parseCmd("ls /")
	if gotCmd != wantCmd {
		t.Errorf("got: %v, want: %v", gotCmd, wantCmd)
	}
	if gotArgs[0] != wantArgs[0] {
		t.Errorf("got: %v, want: %v", gotArgs, wantArgs)
	}
	
	fmt.Println("parseCmd should handle two arguments")
	wantArgs = append([]string{"-la"}, wantArgs...)
	gotCmd, gotArgs, err = parseCmd("ls -la /")
	if gotCmd != wantCmd {
		t.Errorf("got: %v, want: %v", gotCmd, wantCmd)
	}
	if gotArgs[0] != wantArgs[0] {
		t.Errorf("got: %v, want: %v", gotArgs, wantArgs)
	}
	
	fmt.Println("parseCmd should return errNoCommand on empty cmd")
	wErr := errors.New("errNoCommand")
	gotCmd, gotArgs, err = parseCmd("")
	if err.Error() != wErr.Error() {
		t.Errorf("got: %v, want: %v", err.Error(), wErr.Error())
	}

}