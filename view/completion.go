package view

import (
	"fmt"
	"log"
	"os"

	"github.com/awesome-gocui/gocui"
	"github.com/xenobyter/xbsh/cfg"
)

type completionView struct {
	name, line  string
	view        *gocui.View
	completions []os.FileInfo
}

//Completion opens the completion view if  more then one completion is found
func Completion(completions []os.FileInfo) string {
	v := newCompletionView("Completion", completions)
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

func newCompletionView(name string, completions []os.FileInfo) *completionView {
	return &completionView{name: name, completions: completions}
}

func (i *completionView) layout(g *gocui.Gui) error {
	var (
		err error
		out string
	)
	maxX, maxY := g.Size()
	if i.view, err = g.SetView(i.name, 0, 0, maxX-1, maxY-2, 0); err != nil {
		if !gocui.IsUnknownView(err) {
			return err
		}
		i.view.Editor = i
		i.view.Editable = true
		i.view.Wrap = true
		i.view.Frame = true
		i.view.Highlight = true
		i.view.Title = i.name

		for _, c := range i.completions {
			if c.IsDir() {
				out += cfg.PathSep + c.Name() + "\n"
			} else {
				out += " " + c.Name() + "\n"
			}
		}
		fmt.Fprint(i.view, out)
		i.view.SetCursor(0, 0)
		if _, err := g.SetCurrentView(i.name); err != nil {
			return err
		}
	}
	return nil
}

func (i *completionView) quit(g *gocui.Gui, v *gocui.View) error {
	_, cy := v.Cursor()
	_, oy := v.Origin()
	i.line = v.BufferLines()[cy+oy][1:]
	return gocui.ErrQuit
}

func (i *completionView) keybindings(g *gocui.Gui) {
	if err := g.SetKeybinding("", gocui.KeyEsc, gocui.ModNone, quit); err != nil {
		log.Panicln(err)
	}
	if err := g.SetKeybinding("", gocui.KeyTab, gocui.ModNone, quit); err != nil {
		log.Panicln(err)
	}
	if err := g.SetKeybinding("", gocui.KeyEnter, gocui.ModNone, i.quit); err != nil {
		log.Panicln(err)
	}
}

func (i *completionView) Edit(view *gocui.View, key gocui.Key, char rune, mod gocui.Modifier) {
	_, cy := view.Cursor()
	_, oy := view.Origin()
	switch {
	case key == gocui.KeyArrowDown:
		if cy+oy < len(i.completions)-1 {
			view.MoveCursor(0, 1, true)
		}
	case key == gocui.KeyArrowUp:
		view.MoveCursor(0, -1, true)
	}
}
