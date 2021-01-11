package view

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"time"
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
		i.rView.Autoscroll = false
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
	if err := g.SetKeybinding("", gocui.KeyTab, gocui.ModNone, i.tabSwitch); err != nil {
		log.Panicln(err)
	}
	if err := g.SetKeybinding("Output", gocui.KeyArrowDown, gocui.ModNone, i.arrowDown); err != nil {
		log.Panicln(err)
	}
	if err := g.SetKeybinding("Output", gocui.KeyArrowUp, gocui.ModNone, i.arrowUp); err != nil {
		log.Panicln(err)
	}
	if err := g.SetKeybinding("", gocui.KeyF6, gocui.ModNone, i.rename); err != nil {
		log.Panicln(err)
	}
	return nil
}

func (i *batchView) arrowDown(g *gocui.Gui, v *gocui.View) error {
	_, cy := v.Cursor()
	_, oy := v.Origin()
	if cy+oy < len(v.ViewBufferLines())-1 {
		v.MoveCursor(0, 1, true)
	}
	return nil
}
func (i *batchView) arrowUp(g *gocui.Gui, v *gocui.View) error {
	_, cy := v.Cursor()
	_, oy := v.Origin()
	if cy+oy > 0 {
		v.MoveCursor(0, -1, true)
	}
	return nil
}

func (i *batchView) tabSwitch(g *gocui.Gui, v *gocui.View) error {
	if v == i.rView {
		g.SetCurrentView(i.name)
	} else {
		g.SetCurrentView("Output")
	}
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

func (i *batchView) rename(g *gocui.Gui, v *gocui.View) error {
	wd, err := os.Getwd()
	if err != nil {
		return err
	}
	files := rename(wd, i.lView.BufferLines())
	for _, f := range files {
		fmt.Fprintln(i.rView, f)
	}
	return nil
}

func (i *batchView) Edit(v *gocui.View, key gocui.Key, ch rune, mod gocui.Modifier) {
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
	case key == gocui.KeyEnd:
		line, err := v.Line(cy)
		if err != nil {
			break
		}
		v.SetCursor(utf8.RuneCountInString(line), cy)
	case key == gocui.KeyHome:
		v.SetCursor(0, cy)
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
		right[i] = doRules(left[i], rules, i, f.ModTime(), f.IsDir())
	}
	for i, f := range files {
		out = append(out, left[i]+strings.Repeat(" ", longestName-len(f.Name()))+" => "+right[i])

	}
	return
}

func rename(dir string, rules []string) (out []string) {
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		return
	}
	for i, f := range files {
		new := doRules(f.Name(), rules, i, f.ModTime(), f.IsDir())
		os.Rename(filepath.Join(dir,f.Name()),filepath.Join(dir,new))
		out=append(out, new)
	}
	return
}

// doRules uses the following syntax: function [args]
// function can be one of
// "ins" for insert, "del" for delete, "rep" for replace
func doRules(name string, rules []string, cnt int, t time.Time, isDir bool) string {
	include := regexp.MustCompile(".*")
	exclude := regexp.MustCompile("!.*")
	mode := "a"
	for _, r := range rules {
		fields := strings.Fields(r)
		l := len(fields)
		if l == 0 {
			return name
		}
		f := fields[0]
		switch {
		case f == "mod":
			mode = mod(fields)
		case isDir && mode == "f":
			break
		case !isDir && mode == "d":
			break
		case f == "ins" && include.MatchString(name) && !exclude.MatchString(name):
			if l < 3 {
				return name
			}
			name = ins(fields[2], name, fields, cnt)
		case f == "del" && include.MatchString(name) && !exclude.MatchString(name):
			name = del(name, fields)
		case f == "rep" && include.MatchString(name) && !exclude.MatchString(name):
			if l < 3 {
				return name
			}
			name = rep(name, fields, cnt)
		case f == "inc":
			include = inc(fields)
		case f == "exc":
			exclude = exc(fields)
		case f == "cas":
			name = cas(name, fields)
		case f == "dat":
			name = dat(name, fields, t)
		}
	}
	return name
}

// ins can be used with: "pre", "suf","pos" for a fixed position or "aft" to insert after a specific  substring.
// ins string pre
// ins string suf
// ins string pos 2
// ins string aft substring
// If string does parse as int, ins increments string for each file.
// See tests fpr details.
func ins(place, in string, fields []string, cnt int) (out string) {
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

// del can be used to delete a string.
// del string (first occurence)
// del string pre
// del string suf
// del string any (any occurence)
// del 1 (from position 1 to end of filename)
// del 1 3 (1 to 3)
// del -3 (last 3 runes)
func del(name string, fields []string) string {
	switch len(fields) {
	case 2:
		subPos := strings.Index(name, fields[1])
		subLen := utf8.RuneCountInString(fields[1])
		from, err := strconv.ParseInt(fields[1], 10, 0)
		switch {
		case subPos > -1:
			return name[:subPos] + name[subPos+subLen:]
		case err == nil && from > 0:
			return name[:from-1]
		case err == nil && from < 0:
			offset := from + int64(utf8.RuneCountInString(name))
			if offset > 0 {
				return name[:offset]
			}
		}
	case 3:
		from, e1 := strconv.ParseInt(fields[1], 10, 0)
		to, e2 := strconv.ParseInt(fields[2], 10, 0)
		switch {
		case fields[2] == "pre":
			return strings.TrimPrefix(name, fields[1])
		case fields[2] == "suf":
			return strings.TrimSuffix(name, fields[1])
		case fields[2] == "any":
			return strings.ReplaceAll(name, fields[1], "")
		case e1 != nil || e2 != nil:
			return name
		case from > 0 && to > from && to <= int64(utf8.RuneCountInString(name)):
			return name[:from-1] + name[to:]
		}
	}
	return name
}

// rep is used to replace substrings
// rep oldStr newStr (first occurence)
// rep oldStr newStr pre
// rep oldStr newStr suf
// rep oldStr newStr any
// If newStr parses as Integer, rep will increment it for each file
func rep(name string, fields []string, cnt int) string {
	l := len(fields)
	if l >= 3 {
		if start, e := strconv.ParseInt(fields[2], 10, 0); e == nil {
			l := fmt.Sprint(len(fields[2]))
			fields[2] = fmt.Sprintf("%0"+l+"d", int(start)+cnt)
		}
	}
	switch l {
	case 3:
		return strings.Replace(name, fields[1], fields[2], 1)
	case 4:
		switch fields[3] {
		case "pre":
			if strings.HasPrefix(name, fields[1]) {
				return strings.Replace(name, fields[1], fields[2], 1)
			}
		case "suf":
			if strings.HasSuffix(name, fields[1]) {
				return strings.TrimSuffix(name, fields[1]) + fields[2]
			}
		case "any":
			return strings.ReplaceAll(name, fields[1], fields[2])
		}
	}
	return name
}

func inc(fields []string) (reg *regexp.Regexp) {
	reg = regexp.MustCompile(".*")
	if len(fields) < 2 {
		return
	}
	reg, _ = regexp.Compile(fields[1])
	return
}

func exc(fields []string) (reg *regexp.Regexp) {
	reg = regexp.MustCompile("!.*")
	if len(fields) < 2 {
		return
	}
	reg, _ = regexp.Compile(fields[1])
	return
}

// cas is used to change cases between upper and lower
// cas upp (first letter)
// cas upp wrd (every word)
// cas upp any [string]
// cas low (first letter)
// cas low wrd
// cas low any [string]
func cas(name string, fields []string) string {
	l := len(fields)
	if l < 2 || name == "" {
		return name
	}
	f := fields[1]
	switch {
	case f == "upp":
		switch {
		case l == 2:
			r := strings.Split(name, "")
			return strings.ToUpper(r[0]) + strings.Join(r[1:], "")
		case l == 3 && fields[2] == "wrd":
			return strings.Title(name)
		case l == 3 && fields[2] == "any":
			return strings.ToUpper(name)
		case l == 4 && fields[2] == "any":
			return strings.ReplaceAll(name, fields[3], strings.ToUpper(fields[3]))
		}
	case f == "low":
		switch {
		case l == 2:
			r := strings.Split(name, "")
			return strings.ToLower(r[0]) + strings.Join(r[1:], "")
		case l == 3 && fields[2] == "wrd":
			res := ""
			for _, w := range strings.Fields(name) {
				r := strings.Split(w, "")
				res += strings.ToLower(string(r[0])) + strings.Join(r[1:], "") + " "
			}
			return strings.TrimSpace(res)
		case l == 3 && fields[2] == "any":
			return strings.ToLower(name)
		case l == 4 && fields[2] == "any":
			return strings.ReplaceAll(name, fields[3], strings.ToLower(fields[3]))
		}
	default:
		return name
	}
	return name
}

// dat is used to add the a date as part of the filename
// dat (use date as name)
// dat pre (use date as prefix)
// dat suf (use date as suffix)
func dat(name string, fields []string, t time.Time) string {
	l := len(fields)
	switch {
	case l == 1:
		return t.Local().Format("2006-01-02 15:04:05")
	case l == 2 && fields[1] == "pre":
		return t.Local().Format("2006-01-02 15:04:05") + " " + name
	case l == 2 && fields[1] == "suf":
		return name + " " + t.Local().Format("2006-01-02 15:04:05")
	}
	return name

}

// mod is used to restrict rules to directorys or files
// mod dir
// mod fil
// mod all
func mod(fields []string) string {
	l := len(fields)
	switch {
	case l == 2 && fields[1] == "fil":
		return "f"
	case l == 2 && fields[1] == "dir":
		return "d"
	default:
		return "a"
	}
}
