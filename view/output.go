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

var vOutput *gocui.View

// Output returns a pointer to the main output view
func Output(name string, width, bottomMargin int, gui *gocui.Gui) *OutputView {
	return &OutputView{name: name, width: width, bottomMargin: bottomMargin, gui: gui}
}

// Layout implements the layout for Output
func (i *OutputView) Layout(gui *gocui.Gui) error {
	var err error
	maxX, maxY := gui.Size()
	vOutput, err = gui.SetView(i.name, maxX-i.width, 0, maxX-1, maxY-i.bottomMargin-3)

	if err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		vOutput.Autoscroll = true
		vOutput.Title = "Output"
	}
	return nil
}

// OutputWrite writes to OutputView
func OutputWrite(out []byte){
	vOutput.Write(out)
}
