package view

import (
	"github.com/jroimartin/gocui"
)

// An SearchView represents the main output view at the right
type SearchView struct {
	name         string
	width        int
	bottomMargin int
	gui          *gocui.Gui
}

// Search returns a pointer to the search view
func Search(name string, width, bottomMargin int, gui *gocui.Gui) *SearchView {
	return &SearchView{name: name, width: width, bottomMargin: bottomMargin, gui: gui}
}

// Layout implements the layout for Search
func (i *SearchView) Layout(gui *gocui.Gui) error {
	_, maxY := gui.Size()
	view, err := gui.SetView(i.name, 0, 0, i.width-1, maxY-i.bottomMargin-3)

	if err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		view.Autoscroll = true
	view.Title = "Search"
	}
	return nil
}
