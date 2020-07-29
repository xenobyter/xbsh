package view

import (
	"log"
	"github.com/jroimartin/gocui"
)

// MainGui initializes the main gui and calls the views
func MainGui() {
	gui, err := gocui.NewGui(gocui.OutputNormal)
	if err != nil {
		log.Panicln(err)
	}
	defer gui.Close()

	gui.Cursor = true
	maxX, maxY := gui.Size()
	inputView := Input("inputView", 1, maxY-6, maxX-1, maxY-1, 512)
	focus := gocui.ManagerFunc(SetFocus("inputView"))
	gui.SetManager(inputView, focus)

	if err := gui.SetKeybinding("", gocui.KeyCtrlC, gocui.ModNone, quit); err != nil {
		log.Panicln(err)
	}

	if err := gui.MainLoop(); err != nil && err != gocui.ErrQuit {
		log.Panicln(err)
	}
}
func quit(gui *gocui.Gui, v *gocui.View) error {
	return gocui.ErrQuit
}