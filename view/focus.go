package view

import "github.com/jroimartin/gocui"

// SetFocus takes the name of a view and set focus
func SetFocus(name string) func(gui *gocui.Gui) error {
	return func(gui *gocui.Gui) error {
		_, err := gui.SetCurrentView(name)
		return err
	}
}