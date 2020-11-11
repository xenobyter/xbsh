package view

import (
	"io/ioutil"
	"log"
	"os"
	"testing"
)

func Test_subDirs(t *testing.T) {
	dir, err := ioutil.TempDir("", "xbsh")
	defer os.RemoveAll(dir)
	if err != nil {
		log.Fatalln(err)
	}
	type args struct {
		path string
	}
	tests := []struct {
		name    string
		args    args
		wantRes string
	}{
		{"find d1, d2", args{dir},"d1\nd2"},
		{"don't find file", args{dir+"/d1"},"d11"},
	}

	os.MkdirAll(dir + "/d1/d11/d111", os.ModePerm)
	os.MkdirAll(dir + "/d2/d21/d211", os.ModePerm)
	os.Create(dir + "d1/file")

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotRes := subDirs(tt.args.path); gotRes != tt.wantRes {
				t.Errorf("subDirs() = %v, want %v", gotRes, tt.wantRes)
			}
		})
	}
}
