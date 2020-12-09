/*
Package cfg is used to manage load, store and export configuration
*/
package cfg

import (
	"log"
	"os"
	"path/filepath"

	"github.com/sasbury/mini"
)

// Export path and filename for config
const (
	Directory = ".xbsh"
	DBFileName  = "storage.sqlite"
	INIFileName = "xbshrc"
)

// Export config
var (
	//[seperators]
	PipeSep string
	PathSep string

	//[history]
	HistoryMaxEntries int64
	HistoryDelExitCmd string
)

func init() {
	home, _ := os.UserHomeDir()
	confFile := filepath.Join(home, Directory, INIFileName)
	loadConfig(confFile)
}

func loadConfig(confFile string) {
	if _, err := os.Stat(confFile); os.IsNotExist(err) {
		f, _ := os.Create(confFile)
		f.WriteString("#xbsh config\n")
		f.Close()
	}
	conf, err := mini.LoadConfiguration(confFile)
	if err != nil {
		log.Fatalln(err)
	}

	PathSep = conf.StringFromSection("separators", "pathsep", "/")
	PipeSep = conf.StringFromSection("separators", "pipesep", "|")
	HistoryMaxEntries = conf.IntegerFromSection("history", "maxentries", 1000)
	HistoryDelExitCmd = conf.StringFromSection("history", "delExitCmd", "exit")
}
