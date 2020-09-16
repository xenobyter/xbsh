package cmd

import (
	"log"
	"os/user"
	"errors"
	"fmt"
	"testing"
)

func TestParseCmd(t *testing.T) {
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

func TestChangeDir(t *testing.T) {
	usr, err := user.Current()
	if err != nil {
		log.Fatal(err)
	}
	dir := usr.HomeDir
	t.Run("should return homeDir + \\n when called without argument", func(t *testing.T) {
		want := dir + "\n"
		got, err := changeDir([]string{})
		if err != nil {
			t.Errorf("Error: %v", err)
		}
		if string(got) != want {
			t.Errorf("got: %v, want: %v", string(got), want)
		}
	})
	t.Run("should return error on stderr when called with non-existing dir", func(t *testing.T){
		want := "chdir xx: no such file or directory\n"
		_, got := changeDir([]string{"xx"})
		if string(got) != want {
			t.Errorf("got: %v, want: %v", string(got), want)
		}
	})
	t.Run("should return dir", func(t *testing.T){
		want := "/\n"
		got, _ := changeDir([]string{"/"})
		if string(got) != want {
			t.Errorf("got: %v, want: %v", string(got), want)
		}
	})
}
