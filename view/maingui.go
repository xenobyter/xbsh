package view

import (
	"log"

	"github.com/jroimartin/gocui"
)

var (
	gui            *gocui.Gui
	vMainView      *tMainView
	vInputView     *tInputView
	vStatusBarView *tStatusBarView
	vHelpView      *tHelpView
)

const inputHeight = 4

// MainGui initializes the main gui and calls the views
func MainGui() {
	var err error
	gui, err = gocui.NewGui(gocui.OutputNormal)
	if err != nil {
		log.Panicln(err)
	}
	defer gui.Close()

	gui.Highlight = true
	gui.Cursor = true
	_, guiHeight := gui.Size()

	// Initialize views
	vMainView = vMain("MainView", guiHeight-inputHeight)
	vStatusBarView = vStatusBar("StatusBar")
	vInputView = vInput("Inputview", inputHeight)
	focus := gocui.ManagerFunc(SetFocus("Inputview"))
	vHelpView = vHelp("HelpView")

	//Start ticker for Statusbar
	go updateStatus(gui)

	gui.SetManager(vMainView, vInputView, focus, vStatusBarView, vHelpView)
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
