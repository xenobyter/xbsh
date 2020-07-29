package view

import (
	"github.com/jroimartin/gocui"
)

type InputView struct {
	name      string
	x1, y1    int
	x2, y2    int
	maxLength int
}

// Input ...
func Input(name string, x1, y1, x2, y2, maxLength int) *InputView {
	return &InputView{name: name, x1: x1, y1: y1, x2: x2, y2: y2, maxLength: maxLength}
}

func (i *InputView) Layout(gui *gocui.Gui) error {
	v, err := gui.SetView(i.name, i.x1, i.y1, i.x2, i.y2)
	if err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		v.Editor = i
		v.Editable = true
	}
	return nil
}

func (i *InputView) Edit(v *gocui.View, key gocui.Key, ch rune, mod gocui.Modifier) {
	cx, _ := v.Cursor()
	ox, _ := v.Origin()
	limit := ox+cx+1 > i.maxLength
	switch {
	case ch != 0 && mod == 0 && !limit:
		v.EditWrite(ch)
	case key == gocui.KeySpace && !limit:
		v.EditWrite(' ')
	case key == gocui.KeyBackspace || key == gocui.KeyBackspace2:
		v.EditDelete(true)
	}
}