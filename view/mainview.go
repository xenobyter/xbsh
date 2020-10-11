package view

import (
	"fmt"

	"github.com/jroimartin/gocui"
)

type tMainView struct {
	name        string
	inputHeight int
	view        *gocui.View
}

func newMainView(name string, inputHeight int) *tMainView {
	return &tMainView{name: name, inputHeight: inputHeight}
}

// Layout is called every time the GUI is redrawn
func (i *tMainView) Layout(gui *gocui.Gui) error {
	var err error
	maxX, maxY := gui.Size()

	i.view, err = gui.SetView(i.name, 0, 0, maxX-1, maxY-i.inputHeight-3)
	if err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}

		//Startup tasks
		i.view.Autoscroll = false
		i.view.Title = i.name
	}
	return nil
}

func (i *tMainView) print(stdout, stderr []byte) {
	fmt.Fprint(i.view, ansiStdErr)
	i.view.Write(stderr)
	fmt.Fprint(i.view, ansiNormal)
	i.view.Write(stdout)

	//reposition origin if needed
	l := len(i.view.BufferLines())
	_, y := i.view.Size()
	if l > y {
		i.view.SetOrigin(0, l-y-1)
	}
}

func (i *tMainView) scrollMain(cnt int) {
	x, y := i.view.Origin()
	i.view.SetOrigin(x, y+cnt)
}

//TODO: #60 MainView Clearing