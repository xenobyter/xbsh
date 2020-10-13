package main

import (
	"os"

	"github.com/xenobyter/xbsh/storage"
	"github.com/xenobyter/xbsh/view"
)

func main() {
	wd,_:=os.Getwd()
	go storage.WorkDirScan(wd)
	//initialize gui
	view.MainGui()
}

