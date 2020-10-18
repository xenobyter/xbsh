package view

import (
	"log"
	"os"

	"github.com/jroimartin/gocui"
)

var (
	gui *gocui.Gui

	vMainView       *tMainView
	vInputView      *tInputView
	vStatusBarView  *tStatusBarView
	vHelpView       *tHelpView
	vHistoryView    *tHistoryView
	vCompletionView *tCompletionView

	overlayCoord = map[string]int{"x0": 3, "y0": 2, "x1": -4, "y1": -9}
)

const (
	inputHeight = 4
)

// MainGui initializes the main gui and calls the views
func MainGui() {
	var err error
	gui, err = gocui.NewGui(gocui.Output256)
	if err != nil {
		log.Panicln(err)
	}
	defer gui.Close()

	gui.Highlight = true
	gui.SelFgColor = gocui.ColorGreen
	gui.Cursor = true
	gui.Mouse = true
	gui.InputEsc = true

	// Initialize views
	vMainView = newMainView("Main", inputHeight)
	vStatusBarView = newStatusBarView("StatusBar")
	vInputView = newInputView("Input", inputHeight)
	focus := gocui.ManagerFunc(SetFocus("Input"))
	vHelpView = newHelpView("Help")
	vHistoryView = newHistoryView("History")
	vCompletionView = newCompletionView("Completion")

	//Start ticker for Statusbar
	go updateStatus(gui)

	gui.SetManager(vMainView, vInputView, focus, vStatusBarView, vHelpView, vHistoryView, vCompletionView)
	if err := gui.SetKeybinding("", gocui.KeyCtrlC, gocui.ModNone, quit); err != nil {
		log.Panicln(err)
	}

	//Mouse
	if err := gui.SetKeybinding("Main", gocui.MouseWheelUp, gocui.ModNone, wheelUp); err != nil {
		log.Panicln(err)
	}
	if err := gui.SetKeybinding("Main", gocui.MouseWheelDown, gocui.ModNone, wheelDown); err != nil {
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
