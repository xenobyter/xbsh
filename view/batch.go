package view

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strconv"
	"strings"
	"unicode/utf8"

	"github.com/awesome-gocui/gocui"
	"github.com/xenobyter/xbsh/db"
)

type batchView struct {
	name         string
	lView, rView *gocui.View
}

// Batch opens the rename view
func Batch(string) string {
	v := newBatchView("Batch")
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

func newBatchView(name string) *batchView {
	return &batchView{name: name}
}

func (i *batchView) layout(g *gocui.Gui) error {
	var (
		err error
	)
	maxX, maxY := g.Size()
	if i.lView, err = g.SetView(i.name, 0, 0, 25, maxY-2, 0); err != nil {
		if !gocui.IsUnknownView(err) {
			return err
		}
		i.lView.Wrap = false
		i.lView.Frame = true
		i.lView.Highlight = true
		i.lView.Title = i.name
		i.lView.Editable = true
		i.lView.Editor = i

		//read rules
		rules := db.ReadBatchRules()
		for _, rule := range rules {
			fmt.Fprintln(i.lView, rule)
		}

		i.lView.SetCursor(0, 0)
		if _, err := g.SetCurrentView(i.name); err != nil {
			return err
		}
	}
	if i.rView, err = g.SetView("Output", 26, 0, maxX-1, maxY-2, 0); err != nil {
		if !gocui.IsUnknownView(err) {
			return err
		}
		i.rView.Wrap = false
		i.rView.Frame = true
		i.rView.Highlight = true
		i.rView.Title = "Output"
		i.rView.Autoscroll = true
	}
	i.rView.Clear()
	i.preview(g, i.lView)
	return nil
}

func (i *batchView) keybindings(g *gocui.Gui) error {
	if err := g.SetKeybinding("", gocui.KeyEsc, gocui.ModNone, i.quit); err != nil {
		log.Panicln(err)
	}
	if err := g.SetKeybinding("", gocui.KeyF5, gocui.ModNone, i.quit); err != nil {
		log.Panicln(err)
	}
	//TODO: implement tab-switching
	//TODO: implement scrolling in results
	return nil
}

func (i *batchView) quit(g *gocui.Gui, v *gocui.View) error {
	db.WriteBatchRules(i.lView.BufferLines())
	return gocui.ErrQuit
}

func (i *batchView) preview(g *gocui.Gui, v *gocui.View) error {
	wd, err := os.Getwd()
	if err != nil {
		return err
	}
	files := preview(wd, i.lView.BufferLines())
	for _, f := range files {
		fmt.Fprintln(i.rView, f)
	}
	return nil
}

func (i *batchView) Edit(v *gocui.View, key gocui.Key, ch rune, mod gocui.Modifier) { //TODO: Implement Home and End
	ox, oy := v.Origin()
	cx, cy := v.Cursor()
	_, my := v.Size()
	switch {
	case key == 0:
		v.EditWrite(ch)
	case key == gocui.KeySpace:
		v.EditWrite(' ')
	case key == gocui.KeyEnter:
		v.EditNewLine()
		if cy == my-1 {
			v.MoveCursor(0, 1, true)
		}
	case key == gocui.KeyBackspace || key == gocui.KeyBackspace2:
		v.EditDelete(true)
	case key == gocui.KeyDelete:
		v.EditDelete(false)
	case key == gocui.KeyArrowLeft:
		if ox+cx > 0 {
			v.MoveCursor(-1, 0, true)
		}
	case key == gocui.KeyArrowRight:
		line, _ := v.Line(oy + cy)
		if ox+cx < utf8.RuneCountInString(line) {
			v.MoveCursor(1, 0, true)
		}
	case key == gocui.KeyArrowUp:
		if oy+cy > 0 {
			v.MoveCursor(0, -1, true)
		}
	case key == gocui.KeyArrowDown:
		if cy+oy < len(v.BufferLines())-1 {
			v.MoveCursor(0, 1, true)
		}
	}
}

func preview(dir string, rules []string) (out []string) {
	var longestName int
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		return
	}
	left := make([]string, len(files))
	right := make([]string, len(files))
	for i, f := range files {
		left[i] = f.Name()
		if l := utf8.RuneCountInString(f.Name()); l > longestName {
			longestName = l
		}
		right[i] = doRules(left[i], rules, i)
	}
	for i, f := range files {
		out = append(out, left[i]+strings.Repeat(" ", longestName-len(f.Name()))+" => "+right[i])

	}
	return
}

// doRules uses the following syntax
// Insert: ins $place string
// $place can be one of: "pre", "suf","pos" for a fixed position or "aft $string" to insert after a specific  substring
func doRules(name string, rules []string, cnt int) string {
	for _, r := range rules {
		fields := strings.Fields(r)

		if len(fields) < 3 {
			return name
		}
		switch fields[0] {
		case "ins":
			name = insert(fields[2], name, fields, cnt)
		}
	}
	return name
}

// Insert can be used with: "pre", "suf","pos" for a fixed position or "aft" to insert after a specific  substring
// ins string pre
// ins string suf
// ins string pos 2
// ins string aft substring
// If string does parse as int, insert increments string for each file
// For details see tests.
func insert(place, in string, fields []string, cnt int) (out string) {
	//handle increments
	if start, err := strconv.ParseInt(fields[1], 10, 0); err == nil {
		l := fmt.Sprint(len(fields[1]))
		fields[1] = fmt.Sprintf("%0"+l+"d", int(start)+cnt)
	}
	switch place {
	case "pre":
		out = fields[1] + in
	case "suf":
		out = in + fields[1]
	case "pos":
		if len(fields) < 4 {
			return in
		}
		pos, err := strconv.ParseInt(fields[3], 10, 0)
		if err != nil || int(pos) > utf8.RuneCountInString(in) || pos < 0 {
			return in
		}
		out = in[:pos] + fields[1] + in[pos:]
	case "aft":
		if len(fields) < 4 {
			return in
		}
		if pos := strings.Index(in, fields[3]); pos != -1 {
			pos += +utf8.RuneCountInString(fields[3])
			out = in[:pos] + fields[1] + in[pos:]
		} else {
			return in
		}
	default:
		return in
	}
	return
}
