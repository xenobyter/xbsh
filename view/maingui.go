package view

import (
	"log"

	"github.com/jroimartin/gocui"
)

var (
	gui            *gocui.Gui
	vMainView      *tMainView
	vStatusBarView *tStatusBarView
	vHelpView      *tHelpView
)

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

	// Initialize views
	vMainView = vMain("MainView")
	vStatusBarView = vStatusBar("StatusBar")
	focus := gocui.ManagerFunc(SetFocus("MainView"))
	vHelpView = vHelp("HelpView")

	//Start ticker for Statusbar
	go updateStatus(gui)

	gui.SetManager(vMainView, focus, vStatusBarView, vHelpView)
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
