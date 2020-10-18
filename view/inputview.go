package view

import (
	"fmt"
	"strings"
	"unicode/utf8"

	"github.com/jroimartin/gocui"

	"github.com/xenobyter/xbsh/cmd"
	"github.com/xenobyter/xbsh/storage"
)

type tInputView struct {
	name   string
	height int
	prompt string
	view   *gocui.View
	hPos   int64
}

const (
	ansiPrompt = "\033[0;32m"
	ansiStdErr = "\033[0;31m"
	ansiNormal = "\033[0m"
)

func newInputView(name string, height int) *tInputView {
	return &tInputView{name: name, height: height, hPos: 0}
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
		i.view.Autoscroll = false
		i.view.Wrap = true
		i.view.Title = i.name
		i.setPrompt()
	}
	return nil
}

// Edit implements the main editor and calls functions for keyhandling
func (i *tInputView) Edit(view *gocui.View, key gocui.Key, char rune, mod gocui.Modifier) {
	cx, cy := view.Cursor()
	_, oy := view.Origin()
	mx, _ := view.Size()
	offset := utf8.RuneCountInString(cmd.GetPrompt())
	length := utf8.RuneCountInString(i.view.BufferLines()[0])
	switch {
	case char != 0 && mod == 0:
		view.EditWrite(char)
	case key == gocui.KeySpace:
		view.EditWrite(' ')
	case key == gocui.KeyBackspace || key == gocui.KeyBackspace2:
		if cx > offset || cy > 0 || oy > 0 {
			view.EditDelete(true)
		}
	case key == gocui.KeyDelete:
		if cx >= offset || cy > 0 || oy > 0 {
			view.EditDelete(false)
		}
	case key == gocui.KeyF1:
		vHelpView.toggle()
	case key == gocui.KeyF2:
		vHistoryView.toggle()
	case key == gocui.KeyArrowLeft:
		i.cursorLeft()
	case key == gocui.KeyArrowRight:
		if cx+(mx+1)*(cy+oy) < length {
			i.view.MoveCursor(1, 0, true)
		}
	case key == gocui.KeyEnd:
		i.positionEnd()
	case key == gocui.KeyHome:
		view.SetOrigin(0, 0)
		view.SetCursor(offset, 0)
	case key == gocui.KeyArrowUp && mod == gocui.ModAlt:
		vMainView.scrollMain(-1)
	case key == gocui.KeyArrowDown && mod == gocui.ModAlt:
		vMainView.scrollMain(1)
	case key == gocui.KeyPgup:
		vMainView.scrollMain(-8)
	case key == gocui.KeyPgdn:
		vMainView.scrollMain(8)
	case key == gocui.KeyArrowUp:
		i.scrollHistory(-1)
	case key == gocui.KeyArrowDown:
		i.scrollHistory(1)
	case key == gocui.KeyTab:
		go i.tabComplete()
	case key == gocui.KeyEnter:
		cmdString := trimLine(view.BufferLines())
		fmt.Fprintln(vMainView.view, ansiPrompt+cmd.GetPrompt()+ansiNormal+cmdString)
		vMainView.print(cmd.ExecCmd(cmdString))
		go func() { i.hPos = storage.HistoryWrite(cmdString) + 1 }()
		i.view.Clear()
		i.view.SetOrigin(0, 0)
		i.setPrompt()
	}
}

func (i *tInputView) setPrompt() {
	i.prompt = cmd.GetPrompt()
	fmt.Fprintf(i.view, ansiPrompt+i.prompt+ansiNormal)
	i.view.SetCursor(utf8.RuneCountInString(i.prompt), 0)
}

func (i *tInputView) scrollHistory(inc int64) {
	var history string
	history, i.hPos = storage.HistoryRead(i.hPos + inc)
	i.view.Clear()
	i.setPrompt()
	fmt.Fprint(i.view, history)
	mx, my := i.view.Size()
	length := utf8.RuneCountInString(i.view.BufferLines()[0])
	cx, cy, oy := calculateCursor(length, mx, my)
	i.view.SetCursor(cx, cy)
	i.view.SetOrigin(0, oy)
	return
}

func (i *tInputView) cursorLeft() {
	cx, cy := i.view.Cursor()
	mx, _ := i.view.Size()
	_, oy := i.view.Origin()
	p := utf8.RuneCountInString(i.prompt)
	switch {
	case cy == 0 && oy == 0:
		if cx > p {
			i.view.MoveCursor(-1, 0, true)
		}
	case cy == 0 && cx == 0 && oy > 0:
		i.view.SetCursor(mx-10, 0)
		i.view.SetOrigin(0, oy-1)
	case cy > 0 || oy > 0:
		i.view.MoveCursor(-1, 0, true)
	}
}

func (i *tInputView) tabComplete() {
	item, path := splitCmd(i.view.Buffer())
	vCompletionView.itemLength = len(item)
	vCompletionView.completions = completionSearch(item, path)
	switch len(vCompletionView.completions) {
	case 0:
		return
	case 1:
		fmt.Fprint(i.view, vCompletionView.completions[0].Name()[vCompletionView.itemLength:])
		if vCompletionView.completions[0].IsDir() {
			fmt.Fprint(i.view, "/")
		}
		i.positionEnd()
	default:
		vCompletionView.toggle()
	}
}

func (i *tInputView) positionEnd() {
	mx, my := i.view.Size()
	l := len(i.view.Buffer()) - 1
	x, y, o := calculateCursor(l, mx, my)
	i.view.SetCursor(x, y)
	i.view.SetOrigin(0, o)
}

//trimLine cuts off any non-command parts from input
//it returns the command only
func trimLine(bufferLines []string) string {
	const sep = "$ "
	buffer := bufferLines[len(bufferLines)-1]
	if i := strings.Index(buffer, sep); i != -1 {
		buffer = buffer[i+len(sep):]
	}
	buffer = strings.TrimSpace(buffer)
	return strings.TrimSuffix(buffer, "\n")
}

//calculateCursor takes the length of a string and x/y dimensions of a view
//it returns the cursors x/y and origins y to position the cursor in the
//last line
func calculateCursor(l, mx, my int) (cx, cy, oy int) {
	lines := l / mx
	switch {
	case lines == 0:
		cx = l
	case lines < my:
		cy = lines
		cx = l - lines*mx
	case lines >= my:
		cy = my - 1
		oy = lines - my + 1
		cx = l - lines*mx
	}
	return cx, cy, oy
}

// splitCmd takes the input buffer and returns the path and the item for completion
// item is the word following the last item seperator witout any path
// path is either the workdir or an absolut path from the item
func splitCmd(cmd string) (item, path string) {
	var itemSep = []string{" "} //TODO: #63 Add itemSep to config management
	var pathSep = []string{"/"} //TODO: #64 Add pathSep to config management
	var iStart, pStart, pEnd int

	cmd = strings.TrimSuffix(cmd, "\n")
	iStart = findLastSep(cmd, itemSep) + 1
	switch {
	case cmd == "":
		return
	case iStart == len(cmd):
		//prompt only
		pStart, pEnd = strings.Index(cmd, ":")+1, strings.Index(cmd, "$ ")+1
	case cmd[iStart] == '/':
		fmt.Println("hier")
		pStart = iStart
		pEnd = findLastSep(cmd, pathSep) + 1
		iStart = pEnd
	case len(cmd) >= iStart+2 && cmd[iStart:iStart+2] == "./":
		iStart += 2
		pStart, pEnd = strings.Index(cmd, ":")+1, strings.Index(cmd, "$ ")+1
	default:
		pStart, pEnd = strings.Index(cmd, ":")+1, strings.Index(cmd, "$ ")+1
	}
	path = cmd[pStart : pEnd-1]
	item = cmd[iStart:]

	//seperate remaining paths from item
	if iStart = findLastSep(item, pathSep); iStart != -1 {
		path = path + "/" + item[:iStart]
		item = item[iStart+1:]
	}

	//fix for absolute paths starting from root
	if path == "" {
		path = "/"
	}

	return
}

func findLastSep(str string, sep []string) (max int) {
	max = -1
	for _, s := range sep {
		if i := strings.LastIndex(str, s); i > max {
			max = i
		}
	}
	return
}
