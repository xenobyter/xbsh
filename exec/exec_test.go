package exec

import (
	"io/ioutil"
	"log"
	"os"
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
		{"Command only", "c", "c", []string{}, false},
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
