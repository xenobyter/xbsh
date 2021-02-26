/*
Package view is used to implement all gocui-based views
*/
package view

import (
	"fmt"
	"log"

	"github.com/awesome-gocui/gocui"
	"github.com/xenobyter/xbsh/term"
)

// Help displays a simple help text
func Help(string) string {
	term.SetStatus("| ESC: Exit | F1: Exit")
	h := newHelpView("Help")
	g, err := gocui.NewGui(gocui.OutputNormal, false)
	if err != nil {
		log.Panicln(err)
	}
	defer g.Close()

	g.SetManagerFunc(h.layout)
	g.Cursor = true
	if err := h.keybindings(g); err != nil {
		log.Panicln(err)
	}
	if err := g.MainLoop(); err != nil && !gocui.IsQuit(err) {
		log.Panicln(err)
	}

	return h.line
}

type helpView struct {
	name, line string
}

func newHelpView(name string) *helpView {
	return &helpView{name: name, line: ""}
}

func (h *helpView) layout(g *gocui.Gui) error {
	maxX, maxY := g.Size()
	if v, err := g.SetView(h.name, 0, 0, maxX-1, maxY-2, 0); err != nil {
		if !gocui.IsUnknownView(err) {
			return err
		}
		v.Autoscroll = false
		v.Editable = true
		v.Frame = true
		v.Title = h.name
		v.Wrap = true
		fmt.Fprintln(v, helptext)
		v.SetCursor(0, 0)
		if _, err := g.SetCurrentView(h.name); err != nil {
			return err
		}
	}
	return nil
}

func quit(g *gocui.Gui, v *gocui.View) error {
	term.ClearStatus()
	return gocui.ErrQuit
}

func (h *helpView) keybindings(g *gocui.Gui) error {
	if err := g.SetKeybinding("", gocui.KeyEsc, gocui.ModNone, quit); err != nil {
		return err
	}
	if err := g.SetKeybinding("", gocui.KeyF1, gocui.ModNone, quit); err != nil {
		return err
	}
	if err := g.SetKeybinding("", gocui.KeyArrowDown, gocui.ModNone,
		func(g *gocui.Gui, v *gocui.View) error {
			return h.scrollView(v, 1)
		}); err != nil {
		return err
	}
	if err := g.SetKeybinding("", gocui.KeyArrowUp, gocui.ModNone,
		func(g *gocui.Gui, v *gocui.View) error {
			return h.scrollView(v, -1)
		}); err != nil {
		return err
	}
	return nil
}

func (h *helpView) scrollView(v *gocui.View, dy int) error {
	_, cy := v.Cursor()
	_, oy := v.Origin()
	switch ny := cy + oy + dy; {
	case ny < 0:
		v.SetCursor(0, 0)
		v.SetOrigin(0, 0)
	case ny < len(v.ViewBufferLines())-1:
		v.MoveCursor(0, dy, true)
	}
	return nil
}
