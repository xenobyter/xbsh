package view

import (
	"fmt"
	"strings"
	"unicode/utf8"

	"github.com/jroimartin/gocui"

	"github.com/xenobyter/xbsh/cmd"
	"github.com/xenobyter/xbsh/storage"
)

type tInputView struct {
	name   string
	height int
	prompt string
	view   *gocui.View
	hPos   int64
}

const (
	ansiPrompt = "\033[0;32m"
	ansiStdErr = "\033[0;31m"
	ansiNormal = "\033[0m"
)

func newInputView(name string, height int) *tInputView {
	return &tInputView{name: name, height: height, hPos: 0}
}

// Layout is called every time the GUI is redrawn
func (i *tInputView) Layout(gui *gocui.Gui) error {
	var err error
	maxX, maxY := gui.Size()

	i.view, err = gui.SetView(i.name, 0, maxY-i.height-2, maxX-1, maxY-2)
	if err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}

		//Startup tasks
		i.view.Editor = i
		i.view.Editable = true
		i.view.Autoscroll = false
		i.view.Wrap = true
		i.view.Title = i.name
		i.setPrompt()
	}
	return nil
}

// Edit implements the main editor and calls functions for keyhandling
func (i *tInputView) Edit(view *gocui.View, key gocui.Key, char rune, mod gocui.Modifier) {
	cx, cy := view.Cursor()
	_, oy := view.Origin()
	mx, my := view.Size()
	offset := utf8.RuneCountInString(cmd.GetPrompt())
	length := utf8.RuneCountInString(i.view.BufferLines()[0])
	switch {
	case char != 0 && mod == 0:
		view.EditWrite(char)
	case key == gocui.KeySpace:
		view.EditWrite(' ')
	case key == gocui.KeyBackspace || key == gocui.KeyBackspace2:
		if cx > offset || cy > 0 || oy > 0 {
			view.EditDelete(true)
		}
	case key == gocui.KeyDelete:
		if cx >= offset || cy > 0 || oy > 0 {
			view.EditDelete(false)
		}
	case key == gocui.KeyF1:
		vHelpView.toggle()
	case key == gocui.KeyF2:
		vHistoryView.toggle()
	case key == gocui.KeyArrowLeft:
		i.cursorLeft()
	case key == gocui.KeyArrowRight:
		if cx+(mx+1)*(cy+oy) < length {
			i.view.MoveCursor(1, 0, true)
		}
	case key == gocui.KeyEnd:
		x, y, o := caclulateCursor(length, mx, my)
		view.SetCursor(x, y)
		view.SetOrigin(0, o)
	case key == gocui.KeyHome:
		view.SetOrigin(0, 0)
		view.SetCursor(offset, 0)
	case key == gocui.KeyArrowUp && mod == gocui.ModAlt:
		vMainView.scrollMain(-1)
	case key == gocui.KeyArrowDown && mod == gocui.ModAlt:
		vMainView.scrollMain(1)
	case key == gocui.KeyPgup:
		vMainView.scrollMain(-8)
	case key == gocui.KeyPgdn:
		vMainView.scrollMain(8)
	case key == gocui.KeyArrowUp:
		i.scrollHistory(-1)
	case key == gocui.KeyArrowDown:
		i.scrollHistory(1)
	case key == gocui.KeyEnter:
		cmdString := trimLine(view.BufferLines())
		fmt.Fprintln(vMainView.view, ansiPrompt+cmd.GetPrompt()+ansiNormal+cmdString)
		vMainView.print(cmd.ExecCmd(cmdString))
		go func() { i.hPos = storage.HistoryWrite(cmdString) + 1 }()
		i.view.Clear()
		i.view.SetOrigin(0, 0)
		i.setPrompt()
	}
}

func trimLine(bufferLines []string) string {
	const sep = "$ "
	buffer := bufferLines[len(bufferLines)-1]
	if i := strings.Index(buffer, sep); i != -1 {
		buffer = buffer[i+len(sep):]
	}
	return strings.TrimSuffix(buffer, "\n")
}

func (i *tInputView) setPrompt() {
	i.prompt = cmd.GetPrompt()
	fmt.Fprintf(i.view, ansiPrompt+i.prompt+ansiNormal)
	i.view.SetCursor(utf8.RuneCountInString(i.prompt), 0)
}

func (i *tInputView) scrollHistory(inc int64) {
	var history string
	history, i.hPos = storage.HistoryRead(i.hPos + inc)
	i.view.Clear()
	i.setPrompt()
	fmt.Fprint(i.view, history)
	mx, my := i.view.Size()
	length := utf8.RuneCountInString(i.view.BufferLines()[0])
	cx, cy, oy := caclulateCursor(length, mx, my)
	i.view.SetCursor(cx, cy)
	i.view.SetOrigin(0, oy)
	return
}

func caclulateCursor(l, mx, my int) (cx, cy, oy int) {
	lines := l / mx
	switch {
	case lines == 0:
		cx = l
	case lines < my:
		cy = lines
		cx = l - lines*mx
	case lines >= my:
		cy = my - 1
		oy = lines - my + 1
		cx = l - lines*mx
	}
	return cx, cy, oy
}

func (i *tInputView) cursorLeft() {
	cx, cy := i.view.Cursor()
	mx, _ := i.view.Size()
	_, oy := i.view.Origin()
	p := utf8.RuneCountInString(i.prompt)
	switch {
	case cy == 0 && oy == 0:
		if cx > p {
			i.view.MoveCursor(-1, 0, true)
		}
	case cy == 0 && cx == 0 && oy > 0:
		i.view.SetCursor(mx-10, 0)
		i.view.SetOrigin(0, oy-1)
	case cy > 0 || oy > 0:
		i.view.MoveCursor(-1, 0, true)
	}
}
