package view

import (
	"fmt"
	"strings"
	"unicode/utf8"

	"github.com/jroimartin/gocui"

	"github.com/xenobyter/xbsh/cmd"
)

type tInputView struct {
	name   string
	height int
	prompt string
	view   *gocui.View
}

const (
	ansiPrompt = "\033[0;32m"
	ansiStdErr = "\033[0;31m"
	ansiNormal = "\033[0m"
)

func vInput(name string, height int) *tInputView {
	return &tInputView{name: name, height: height}
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
		i.view.Autoscroll = true //TODO: #39 #38 Long lines break prompt
		i.view.Title = i.name
		i.setPrompt()
	}
	return nil
}

// Edit implements the main editor and calls functions for keyhandling
func (i *tInputView) Edit(view *gocui.View, key gocui.Key, char rune, mod gocui.Modifier) {
	// view.Wrap = true
	switch {
	case char != 0 && mod == 0:
		view.EditWrite(char)
	case key == gocui.KeySpace:
		view.EditWrite(' ')
	case key == gocui.KeyBackspace || key == gocui.KeyBackspace2:
		view.EditDelete(true) //TODO: #40 Dont delete the prompt
	case key == gocui.KeyF1:
		vHelpView.toggle()
	case key == gocui.KeyArrowLeft:
		i.cursorLeft()
	case key == gocui.KeyArrowRight:
		i.cursorRight()
	case key == gocui.KeyEnd: //TODO: #45 Implement gocui.KeyHome
		i.cursorEnd()
	case key == gocui.KeyArrowUp && mod == gocui.ModAlt:
		vMainView.scrollMain(-1)
	case key == gocui.KeyArrowDown && mod == gocui.ModAlt:
		vMainView.scrollMain(1)
	case key == gocui.KeyPgup:
		vMainView.scrollMain(-8)
	case key == gocui.KeyPgdn:
		vMainView.scrollMain(8)
	case key == gocui.KeyEnter:
		cmdString := trimLine(view.BufferLines())
		fmt.Fprintln(vMainView.view, ansiPrompt+cmd.GetPrompt()+ansiNormal+cmdString)
		vMainView.print(cmd.ExecCmd(cmdString))
		i.view.Clear()
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

func (i *tInputView) cursorLeft() {
	if pos, _ := i.view.Cursor(); pos > utf8.RuneCountInString(i.prompt) {
		i.view.MoveCursor(-1, 0, true)
	}
}

func (i *tInputView) cursorRight() {
	if pos, _ := i.view.Cursor(); pos < i.bufferLength() {
		i.view.MoveCursor(1, 0, true)
	}
}

func (i *tInputView) cursorEnd() {
	x, y := i.view.Cursor()
	i.view.MoveCursor(i.bufferLength()-x, y, true)
}

func (i *tInputView) bufferLength() int {
	return utf8.RuneCountInString(i.view.ViewBuffer())-1
}
