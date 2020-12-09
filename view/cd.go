package view

import (
	"fmt"
	"io/ioutil"
	"log"
	"strings"

	"github.com/awesome-gocui/gocui"
	"github.com/xenobyter/xbsh/cfg"
)

//CD displays a view
func CD(path string) string {
	v := newCDView(path)
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

type cdView struct {
	name, line, path string
	view             *gocui.View
}

func newCDView(path string) *cdView {
	return &cdView{path: path, line: "", name: "cd"}
}

func (i *cdView) layout(g *gocui.Gui) error {
	var err error
	maxX, maxY := g.Size()
	if i.view, err = g.SetView(i.name, 0, 0, maxX-1, maxY-2, 0); err != nil {
		if !gocui.IsUnknownView(err) {
			return err
		}
		i.view.Editable = true
		i.view.Editor = i
		i.view.Wrap = true
		i.view.Frame = true
		i.view.Highlight = true
		if _, err := g.SetCurrentView(i.name); err != nil {
			return err
		}
	}
	i.view.Clear()
	fmt.Fprint(i.view, subDirs(i.path))
	i.view.Title = i.path
	return nil
}

func (i *cdView) Edit(view *gocui.View, key gocui.Key, char rune, mod gocui.Modifier) {
	switch {
	case char != 0 && mod == 0:
		i.jump(char)
	case key == gocui.KeyArrowDown:
		i.moveCursor(1)
	case key == gocui.KeyPgdn:
		i.moveCursor(8)
	case key == gocui.KeyArrowUp:
		i.moveCursor(-1)
	case key == gocui.KeyPgup:
		i.moveCursor(-8)
	case key == gocui.KeyArrowLeft:
		i.arrowLeft()
	case key == gocui.KeyArrowRight:
		i.arrowRight()
	}
}

func (i *cdView) keybindings(g *gocui.Gui) {
	if err := g.SetKeybinding("", gocui.KeyEsc, gocui.ModNone, quit); err != nil {
		log.Panicln(err)
	}
	if err := g.SetKeybinding("", gocui.KeyTab, gocui.ModNone, quit); err != nil {
		log.Panicln(err)
	}
	if err := g.SetKeybinding("", gocui.KeyF3, gocui.ModNone, quit); err != nil {
		log.Panicln(err)
	}
	if err := g.SetKeybinding("", gocui.KeyEnter, gocui.ModNone, i.quit); err != nil {
		log.Panicln(err)
	}
}

func (i *cdView) quit(g *gocui.Gui, v *gocui.View) error {
	_, cy := v.Cursor()
	_, oy := v.Origin()
	i.line = i.path
	if len(v.Buffer()) > 0 {
		if !strings.HasSuffix(i.line, cfg.PathSep) {
			i.line += cfg.PathSep
		}
		i.line += v.BufferLines()[cy+oy]
	}
	return gocui.ErrQuit
}

func (i *cdView) jump(char rune) {
	ny := -1
	_, my := i.view.Size()
	for i, l := range i.view.ViewBufferLines() {
		if strings.HasPrefix(l, string(char)) {
			ny = i
			break
		}
	}
	if ny == -1 {
		return
	}
	i.view.SetOrigin(0, 0)
	if ny > my {
		i.view.SetOrigin(0, ny-my+1)
		i.view.SetCursor(0, my-1)
		return
	}
	i.view.SetCursor(0, ny)
}

func (i *cdView) moveCursor(dy int) error {
	_, cy := i.view.Cursor()
	_, oy := i.view.Origin()
	switch ny := cy + oy + dy; {
	case ny < 0:
		i.view.SetCursor(0, 0)
		i.view.SetOrigin(0, 0)
	case ny < len(i.view.ViewBufferLines()):
		i.view.MoveCursor(0, dy, true)
	}
	return nil
}

func (i *cdView) arrowUp(g *gocui.Gui, v *gocui.View) error {
	_, cy := v.Cursor()
	_, oy := v.Origin()
	if cy+oy > 0 {
		v.MoveCursor(0, -1, true)
	}
	return nil
}

func (i *cdView) arrowRight() error {
	_, cy := i.view.Cursor()
	_, oy := i.view.Origin()
	if !strings.HasSuffix(i.path, cfg.PathSep) {
		i.path += cfg.PathSep
	}
	i.path += i.view.ViewBufferLines()[cy+oy]
	i.view.SetOrigin(0, 0)
	i.view.SetCursor(0, 0)
	return nil
}

func (i *cdView) arrowLeft() error {
	idx := strings.LastIndex(i.path, cfg.PathSep)
	switch idx {
	case -1:
		return nil
	case 0:
		i.path = cfg.PathSep
	default:
		i.path = i.path[:idx]
	}
	i.view.SetOrigin(0, 0)
	i.view.SetCursor(0, 0)
	return nil
}

func subDirs(path string) (res string) {
	files, err := ioutil.ReadDir(path)
	if err != nil {
		return
	}
	for _, f := range files {
		if f.IsDir() {
			res += f.Name() + "\n"
		}
	}
	return strings.TrimSuffix(res, "\n")
}
