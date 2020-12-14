package view

import (
	"log"
	"unicode/utf8"

	"github.com/awesome-gocui/gocui"
)

type renameView struct {
	name         string
	lView, rView *gocui.View
}

// Rename opens the rename view
func Rename(string) string {
	v := newRenameView("Rename")
	g, err := gocui.NewGui(gocui.OutputNormal, false)
	if err != nil {
		log.Panicln(err)
	}
	defer g.Close()

	g.SetManagerFunc(v.layout)
	v.keybindings(g)
	g.Cursor = true
	g.Highlight = true
	if err := g.MainLoop(); err != nil && !gocui.IsQuit(err) {
		log.Panicln(err)
	}
	return ""
}

func newRenameView(name string) *renameView {
	return &renameView{name: name}
}

func (i *renameView) layout(g *gocui.Gui) error { //TODO: Integrate db.rename
	var (
		err error
	)
	maxX, maxY := g.Size()
	if i.lView, err = g.SetView(i.name, 0, 0, 25, maxY-2, 0); err != nil {
		if !gocui.IsUnknownView(err) {
			return err
		}
		i.lView.Wrap = false
		i.lView.Frame = true
		i.lView.Highlight = true
		i.lView.Title = i.name
		i.lView.Editable = true
		i.lView.Editor = i

		i.lView.SetCursor(0, 0)
		if _, err := g.SetCurrentView(i.name); err != nil {
			return err
		}
	}
	if i.rView, err = g.SetView("Output", 26, 0, maxX-1, maxY-2, 0); err != nil {
		if !gocui.IsUnknownView(err) {
			return err
		}
		i.rView.Wrap = false
		i.rView.Frame = true
		i.rView.Highlight = true
		i.rView.Title = "Output"
		i.rView.Autoscroll = true
	}
	return nil
}

func (i *renameView) keybindings(g *gocui.Gui) error {
	if err := g.SetKeybinding("", gocui.KeyEsc, gocui.ModNone, quit); err != nil {
		log.Panicln(err)
	}
	if err := g.SetKeybinding("", gocui.KeyF5, gocui.ModNone, quit); err != nil {
		log.Panicln(err)
	}
	return nil
}

func (i *renameView) Edit(v *gocui.View, key gocui.Key, ch rune, mod gocui.Modifier) {
	ox, oy := v.Origin()
	cx, cy := v.Cursor()
	_, my := v.Size()
	switch {
	case key == 0:
		v.EditWrite(ch)
	case key == gocui.KeySpace:
		v.EditWrite(' ')
	case key == gocui.KeyEnter:
		v.EditNewLine()
		if cy == my-1 {
			v.MoveCursor(0, 1, true)
		}
	case key == gocui.KeyBackspace || key == gocui.KeyBackspace2:
		v.EditDelete(true)
	case key == gocui.KeyDelete:
		v.EditDelete(false)
	case key == gocui.KeyArrowLeft:
		if ox+cx > 0 {
			v.MoveCursor(-1, 0, true)
		}
	case key == gocui.KeyArrowRight:
		line, _ := v.Line(oy + cy)
		if ox+cx < utf8.RuneCountInString(line) {
			v.MoveCursor(1, 0, true)
		}
	case key == gocui.KeyArrowUp:
		if oy+cy > 0 {
			v.MoveCursor(0, -1, true)
		}
	case key == gocui.KeyArrowDown:
		if cy+oy < len(v.BufferLines())-1 {
			v.MoveCursor(0, 1, true)
		}
	}
}
