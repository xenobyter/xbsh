package view

import (
	"strings"

	"github.com/jroimartin/gocui"

	"github.com/xenobyter/xbsh/cmd"
)

// An InputView represents the main editor view at the bottom
type InputView struct {
	name   string
	height int
	gui    *gocui.Gui
}

// Input returns a pointer to the main editor view
func Input(name string, height int, gui *gocui.Gui) *InputView {
	return &InputView{name: name, height: height, gui: gui}
}

// Layout implements the layout for Input
func (i *InputView) Layout(gui *gocui.Gui) error {
	maxX, maxY := gui.Size()

	view, err := gui.SetView(i.name, 0, maxY-i.height-2, maxX-1, maxY-2)

	//Delete Inputview if terminal height too small
	if maxY<i.height+2 { gui.DeleteView(i.name) }
	if err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		view.Editor = i
		view.Editable = true
	}
	view.Title = getPrompt()
	return nil
}

// Edit implements the main editor and calls functions for keyhandling
func (i *InputView) Edit(view *gocui.View, key gocui.Key, char rune, mod gocui.Modifier) {
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
	case key == gocui.KeyEnter:
		cmdString := trimCmdString(view.Buffer())
		view.Clear()
		view.SetCursor(0,0)
		res := cmd.ExecCmd(cmdString)
		OutputWrite(res)
	}
}

func trimCmdString(buffer string) string {
	return strings.TrimSuffix(buffer, "\n")
}

func getPrompt() string {
	return cmd.GetWorkDir() + ":"
}
