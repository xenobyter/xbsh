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
		i.view.Autoscroll = true
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
		view.EditDelete(true)
	case key == gocui.KeyF1:
		vHelpView.toggle()
	case key == gocui.KeyArrowLeft:
		i.cursorLeft()
	case key == gocui.KeyArrowRight:
		i.cursorRight()
	case key == gocui.KeyEnd:
		i.cursorEnd()
	case key == gocui.KeyEnter:
		cmdString := trimLine(view.BufferLines())
		fmt.Fprintln(vMainView.view, ansiPrompt+cmd.GetPrompt()+ansiNormal+cmdString)
		vMainView.print(cmd.ExecCmd(cmdString))
		i.view.Clear()
		i.setPrompt()
	}
}

func trimLine(bufferLines []string) string {
	buffer := bufferLines[len(bufferLines)-1]
	if i := strings.Index(buffer, "$"+"\u202f"); i != -1 {
		buffer = buffer[i+4:]
	}
	return strings.TrimSuffix(buffer, "\n")
}

func (i *tInputView) setPrompt() {
	i.prompt = cmd.GetPrompt()
	fmt.Fprintf(i.view, ansiPrompt+i.prompt+ansiNormal)
	_, y := i.view.Cursor()
	i.view.SetCursor(utf8.RuneCountInString(i.prompt), y)
}

func (i *tInputView) cursorLeft() {
	if pos, _ := i.view.Cursor(); pos > utf8.RuneCountInString(i.prompt) {
		i.view.MoveCursor(-1, 0, true)
	}
}

func (i *tInputView) cursorRight() {
	_, length := i.getLastLine()
	if pos, _ := i.view.Cursor(); pos < length {
		i.view.MoveCursor(1, 0, true)
	}
}

func (i *tInputView) cursorEnd() {
	_, length := i.getLastLine()
	x, y := i.view.Cursor()
	i.view.MoveCursor(length-x, y, true)
}

func (i *tInputView) getLastLine() (line string, length int) { //ToDo: #31 Get rid of getLastLine and ViewBufferLines
	line = i.view.ViewBufferLines()[len(i.view.ViewBufferLines())-1]
	length = utf8.RuneCountInString(line)
	return
}
