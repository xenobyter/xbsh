package view

import (
	"io/ioutil"
	"log"
	"os"
	"reflect"
	"testing"
)

func TestReadDir(t *testing.T) {
	//setup
	tDir, err := ioutil.TempDir("", "xbsh")
	if err != nil {
		log.Fatal(err)
	}
	os.Mkdir(tDir+"/dir1", os.ModePerm)
	d1I, _ := os.Stat(tDir + "/dir1")
	os.Mkdir(tDir+"/dir2", os.ModePerm)
	d2I, _ := os.Stat(tDir + "/dir2")
	os.Create(tDir + "/dir1file")
	f1I, _ := os.Stat(tDir + "/dir1file")
	os.Create(tDir + "/nodirfile")
	f2I, _ := os.Stat(tDir + "/nodirfile")
	defer os.RemoveAll(tDir)

	type args struct {
		srch string
		dir  string
	}
	tests := []struct {
		name            string
		args            args
		wantCompletions []os.FileInfo
	}{
		{"return nil on wrong dir", args{"dir", tDir + "/dir3"}, nil},
		{"return nodirfile", args{"nodir", tDir}, []os.FileInfo{f2I}},
		{"return dir2", args{"dir2", tDir}, []os.FileInfo{d2I}},
		{"return dir1 and dir1file", args{"dir1", tDir}, []os.FileInfo{d1I, f1I}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotCompletions := completionSearch(tt.args.srch, tt.args.dir); !reflect.DeepEqual(gotCompletions, tt.wantCompletions) {
				t.Errorf("ReadDir() = %v, want %v", gotCompletions, tt.wantCompletions)
			}
		})
	}
}
