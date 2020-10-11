package view

import (
	"fmt"
	"unicode/utf8"

	"github.com/xenobyter/xbsh/storage"

	"github.com/jroimartin/gocui"
)

type tHistoryView struct {
	name    string
	view    *gocui.View
	visible bool
}

const (
	offset = 8
)

func newHistoryView(name string) *tHistoryView {
	return &tHistoryView{name: name, visible: false}
}

func (i *tHistoryView) Layout(gui *gocui.Gui) error {
	var err error
	maxX, maxY := gui.Size()

	i.view, err = gui.SetView(i.name, overlayCoord["x0"], overlayCoord["y0"], maxX+overlayCoord["x1"], maxY+overlayCoord["y1"])
	if err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}

		//Startup tasks
		i.view.Editor = i
		i.view.Editable = true
		i.view.Title = i.name
		i.view.Wrap = true
		i.view.Highlight = true
	}
	if !i.visible {
		//Hide the view
		gui.DeleteView(i.name)
	} else {
		gui.SetCurrentView(i.view.Name())
	}
	return nil
}

func (i *tHistoryView) Edit(view *gocui.View, key gocui.Key, char rune, mod gocui.Modifier) {
	cx, cy := view.Cursor()
	ox, oy := view.Origin()
	switch {
	case char != 0 && mod == 0:
		i.resetPos(cx)
		view.EditWrite(char)
		i.search(view.BufferLines()[0][offset:])
	case key == gocui.KeySpace:
		i.resetPos(cx)
		view.EditWrite(' ')
		i.search(view.BufferLines()[0][offset:])
	case key == gocui.KeyBackspace || key == gocui.KeyBackspace2:
		i.resetPos(cx)
		if cx > offset {
			view.EditDelete(true)
		}
		i.search(view.BufferLines()[0][offset:])
	case key == gocui.KeyDelete:
		i.resetPos(cx)
		if cx >= offset && cx < utf8.RuneCountInString(view.BufferLines()[0]) {
			view.EditDelete(false)
		}
		i.search(view.BufferLines()[0][offset:])
	case key == gocui.KeyF2:
		i.toggle()
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
	case key==gocui.KeyHome:
		view.SetOrigin(0,0)
		view.SetCursor(offset,0)
	case key == gocui.KeyPgup:
		if oy > 8 {
			view.SetOrigin(ox, oy-8)
		} else {
			i.resetPos(cx)
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
	case key == gocui.KeyEnter:
		cmdString := view.BufferLines()[cy]
		if cy == 0 {
			cmdString = view.BufferLines()[cy][offset:]
		}
		i.toggle()
		vInputView.view.Clear()
		vInputView.setPrompt()
		fmt.Fprintf(vInputView.view, cmdString)
		vInputView.view.MoveCursor(utf8.RuneCountInString(cmdString), 0, true)
	}
}

func (i *tHistoryView) toggle() {
	i.visible = !i.visible
	//Trigger initial search
	if i.visible {
		go func() {
			srch := trimLine(vInputView.view.BufferLines())
			i.search(srch)
			i.view.SetCursor(8+utf8.RuneCountInString(srch), 0)
		}()
	}
}

func (i *tHistoryView) search(srch string) {
	res := storage.HistorySearch(srch)
	i.view.Clear()
	fmt.Fprintln(i.view, "Search:", srch)
	for _, r := range res {
		fmt.Fprintln(i.view, r)
	}
}

func (i *tHistoryView) resetPos(cx int) {
	i.view.SetOrigin(0, 0)
	i.view.SetCursor(cx, 0)
}
