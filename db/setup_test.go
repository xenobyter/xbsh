package db

import (
	"io/ioutil"
	"log"
	"os"
	"os/user"
	"strings"
	"testing"

	_ "github.com/mattn/go-sqlite3"
	"github.com/xenobyter/xbsh/cfg"
)

func Test_getDotDir(t *testing.T) {
	t.Run("should end with dotDirSuffix", func(t *testing.T) {
		got := getDotDir()
		if !strings.HasSuffix(got, cfg.Directory) {
			t.Errorf("getDotDir() = %v, want path to end with %v", got, cfg.Directory)
		}
	})
	t.Run("should start with homedir", func(t *testing.T) {
		got := getDotDir()
		usr, _ := user.Current()
		want := usr.HomeDir
		if !strings.HasPrefix(got, want) {
			t.Errorf("getDotDir() = %v, want path to start with %v", got, want)
		}
	})
}

func tempDirHelper() string {
	dir, err := ioutil.TempDir("", "xbsh")
	if err != nil {
		log.Fatal(err)
	}
	return dir
}
func Test_makeDotDir(t *testing.T) {
	type args struct {
		dir string
	}
	dir := tempDirHelper()
	defer os.RemoveAll(dir)

	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{"should return an error when called with empty string", args{""}, true},
		{"should return without error when called with TempDir", args{dir}, false},
		{"should return without error when called again with TempDir", args{dir}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := makeDotDir(tt.args.dir); (err != nil) != tt.wantErr {
				t.Errorf("makeDotDir() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
	t.Run("should have created TempDir", func(t *testing.T) {
		if _, err := os.Stat(dir); os.IsNotExist(err) {
			t.Errorf("makeDotDir() hsould have created %v, got: %v", dir, err)
		}
	})
}

func Test_openDB(t *testing.T) {
	type args struct {
		dbPath string
	}
	var arg, wrongArg args
	dir := tempDirHelper()
	defer os.RemoveAll(dir)
	arg.dbPath = dir + "/" + "test.sqlite"
	wrongArg.dbPath = "/xx"
	tests := []struct {
		name     string
		args     args
		wantPing bool
		wantErr  bool
	}{
		{"should return a db without error on first open", arg, true, false},
		{"should return a db without error on second open", arg, true, false},
		{"should return no db and PingError on impossible path", wrongArg, false, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotDb, err := openDB(tt.args.dbPath)
			if (err != nil) != tt.wantErr {
				t.Errorf("openDB() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			gotPing := gotDb.Ping()
			if (gotPing==nil)!=tt.wantPing {
				t.Errorf("openDB() PingError = %v, want %v", gotPing, tt.wantPing)
			}
		})
	}
}