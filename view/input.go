package view

import (
	"github.com/jroimartin/gocui"
)

// An InputView represents the main editor view at the bottom
type InputView struct {
	name      string
	x1, y1    int
	x2, y2    int
	maxLength int
}

// Input returns a pointer to the main editor view
func Input(name string, x1, y1, x2, y2, maxLength int) *InputView {
	return &InputView{name: name, x1: x1, y1: y1, x2: x2, y2: y2, maxLength: maxLength}
}

// Layout implements the layout for Input
func (i *InputView) Layout(gui *gocui.Gui) error {
	view, err := gui.SetView(i.name, i.x1, i.y1, i.x2, i.y2)
	if err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		view.Editor = i
		view.Editable = true
	}
	return nil
}

// Edit implements the main editor and calls functions for keyhandling
func (i *InputView) Edit(view *gocui.View, key gocui.Key, char rune, mod gocui.Modifier) {
	cx, cy := view.Cursor()
	ox, _ := view.Origin()
	maxX, _ := view.Size()
	limit := ox+maxX*cy+cx-1 > i.maxLength
	view.Wrap = true
	switch {
	case char != 0 && mod == 0 && !limit:
		view.EditWrite(char)
	case key == gocui.KeySpace && !limit:
		view.EditWrite(' ')
	case key == gocui.KeyBackspace || key == gocui.KeyBackspace2:
		view.EditDelete(true)
	case key == gocui.KeyArrowDown:
		view.BgColor = gocui.ColorBlue
	}
}
