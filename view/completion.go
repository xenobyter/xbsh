package view

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"github.com/jroimartin/gocui"
)

const (
	ansiBold      = "\u001b[1m"
	ansiUnderLine = "\u001b[4m"
)

type tCompletionView struct {
	name        string
	visible     bool
	view        *gocui.View
	itemLength  int
	completions []os.FileInfo
}

func newCompletionView(name string) *tCompletionView {
	return &tCompletionView{name: name, visible: false}
}

// Layout is called every time the GUI is redrawn
func (i *tCompletionView) Layout(gui *gocui.Gui) error {
	var err error

	x1, y1, x2, y2, out := i.prepareView()
	i.view, err = gui.SetView(i.name, x1, y1, x2, y2)
	if err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}

		//Startup tasks
		i.view.Autoscroll = false
		i.view.Frame = true
		i.view.Editor = i
		i.view.Editable = true
		i.view.Highlight = true
		fmt.Fprintln(i.view, out)
	}
	if !i.visible {
		//Hide the view
		gui.DeleteView(i.name)
	} else {
		gui.SetCurrentView(i.view.Name())
	}
	return nil
}

func (i *tCompletionView) Edit(view *gocui.View, key gocui.Key, char rune, mod gocui.Modifier) {
	_, cy := view.Cursor()
	_, oy := view.Origin()
	switch {
	case key == gocui.KeyEsc:
		i.toggle()
	case key == gocui.KeyArrowDown:
		if cy+oy < len(i.completions)-1 {
			view.MoveCursor(0, 1, true)
		}
	case key == gocui.KeyArrowUp:
		view.MoveCursor(0, -1, true)
	case key == gocui.KeyEnter || key == gocui.KeyTab:
		fmt.Fprint(vInputView.view, view.BufferLines()[cy+oy][i.itemLength:])
		if i.completions[cy+oy].IsDir() {
			fmt.Fprint(vInputView.view, "/")
		}
		vInputView.positionEnd()
		i.toggle()
	}
}

func (i *tCompletionView) toggle() {
	i.visible = !i.visible
}

func (i *tCompletionView) prepareView() (x1, y1, x2, y2 int, out string) {
	mx, _ := gui.Size()
	_, vy := vMainView.view.Size()
	cx, cy := vInputView.view.Cursor()

	//start above cursor
	x1 = cx
	y1 = vy + 3 - len(i.completions)
	x2 = 0
	y2 = vy + 4 + cy

	//prepare output and search for width
	for _, c := range i.completions {
		if c.IsDir() {
			out += ansiUnderLine + c.Name() + ansiNormal + "\n"
		} else {
			out += c.Name() + "\n"
		}
		if len(c.Name()) > x2 {
			x2 = len(c.Name())
		}
	}

	//add width to start
	x2 += x1 + 1

	//safety checks
	if x2 > mx {
		x2 = mx
	}
	if y1 < 0 {
		y1 = 1
	}

	return
}

//completionSearch takes search term and directory. It returns a slice of files and directorys
//the names of wich start with search term
func completionSearch(srch, dir string) (completions []os.FileInfo) {
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		return
	}

	for _, file := range files {
		if strings.HasPrefix(file.Name(), srch) {
			completions = append(completions, file)
		}
	}
	return
}
