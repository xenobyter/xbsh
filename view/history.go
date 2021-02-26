package view

import (
	"fmt"
	"log"
	"unicode/utf8"

	"github.com/awesome-gocui/gocui"
	"github.com/xenobyter/xbsh/db"
	"github.com/xenobyter/xbsh/term"
)

const (
	offset = 8
)

//History takes a string and uses it to search in db. It displays the filtered
//results. One entry can be selected with Arrowkeys. Enter returns this entry
func History(line string) string {
	term.SetStatus("| Esc: Exit | F2: Exit | Enter: Use Entry | Type to refine search")
	v := newHistoryView("History", line)
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
	return v.line
}

type historyView struct {
	name, line string
	view       *gocui.View
}

func newHistoryView(name, line string) *historyView {
	return &historyView{name: name, line: line}
}

func (h *historyView) keybindings(g *gocui.Gui) {
	if err := g.SetKeybinding("", gocui.KeyEsc, gocui.ModNone, quit); err != nil {
		log.Panicln(err)
	}
	if err := g.SetKeybinding("", gocui.KeyF2, gocui.ModNone, quit); err != nil {
		log.Panicln(err)
	}
	if err := g.SetKeybinding("", gocui.KeyEnter, gocui.ModNone, h.quit); err != nil {
		log.Panicln(err)
	}
}

func (h *historyView) quit(g *gocui.Gui, v *gocui.View) error {
	_, cy := v.Cursor()
	_, oy := v.Origin()
	h.line = v.BufferLines()[cy+oy]
	if cy+oy == 0 {
		h.line = v.BufferLines()[cy][offset:]
	}
	return gocui.ErrQuit
}

func (h *historyView) layout(g *gocui.Gui) error {
	var err error
	maxX, maxY := g.Size()
	if h.view, err = g.SetView(h.name, 0, 0, maxX-1, maxY-2, 0); err != nil {
		if !gocui.IsUnknownView(err) {
			return err
		}
		h.view.Editor = h
		h.view.Editable = true
		h.view.Wrap = true
		h.view.Frame = true
		h.view.Title = h.name
		h.view.Highlight = true
		h.search(h.line)
		h.view.SetCursor(offset, 0)
		if _, err := g.SetCurrentView(h.name); err != nil {
			return err
		}
	}
	return nil
}

func (h *historyView) Edit(view *gocui.View, key gocui.Key, char rune, mod gocui.Modifier) {
	cx, cy := view.Cursor()
	ox, oy := view.Origin()
	switch {
	case char != 0 && mod == 0:
		h.resetPos(cx)
		view.EditWrite(char)
		h.search(view.BufferLines()[0][offset:])
	case key == gocui.KeySpace:
		h.resetPos(cx)
		view.EditWrite(' ')
		h.search(view.BufferLines()[0][offset:])
	case key == gocui.KeyBackspace || key == gocui.KeyBackspace2:
		h.resetPos(cx)
		if cx > offset {
			view.EditDelete(true)
		}
		h.search(view.BufferLines()[0][offset:])
	case key == gocui.KeyDelete:
		h.resetPos(cx)
		if cx >= offset && cx < utf8.RuneCountInString(view.BufferLines()[0]) {
			view.EditDelete(false)
		}
		h.search(view.BufferLines()[0][offset:])

	case key == gocui.KeyArrowLeft:
		if cx > offset {
			view.MoveCursor(-1, 0, true)
		}
	case key == gocui.KeyArrowRight:
		if cx < len(view.BufferLines()[0]) {
			view.MoveCursor(1, 0, true)
		}
	case key == gocui.KeyEnd:
		view.SetCursor(len(view.BufferLines()[0]), 0)
	case key == gocui.KeyHome:
		view.SetOrigin(0, 0)
		view.SetCursor(offset, 0)
	case key == gocui.KeyPgup:
		if oy > 8 {
			view.SetOrigin(ox, oy-8)
		} else {
			h.resetPos(cx)
		}
	case key == gocui.KeyPgdn:
		switch {
		case cy+oy < len(view.BufferLines())-9:
			view.SetOrigin(ox, oy+8)
		case cy+oy < len(view.BufferLines())-2:
			view.MoveCursor(0, 1, true)
		}
	case key == gocui.KeyArrowUp:
		switch {
		case cy > 0:
			view.MoveCursor(0, -1, true)
		case oy > 0:
			view.SetOrigin(ox, oy-1)
		}
	case key == gocui.KeyArrowDown:
		if cy+oy < len(view.BufferLines())-2 {
			view.MoveCursor(0, 1, true)
		}
	}
}

func (h *historyView) search(srch string) {
	res := db.HistorySearch(srch)
	h.view.Clear()
	fmt.Fprintln(h.view, "Search:", srch)
	for _, r := range res {
		fmt.Fprintln(h.view, r)
	}
}

func (h *historyView) resetPos(cx int) {
	h.view.SetOrigin(0, 0)
	h.view.SetCursor(cx, 0)
}
