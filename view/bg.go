package view

import (
	"fmt"
	"log"

	"github.com/awesome-gocui/gocui"
	"github.com/xenobyter/xbsh/exec"
	"github.com/xenobyter/xbsh/term"
)

type bgView struct {
	name         string
	lView, rView *gocui.View
}

//BG opens the background jobs view
func BG(string) string {
	term.SetStatus("| ESC: Exit | F4: Exit | Enter: Exit | Del: Delete | Tab: Left<->Right")
	v := newBgView("Jobs")
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
	return ""
}

func newBgView(name string) *bgView {
	return &bgView{name: name}
}

func (i *bgView) layout(g *gocui.Gui) error {
	var (
		err error
	)
	maxX, maxY := g.Size()
	if i.lView, err = g.SetView(i.name, 0, 0, 13, maxY-2, 0); err != nil {
		if !gocui.IsUnknownView(err) {
			return err
		}
		i.lView.Wrap = false
		i.lView.Frame = true
		i.lView.Highlight = true
		i.lView.Title = i.name

		i.lView.SetCursor(0, 0)
		if _, err := g.SetCurrentView(i.name); err != nil {
			return err
		}

		outList(i.lView)
	}
	if i.rView, err = g.SetView("Output", 14, 0, maxX-1, maxY-2, 0); err != nil {
		if !gocui.IsUnknownView(err) {
			return err
		}
		i.rView.Wrap = true
		i.rView.Frame = true
		i.rView.Highlight = true
		i.rView.Title = "Output"
	}
	outSingleJob(i.lView, i.rView)
	return nil
}

func (i *bgView) keybindings(g *gocui.Gui) {
	if err := g.SetKeybinding("", gocui.KeyEsc, gocui.ModNone, quit); err != nil {
		log.Panicln(err)
	}
	if err := g.SetKeybinding("", gocui.KeyF4, gocui.ModNone, quit); err != nil {
		log.Panicln(err)
	}
	if err := g.SetKeybinding("", gocui.KeyEnter, gocui.ModNone, quit); err != nil {
		log.Panicln(err)
	}
	if err := g.SetKeybinding("", gocui.KeyArrowDown, gocui.ModNone, arrowDown); err != nil {
		log.Panicln(err)
	}
	if err := g.SetKeybinding("", gocui.KeyArrowUp, gocui.ModNone, arrowUp); err != nil {
		log.Panicln(err)
	}
	if err := g.SetKeybinding("", gocui.KeyTab, gocui.ModNone, tab); err != nil {
		log.Panicln(err)
	}
	if err := g.SetKeybinding("Jobs", gocui.KeyDelete, gocui.ModNone, deleteJob); err != nil {
		log.Panicln(err)
	}
}

func outList(v *gocui.View) {
	v.SetCursor(0, 0)
	v.SetOrigin(0, 0)
	v.Clear()
	cmds := exec.BgList()
	for id, cmd := range cmds {
		fmt.Fprintf(v, "%2d: %v\n", id, cmd)
	}
}

func outSingleJob(lv, rv *gocui.View) {
	cy, oy := getY(lv)
	cmd, args, stdout, stderr, finished := exec.BgGet(cy + oy)
	rv.Clear()
	if finished {
		fmt.Fprint(rv, "Finished: ")
	} else {
		fmt.Fprint(rv, "Running: ")
	}
	fmt.Fprintf(rv, "%v ", cmd)
	for _, arg := range args {
		fmt.Fprintf(rv, "%v ", arg)
	}
	fmt.Fprintln(rv)
	rv.Write(stdout)
	fmt.Fprint(rv, term.AnsiErr+string(stderr)+term.AnsiNormal)
}

func deleteJob(g *gocui.Gui, v *gocui.View) error {
	cy, oy := getY(v)
	exec.BgDelete(cy + oy)
	outList(v)
	return nil
}

func tab(g *gocui.Gui, v *gocui.View) error {
	views := g.Views()
	if v == views[0] {
		g.SetCurrentView(views[1].Name())
	} else {
		g.SetCurrentView(views[0].Name())
		views[1].SetOrigin(0, 0)
		views[1].SetCursor(0, 0)
	}
	return nil
}

func arrowDown(g *gocui.Gui, v *gocui.View) error {
	cy, oy := getY(v)
	if cy+oy < len(v.ViewBufferLines())-2 {
		v.MoveCursor(0, 1, true)
	}
	return nil
}

func arrowUp(g *gocui.Gui, v *gocui.View) error {
	cy, oy := getY(v)
	if cy+oy > 0 {
		v.MoveCursor(0, -1, true)
	}
	return nil
}

func getY(v *gocui.View) (cy, oy int) {
	_, cy = v.Cursor()
	_, oy = v.Origin()
	return
}
