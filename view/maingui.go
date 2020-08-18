package view

import (
	"log"
	"github.com/jroimartin/gocui"
)

const inputViewHeight = 4
const searchViewWidth = 50

// MainGui initializes the main gui and calls the views
func MainGui() {
	gui, err := gocui.NewGui(gocui.OutputNormal)
	if err != nil {
		log.Panicln(err)
	}
	defer gui.Close()

	gui.Cursor = true
	outputViewWidth,_ := gui.Size()
	outputViewWidth -= searchViewWidth

	inputView := Input("inputView", inputViewHeight, gui)
	outputView := Output("outputView", outputViewWidth, inputViewHeight, gui)
	searchView := Search("searchView", searchViewWidth, inputViewHeight, gui)
	focus := gocui.ManagerFunc(SetFocus("inputView"))
	gui.SetManager(inputView, focus, outputView, searchView)

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