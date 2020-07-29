package main

import (
	"log"

	"github.com/jroimartin/gocui"

	"github.com/xenobyter/xbsh/view"
)

func setFocus(name string) func(gui *gocui.Gui) error {
	return func(gui *gocui.Gui) error {
		_, err := gui.SetCurrentView(name)
		return err
	}
}

func main() {
	gui, err := gocui.NewGui(gocui.OutputNormal)
	if err != nil {
		log.Panicln(err)
	}
	defer gui.Close()

	gui.Cursor = true
	maxX, maxY := gui.Size()
	inputView := view.Input("inputView", 1, maxY-6, maxX-1, maxY-1, 255)
	focus := gocui.ManagerFunc(setFocus("inputView"))
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
