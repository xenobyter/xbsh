package view

import (
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"reflect"
	"testing"
)

func Test_preview(t *testing.T) {
	//setup
	dir, err := ioutil.TempDir("", "xbsh")
	if err != nil {
		log.Fatalln(err)
	}
	var files = []string{
		filepath.Join(dir, "file1"),
		filepath.Join(dir, "files2"),
	}
	for _, f := range files {
		os.Create(f)
	}
	defer os.RemoveAll(dir)

	type args struct {
		dir   string
		rules []string
	}
	tests := []struct {
		name       string
		args       args
		wantOutput []string
	}{
		{"No dir, no rules", args{"", []string{}}, nil},
		{"Dir, no rules", args{dir, []string{}}, []string{"file1  => file1", "files2 => files2"}},
		{"Insert", args{dir, []string{"ins suf suffix"}}, []string{"file1  => file1suffix", "files2 => files2suffix"}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotOutput := preview(tt.args.dir, tt.args.rules); !reflect.DeepEqual(gotOutput, tt.wantOutput) {
				t.Errorf("preview() = %v, want %v", gotOutput, tt.wantOutput)
			}
		})
	}
}

func Test_doRules(t *testing.T) {
	type args struct {
		name  string
		rules []string
	}
	tests := []struct {
		name    string
		args    args
		wantRes string
	}{
		{"No rules", args{"name", []string{}}, "name"},
		{"Wrong rule", args{"name", []string{""}}, "name"},
		{"Insert prefix", args{"name", []string{"ins pre string"}}, "stringname"},
		{"Insert missing place", args{"name", []string{"ins"}}, "name"},
		{"Insert prefix missing string", args{"name", []string{"ins pre"}}, "name"},
		{"Insert suffix", args{"name", []string{"ins suf string"}}, "namestring"},
		{"Insert position", args{"name", []string{"ins pos 2 string"}}, "nastringme"},
		{"Insert after", args{"name", []string{"ins aft nam string"}}, "namstringe"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotRes := doRules(tt.args.name, tt.args.rules); gotRes != tt.wantRes {
				t.Errorf("doRules() = %v, want %v", gotRes, tt.wantRes)
			}
		})
	}
}

func Test_place(t *testing.T) {
	type args struct {
		place  string
		in     string
		fields []string
	}
	tests := []struct {
		name    string
		args    args
		wantOut string
	}{
		{"Prefix", args{"pre", "name", []string{"ins", "pre", "string"}}, "stringname"},
		{"Suffix", args{"suf", "name", []string{"ins", "suf", "string"}}, "namestring"},
		{"Position", args{"pos", "name", []string{"ins", "pos", "2", "string"}}, "nastringme"},
		{"Position missing pos", args{"pos", "name", []string{"ins", "pos", "string"}}, "name"},
		{"Position missing string", args{"pos", "name", []string{"ins", "pos", "2"}}, "name"},
		{"Position greater string", args{"pos", "name", []string{"ins", "pos", "5", "string"}}, "name"},
		{"Position negative", args{"pos", "name", []string{"ins", "pos", "-1", "string"}}, "name"},
		{"After", args{"aft", "name", []string{"ins", "aft", "am", "string"}}, "namstringe"},
		{"After with substring not found", args{"aft", "name", []string{"ins", "aft", "xx", "string"}}, "name"},
		{"After with string missing", args{"aft", "name", []string{"ins", "aft", "xx"}}, "name"},
		{"After without arguments", args{"aft", "name", []string{"ins", "aft"}}, "name"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotOut := place(tt.args.place, tt.args.in, tt.args.fields); gotOut != tt.wantOut {
				t.Errorf("place() = %v, want %v", gotOut, tt.wantOut)
			}
		})
	}
}
