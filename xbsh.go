package main

import (
	"github.com/xenobyter/xbsh/storage"
	"github.com/xenobyter/xbsh/view"
)

func main() {
	//PATH cache
	go storage.PathCache()
	//initialize gui
	view.MainGui()
}

