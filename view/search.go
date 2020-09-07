package view

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"github.com/jroimartin/gocui"
)

// An SearchView represents the search view at the left
type SearchView struct {
	name         string
	width        int
	bottomMargin int
	gui          *gocui.Gui
}

var vSearch *gocui.View

// Search returns a pointer to the search view
func Search(name string, width, bottomMargin int, gui *gocui.Gui) *SearchView {
	return &SearchView{name: name, width: width, bottomMargin: bottomMargin, gui: gui}
}

// Layout implements the layout for Search
func (i *SearchView) Layout(gui *gocui.Gui) error {
	var err error
	_, maxY := gui.Size()
	vSearch, err = gui.SetView(i.name, 0, 0, i.width-1, maxY-i.bottomMargin-3)

	if err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		vSearch.Autoscroll = true
		vSearch.Title = "Search"
	}
	return nil
}

func getDir(dir string) (files []os.FileInfo) {
	entries, err := ioutil.ReadDir(dir)
	if err != nil {
		log.Fatal(err)
	}
	for _, entry := range entries {
		files = append(files, entry)
	}
	return
}

func printDir(dir string) {
	entries := getDir(dir)
	for _, entry := range entries {
		fmt.Fprintln(vSearch, entry.Name())
	}
}
