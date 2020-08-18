package view

import (
	"github.com/jroimartin/gocui"
)

// An OutputView represents the main output view at the right
type OutputView struct {
	name         string
	width        int
	bottomMargin int
	gui          *gocui.Gui
}

// Output returns a pointer to the main output view
func Output(name string, width, bottomMargin int, gui *gocui.Gui) *OutputView {
	return &OutputView{name: name, width: width, bottomMargin: bottomMargin, gui: gui}
}

// Layout implements the layout for Output
func (i *OutputView) Layout(gui *gocui.Gui) error {
	maxX, maxY := gui.Size()
	view, err := gui.SetView(i.name, maxX-i.width, 0, maxX-1, maxY-i.bottomMargin-3)

	if err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		view.Autoscroll = true
		view.Title = "Output"
	}
	return nil
}
