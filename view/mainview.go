package view

import (
	"fmt"
	"strings"

	"github.com/jroimartin/gocui"

	"github.com/xenobyter/xbsh/cmd"
)

type tMainView struct {
	name, prompt string
	view         *gocui.View
}

func vMain(name string) *tMainView {
	return &tMainView{name: name}
}

// Layout is called every time the GUI is redrawn
func (i *tMainView) Layout(gui *gocui.Gui) error {
	var err error
	maxX, maxY := gui.Size()

	i.view, err = gui.SetView(i.name, 0, 0, maxX-1, maxY-2)
	if err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}

		//Startup tasks
		i.view.Editor = i
		i.view.Editable = true
		i.view.Title = i.name
		i.setPrompt()
	}
	return nil
}

// Edit implements the main editor and calls functions for keyhandling
func (i *tMainView) Edit(view *gocui.View, key gocui.Key, char rune, mod gocui.Modifier) {
	cx, cy := view.Cursor()
	ox, _ := view.Origin()
	maxX, _ := view.Size()
	limit := ox+maxX*cy+cx-1 > 512
	view.Wrap = true
	switch {
	case char != 0 && mod == 0 && !limit:
		view.EditWrite(char)
	case key == gocui.KeySpace && !limit:
		view.EditWrite(' ')
	case key == gocui.KeyBackspace || key == gocui.KeyBackspace2:
		view.EditDelete(true)
	case key == gocui.KeyF1:
		vHelpView.toggle()
	case key == gocui.KeyEnter:
		cmdString := trimLine(view.BufferLines())
		view.EditNewLine()
		i.print(cmd.ExecCmd(cmdString))
	}
}

func trimLine(bufferLines []string) string {
	buffer := bufferLines[len(bufferLines)-1]
	if i:= strings.Index(buffer, "$" + "\u202f"); i!=-1 {
		buffer = buffer[i+4:]
	}
	return strings.TrimSuffix(buffer, "\n")
}

func (i *tMainView) setPrompt() {
	i.prompt = cmd.GetPrompt()
	fmt.Fprintf(i.view, i.prompt)
	x, y := i.view.Cursor()
	i.view.SetCursor(x+len(i.prompt)-2, y)
}

func (i *tMainView) print(stdout, stderr []byte) {
	i.view.Write(stderr)
	i.view.Write(stdout)
}
