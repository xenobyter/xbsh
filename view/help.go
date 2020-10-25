package view

import (
	"fmt"
	"log"

	"github.com/awesome-gocui/gocui"
)

// Help displays a simple help text
func Help(string) string {
	h := newHelpView("Help")
	g, err := gocui.NewGui(gocui.OutputNormal, false)
	if err != nil {
		log.Panicln(err)
	}
	defer g.Close()

	g.SetManagerFunc(h.layout)
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
	if v, err := g.SetView("help", 0, 0, maxX-1, maxY-1, 0); err != nil {
		if !gocui.IsUnknownView(err) {
			return err
		}
		v.Autoscroll = false
		v.Frame = true
		v.Title = h.name
		v.Wrap = true
		fmt.Fprintln(v, helptext)
		if _, err := g.SetCurrentView("help"); err != nil {
			return err
		}
	}
	return nil
}

func quit(g *gocui.Gui, v *gocui.View) error {
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
	if v != nil {
		ox, oy := v.Origin()
		_, size := v.Size()
		lenBuf := len(v.BufferLines()) - size
		oy += dy
		switch {
		case oy < 0:
			oy = 0
		case oy > lenBuf:
			oy = lenBuf
		}
		v.SetOrigin(ox, oy)
	}
	return nil
}
