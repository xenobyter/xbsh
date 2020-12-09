/*
Package cfg is used to manage load, store and export configuration
*/
package cfg

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"
)

func Test_loadConfig(t *testing.T) {
	dir, _ := ioutil.TempDir("", "xbsh")
	conffile := filepath.Join(dir, "xbshrc")
	defer os.RemoveAll(dir)
	loadConfig(conffile)
	t.Run("PathSep default value", func(t *testing.T) {
		if PathSep != "/" {
			t.Errorf("loadConfig() = %v, want %v", PathSep, "/")
		}
		if HistoryMaxEntries != 1000 {
			t.Errorf("loadConfig() = %v, want %v", HistoryMaxEntries, 1000)
		}
	})

	f, _ := os.Create(conffile)
	f.WriteString(`
	[separators]
	pathsep=|
	
	[history]
	maxentries=10
	`)
	loadConfig(conffile)
	t.Run("PathSep new value from file", func(t *testing.T) {
		if PathSep != "|" {
			t.Errorf("loadConfig() = %v, want %v", PathSep, "|")
		}
		if HistoryMaxEntries != 10 {
			t.Errorf("loadConfig() = %v, want %v", HistoryMaxEntries, 10)
		}
	})

}
