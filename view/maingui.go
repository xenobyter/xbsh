package view

import (
	"os"
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
	gui.Mouse = true

	// Initialize views
	vMainView = vMain("MainView", inputHeight)
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

	//Mouse
	if err := gui.SetKeybinding("MainView", gocui.MouseWheelUp, gocui.ModNone, wheelUp); err != nil {
		log.Panicln(err)
	}
	if err := gui.SetKeybinding("MainView", gocui.MouseWheelDown, gocui.ModNone, wheelDown); err != nil {
		log.Panicln(err)
	}

	if err := gui.MainLoop(); err != nil && err != gocui.ErrQuit {
		gui.Close()
		log.Println("Terminal too small for xbsh. Exiting!")
		os.Exit(1)
	}
}

func quit(gui *gocui.Gui, v *gocui.View) error {
	return gocui.ErrQuit
}

func wheelUp(gui *gocui.Gui, v *gocui.View) error {
	vMainView.scrollMain(-1)
	return nil
}
func wheelDown(gui *gocui.Gui, v *gocui.View) error {
	vMainView.scrollMain(1)
	return nil
}
