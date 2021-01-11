package view

import (
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"reflect"
	"testing"
	"time"
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
		{"Insert", args{dir, []string{"ins suffix suf"}}, []string{"file1  => file1suffix", "files2 => files2suffix"}},
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
		cnt   int
	}
	tests := []struct {
		name    string
		args    args
		wantRes string
	}{
		{"No rules", args{"name", []string{}, 0}, "name"},
		{"Wrong rule", args{"name", []string{""}, 0}, "name"},
		{"Insert prefix", args{"name", []string{"ins string pre"}, 0}, "stringname"},
		{"Insert missing place", args{"name", []string{"ins"}, 0}, "name"},
		{"Insert prefix missing string", args{"name", []string{"ins pre"}, 0}, "name"},
		{"Insert suffix", args{"name", []string{"ins string suf"}, 0}, "namestring"},
		{"Insert position", args{"name", []string{"ins string pos 2"}, 0}, "nastringme"},
		{"Insert after", args{"name", []string{"ins string aft nam"}, 0}, "namstringe"},
		{"Insert count", args{"name", []string{"ins 00 pre"}, 0}, "00name"},
		{"Dont Insert after inc", args{"name", []string{"inc test", "ins test pre"}, 0}, "name"},
		{"Insert after inc", args{"name", []string{"inc .am.", "ins test pre"}, 0}, "testname"},
		{"Dont Insert after exc", args{"name", []string{"exc .am.", "ins test pre"}, 0}, "name"},
		{"Uppercase after ins", args{"name", []string{"ins test pre", "cas upp"}, 0}, "Testname"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotRes := doRules(tt.args.name, tt.args.rules, tt.args.cnt, time.Now(), false)
			if gotRes != tt.wantRes {
				t.Errorf("doRules() = %v, want %v", gotRes, tt.wantRes)
			}
		})
	}
}

func Test_ins(t *testing.T) {
	type args struct {
		place  string
		in     string
		fields []string
		cnt    int
	}
	tests := []struct {
		name    string
		args    args
		wantOut string
	}{
		{"Prefix", args{"pre", "name", []string{"ins", "string", "pre"}, 0}, "stringname"},
		{"Suffix", args{"suf", "name", []string{"ins", "string", "suf"}, 0}, "namestring"},
		{"Position", args{"pos", "name", []string{"ins", "string", "pos", "2"}, 0}, "nastringme"},
		{"Position missing pos", args{"pos", "name", []string{"ins", "string", "pos"}, 0}, "name"},
		{"Position missing string", args{"pos", "name", []string{"ins", "pos", "2"}, 0}, "name"},
		{"Position greater string", args{"pos", "name", []string{"ins", "string", "pos", "5"}, 0}, "name"},
		{"Position negative", args{"pos", "name", []string{"ins", "string", "pos", "-1"}, 0}, "name"},
		{"After", args{"aft", "name", []string{"ins", "string", "aft", "am"}, 0}, "namstringe"},
		{"After with substring not found", args{"aft", "name", []string{"ins", "string", "aft", "xx"}, 0}, "name"},
		{"After with string missing", args{"aft", "name", []string{"ins", "aft", "xx"}, 0}, "name"},
		{"After without arguments", args{"aft", "name", []string{"ins", "aft"}, 0}, "name"},
		{"Wrong place", args{"xxx", "name", []string{"ins", "xxx"}, 0}, "name"},
		{"count", args{"pre", "name", []string{"ins", "000", "pre"}, 0}, "000name"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotOut := ins(tt.args.place, tt.args.in, tt.args.fields, tt.args.cnt)
			if gotOut != tt.wantOut {
				t.Errorf("place() = %v, want %v", gotOut, tt.wantOut)
			}
		})
	}
}

func Test_del(t *testing.T) {
	type args struct {
		name   string
		fields []string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{"no args", args{"name", []string{}}, "name"},
		{"del string", args{"name", []string{"del", "me"}}, "na"},
		{"del string at start", args{"name", []string{"del", "na"}}, "me"},
		{"del prefix", args{"name", []string{"del", "n", "pre"}}, "ame"},
		{"del suffix", args{"name", []string{"del", "e", "suf"}}, "nam"},
		{"del any substring", args{"namename", []string{"del", "am", "any"}}, "nene"},
		{"del any no find", args{"namename", []string{"del", "xy", "any"}}, "namename"},
		{"del from pos 3", args{"namename", []string{"del", "3"}}, "na"},
		{"del from pos 1", args{"namename", []string{"del", "1"}}, ""},
		{"del from pos 2", args{".git", []string{"del", "5", "6"}}, ".git"},
		{"del from pos 0", args{"namename", []string{"del", "0"}}, "namename"},
		{"del from pos 2 to 3", args{"namename", []string{"del", "2", "3"}}, "nename"},
		{"del last 3", args{"namename", []string{"del", "-3"}}, "namen"},
		{"del last n with big n", args{"namename", []string{"del", "-10"}}, "namename"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := del(tt.args.name, tt.args.fields); got != tt.want {
				t.Errorf("del() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_rep(t *testing.T) {
	type args struct {
		name   string
		fields []string
		cnt    int
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{"no args", args{"", []string{}, 0}, ""},
		{"First occurence", args{"name", []string{"rep", "am", "yy"}, 0}, "nyye"},
		{"First occurence, no find", args{"name", []string{"rep", "xx", "yy"}, 0}, "name"},
		{"Prefix", args{"name", []string{"rep", "na", "yy", "pre"}, 0}, "yyme"},
		{"Suffix", args{"name", []string{"rep", "me", "yy", "suf"}, 0}, "nayy"},
		{"Any", args{"namename", []string{"rep", "am", "yy", "any"}, 0}, "nyyenyye"},
		{"Incrementing", args{"namename", []string{"rep", "am", "001", "any"}, 1}, "n002en002e"},
		{"Incrementing first", args{"namename", []string{"rep", "am", "001"}, 1}, "n002ename"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := rep(tt.args.name, tt.args.fields, tt.args.cnt); got != tt.want {
				t.Errorf("rep() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_cas(t *testing.T) {
	type args struct {
		name   string
		fields []string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{"Empty args", args{"", []string{}}, ""},
		{"First letter upp", args{"name", []string{"cas", "upp"}}, "Name"},
		{"First letter upp umlauts", args{"äame", []string{"cas", "upp"}}, "Äame"},
		{"First letter low", args{"Name", []string{"cas", "low"}}, "name"},
		{"First letter of first word only upp", args{"name name", []string{"cas", "upp"}}, "Name name"},
		{"First letter of first word only low", args{"NAME NAME", []string{"cas", "low"}}, "nAME NAME"},
		{"Empty name", args{"", []string{"cas", "upp"}}, ""},
		{"Name with one char upp", args{"x", []string{"cas", "upp"}}, "X"},
		{"Name with one char low", args{"X", []string{"cas", "low"}}, "x"},
		{"Every word upp", args{"name name", []string{"cas", "upp", "wrd"}}, "Name Name"},
		{"Every word low", args{"Name name", []string{"cas", "low", "wrd"}}, "name name"},
		{"Every word low umlauts", args{"Äame Äame", []string{"cas", "low", "wrd"}}, "äame äame"},
		{"Every word low one char word", args{"Ä e", []string{"cas", "low", "wrd"}}, "ä e"},
		{"Any char upp", args{"name name", []string{"cas", "upp", "any"}}, "NAME NAME"},
		{"Any char low", args{"NAME NAME", []string{"cas", "low", "any"}}, "name name"},
		{"Any m upp", args{"name name", []string{"cas", "upp", "any", "m"}}, "naMe naMe"},
		{"Any m low", args{"NAME NAME", []string{"cas", "low", "any", "M"}}, "NAmE NAmE"},
		{"Any with empty char", args{"name name", []string{"cas", "upp", "any", ""}}, "name name"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := cas(tt.args.name, tt.args.fields); got != tt.want {
				t.Errorf("cas() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_dat(t *testing.T) {
	type args struct {
		name   string
		fields []string
		t      time.Time
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{"Empty args", args{"", []string{}, time.Now()}, ""},
		{"dat only", args{"name", []string{"dat"}, time.Unix(0, 0)}, "1970-01-01 01:00:00"},
		{"dat pre", args{"name", []string{"dat", "pre"}, time.Unix(0, 0)}, "1970-01-01 01:00:00 name"},
		{"dat suf", args{"name", []string{"dat", "suf"}, time.Unix(0, 0)}, "name 1970-01-01 01:00:00"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := dat(tt.args.name, tt.args.fields, tt.args.t); got != tt.want {
				t.Errorf("dat() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_mod(t *testing.T) {
	type args struct {
		fields []string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{"Empty args", args{[]string{}}, "a"},
		{"mod file", args{[]string{"mod", "fil"}}, "f"},
		{"mod directory", args{[]string{"mod", "dir"}}, "d"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := mod(tt.args.fields); got != tt.want {
				t.Errorf("mode() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_rename(t *testing.T) {
	type args struct {
		dir   string
		rules []string
	}
	dir, err := ioutil.TempDir("", "xbsh")
	if err != nil {
		log.Fatalln(err)
	}
	defer os.RemoveAll(dir)
	os.Create(filepath.Join(dir, "file1"))
	os.Create(filepath.Join(dir, "file2"))
	tests := []struct {
		name    string
		args    args
		wantOut []string
	}{
		{"Simple insert", args{dir, []string{"ins pre pre"}}, []string{"prefile1", "prefile2"}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotOut := rename(tt.args.dir, tt.args.rules); !reflect.DeepEqual(gotOut, tt.wantOut) {
				t.Errorf("rename() = %v, want %v", gotOut, tt.wantOut)
			}
			for i, want := range tt.wantOut {
				log.Println(filepath.Join(dir, want))
				if _, err := os.Stat(filepath.Join(dir, want)); os.IsNotExist(err) {
					t.Errorf("rename() failed for %v", tt.wantOut[i])
				}
			}
		})
	}
}
