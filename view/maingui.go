package view

import (
	"log"
	"github.com/jroimartin/gocui"
)

const (
	vInputHeight = 4
	vSearchWidth = 50
)

// MainGui initializes the main gui and calls the views
func MainGui() {
	gui, err := gocui.NewGui(gocui.OutputNormal)
	if err != nil {
		log.Panicln(err)
	}
	defer gui.Close()

	gui.Highlight = true
	gui.Cursor = true

	guiWidth,_ := gui.Size()
	vOutputWidth := guiWidth-vSearchWidth

	vInput := Input("vInput", vInputHeight, gui)
	vOutput := Output("vOutput", vOutputWidth, vInputHeight, gui)
	vSearch := Search("vSearch", vSearchWidth, vInputHeight, gui)
	focus := gocui.ManagerFunc(SetFocus("vInput"))
	gui.SetManager(vInput, focus, vOutput, vSearch)

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