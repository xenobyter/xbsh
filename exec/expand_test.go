package exec

import (
	"io/ioutil"
	"log"
	"os"
	"reflect"
	"testing"
)

func Test_expandArg(t *testing.T) {
	//setup
	dir, err := ioutil.TempDir("", "xbsh")
	if err != nil {
		log.Fatalln(err)
	}
	os.Create(dir + "/file1")
	os.Create(dir + "/file2")
	os.Create(dir + "/filX2")
	os.Chdir(dir)

	type args struct {
		args []string
	}
	tests := []struct {
		name        string
		args        args
		wantExpArgs []string
	}{
		{"Empty args", args{nil}, nil},
		{"absolute path no wildcard", args{[]string{dir + "/file1"}}, []string{dir + "/file1"}},
		{"absolute path with *", args{[]string{dir + "/file*"}}, []string{dir + "/file1", dir + "/file2"}},
		{"absolute path with ?", args{[]string{dir + "/file?"}}, []string{dir + "/file1", dir + "/file2"}},
		{"no path with ?", args{[]string{"file?"}}, []string{"file1", "file2"}},
		{"relative path with ?", args{[]string{"./file?"}}, []string{"file1", "file2"}},
		{"nothing to expand", args{[]string{"a", "b"}}, []string{"a", "b"}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotExpArgs := expandArg(tt.args.args); !reflect.DeepEqual(gotExpArgs, tt.wantExpArgs) {
				t.Errorf("expandArg() = %v, want %v", gotExpArgs, tt.wantExpArgs)
			}
		})
	}
}
