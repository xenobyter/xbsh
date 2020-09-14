package view

import (
	"github.com/jroimartin/gocui"
)

type tMainView struct {
	name   string
	height int
	view   *gocui.View
}

func vMain(name string, height int) *tMainView {
	return &tMainView{name: name, height: height}
}

// Layout is called every time the GUI is redrawn
func (i *tMainView) Layout(gui *gocui.Gui) error {
	var err error
	maxX, _ := gui.Size()

	i.view, err = gui.SetView(i.name, 0, 0, maxX-1, i.height-3)
	if err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}

		//Startup tasks
		i.view.Autoscroll = true
		i.view.Title = i.name
	}
	return nil
}

func (i *tMainView) print(stdout, stderr []byte) {
	// lines := countLines(stdout, stderr)
	i.view.Write(stderr)
	i.view.Write(stdout)
}
