package view

import (
	"fmt"

	"github.com/jroimartin/gocui"
)

type tHelpView struct {
	name    string
	view    *gocui.View
	visible bool
}

func vHelp(name string) *tHelpView {
	return &tHelpView{name: name, visible: false}
}

func (i *tHelpView) Layout(gui *gocui.Gui) error {
	var err error
	maxX, maxY := gui.Size()

	i.view, err = gui.SetView(i.name, 4, 4, maxX-5, maxY-6)
	if err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}

		//Startup tasks
		i.view.Title = i.name
		i.view.Wrap = true
		fmt.Fprintf(i.view, helptext)
	}
	if !i.visible {
		gui.DeleteView(i.name)
	}
	return nil
}

func (i *tHelpView) toggle() {
	i.visible = !i.visible
}
