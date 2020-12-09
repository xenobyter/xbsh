package exec

import (
	"io"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"reflect"
	"testing"
)

func Test_changeDir(t *testing.T) {
	dir, err := ioutil.TempDir("", "xbsh")
	if err != nil {
		log.Fatal(err)
	}
	defer os.RemoveAll(dir)
	home, _ := os.UserHomeDir()

	type args struct {
		args []string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
		wantDir string
	}{
		{"empty argument", args{}, false, home},
		{"cd to tempdir", args{[]string{dir}}, false, dir},
		{"cd to non existing dir", args{[]string{dir + "/non"}}, true, dir},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := changeDir(tt.args.args); (err != nil) != tt.wantErr {
				t.Errorf("changeDir() error = %v, wantErr %v", err, tt.wantErr)
			}
			if gotDir, _ := os.Getwd(); gotDir != tt.wantDir {
				t.Errorf("changeDir() dir = %v, wantDir %v", gotDir, tt.wantDir)
			}
		})
	}
}

func Test_splitArgs(t *testing.T) {
	tests := []struct {
		name        string
		line        string
		wantCommand string
		wantArgs    []string
		wantErr     bool
	}{
		{"Empty line", "", "", nil, true},
		{"Command only", "c", "c", nil, false},
		{"One arg", "c a", "c", []string{"a"}, false},
		{"Two args", "c a b", "c", []string{"a", "b"}, false},
		{"More blanks", " c  a  b ", "c", []string{"a", "b"}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotCommand, gotArgs, err := splitArgs(tt.line)
			if (err != nil) != tt.wantErr {
				t.Errorf("splitArgs() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotCommand != tt.wantCommand {
				t.Errorf("splitArgs() gotCommand = %v, want %v", gotCommand, tt.wantCommand)
			}
			if !reflect.DeepEqual(gotArgs, tt.wantArgs) {
				t.Errorf("splitArgs() gotArgs = %v, want %v", gotArgs, tt.wantArgs)
			}
		})
	}
}

func TestCmd(t *testing.T) {
	tests := []struct {
		name    string
		line    string
		wantErr bool
	}{
		{"Empty line", "", true},
		{"cd", "cd", false},
		{"wrong args", "pwd -b", true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := Cmd(tt.line); (err != nil) != tt.wantErr {
				t.Errorf("Cmd() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_redirect(t *testing.T) {
	//setup
	os.Chdir(os.TempDir())
	file, err := ioutil.TempFile("", "xbsh")
	if err != nil {
		log.Fatalln(err)
	}
	defer os.Remove(file.Name())

	tests := []struct {
		name       string
		args       []string
		wantStdin  *os.File
		wantStdout *os.File
		wantStderr *os.File
		wantArgs   []string
		wantErr    string
	}{
		{"No args", []string{}, os.Stdin, os.Stdout, os.Stderr, nil, ""},
		{"echo to absolute file", []string{"test", ">" + file.Name()}, os.Stdin, file, os.Stderr, []string{"test"}, ""},
		{"echo to relative file", []string{"test", ">" + filepath.Base(file.Name())}, os.Stdin, file, os.Stderr, []string{"test"}, ""},
		{"echo with space after >", []string{"test", ">", filepath.Base(file.Name())}, os.Stdin, file, os.Stderr, []string{"test"}, ""},
		{"nothing after >", []string{"test", ">"}, os.Stdin, os.Stdout, os.Stderr, []string{"test"}, "redirect error"},
		{"append to absolute file", []string{"test", ">>" + file.Name()}, os.Stdin, file, os.Stderr, []string{"test"}, ""},
		{"echo with stderr", []string{"test", ">!" + file.Name()}, os.Stdin, os.Stdout, file, []string{"test"}, ""},
		{"append stderr", []string{"test", ">>!" + file.Name()}, os.Stdin, os.Stdout, file, []string{"test"}, ""},
		{"redirect stdin", []string{"test", "<" + file.Name()}, file, os.Stdout, os.Stderr, []string{"test"}, ""},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			io.WriteString(file, "t")
			gotStdin, gotStdout, gotStderr, gotArgs, gotErr := redirect(tt.args)
			if gotStdin.Name() != tt.wantStdin.Name() {
				t.Errorf("redirect() gotStdin = %v, want %v", gotStdin.Name(), tt.wantStdin.Name())
			}
			if gotStdout.Name() != tt.wantStdout.Name() {
				t.Errorf("redirect() gotStdout = %v, want %v", gotStdout.Name(), tt.wantStdout.Name())
			}
			if gotStderr.Name() != tt.wantStderr.Name() {
				t.Errorf("redirect() gotStderr = %v, want %v", gotStderr.Name(), tt.wantStderr.Name())
			}
			if !reflect.DeepEqual(gotArgs, tt.wantArgs) {
				t.Errorf("redirect() gotArgs = %v, want %v", gotArgs, tt.wantArgs)
			}
			if gotErr != nil && gotErr.Error() != tt.wantErr {
				t.Errorf("redirect() gotErr = %v, want %v", gotErr, tt.wantErr)
			}
		})
	}
}
