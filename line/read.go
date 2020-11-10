/*
Package line implements the basic command line editor
*/
package line

import (
	"fmt"
	"unicode/utf8"

	"github.com/eiannone/keyboard"

	"github.com/xenobyter/xbsh/cfg"
	"github.com/xenobyter/xbsh/db"
	"github.com/xenobyter/xbsh/term"
	"github.com/xenobyter/xbsh/view"
)

type mockInput struct {
	runes []rune
	keys  []keyboard.Key
	errs  []error
	i     int
}

var (
	mockIn *mockInput
	hist   *history
)

// Mockservices for testing only
func newMockInput() *mockInput {
	return &mockInput{i: 0}
}
func (i *mockInput) next() (rune, keyboard.Key, error) {
	rune := i.runes[i.i]
	key := i.keys[i.i]
	err := i.errs[i.i]
	i.i++
	return rune, key, err
}
func (i *mockInput) set(runes []rune, keys []keyboard.Key, errs []error) {
	i.runes = runes
	i.keys = keys
	i.errs = errs
	i.i = 0
}

func insert(line string, pos int, c rune) (res string) {
	if pos == utf8.RuneCountInString(line) {
		res = line + string(c)
	} else {
		pre := []rune(line)[0:pos]
		suf := []rune(line)[pos:]
		runes := append(pre, c)
		res = string(append(runes, suf...))
	}
	return
}

func delete(line string, pos int) (res string) {
	pre := []rune(line)[0:pos]
	suf := []rune(line)[pos+1:]
	return string(append(pre, suf...))
}

func output(prompt, line string, lx, ox, ld int) (res string, cx int) {
	lp := utf8.RuneCountInString(prompt)
	runes := []rune(line)
	start := ox
	end := ox + ld
	if end > len(runes) {
		end = len(runes)
	}
	return string(runes[start:end]), lp + lx - ox
}

func move(dx, lx, ox, ll, ld int) (nx, no int) {
	//1. move lx
	nx = lx + dx
	//2. respect line borders
	switch {
	case nx > ll:
		nx = ll
	case nx < 0:
		nx = 0
	}
	//3. position ox
	switch {
	case nx-ox > ld:
		no = nx - ld
	case nx < ox:
		no = nx
	default:
		no = ox
	}
	return
}

// Read one line until <Enter> and returns the finished line.
// It handles F-Keys.
func Read(prompt string, test ...mockInput) (line string) {
	var (
		getKey func() (rune, keyboard.Key, error)
		ws     *term.WinSize
		err    error
		lx     int //cursor in line
		ox     int //origin in line
		ld     int //length of display
	)

	if test == nil {
		//setup real input
		ws, err = term.GetWinsize()
		if err != nil {
			panic(err)
		}
		if err := keyboard.Open(); err != nil {
			panic(err)
		}
		defer keyboard.Close()
		getKey = keyboard.GetKey
	} else {
		//setup testing
		ws = &term.WinSize{Row: 10, Col: 100, Xpixel: 0, Ypixel: 0}
		mockIn = newMockInput()
		mockIn.set(test[0].runes, test[0].keys, test[0].errs)
		getKey = mockIn.next
	}

	hist = newHistory()

	ld = int(ws.Col) - utf8.RuneCountInString(prompt) - 1
	fmt.Printf(term.AnsiPrompt + prompt + term.AnsiNormal)
	for {
		c, k, err := getKey()
		ll := utf8.RuneCountInString(line) //don't use when editing line
		switch err == nil {
		//Insert
		case k == 0:
			line = insert(line, lx, c)
			lx, ox = move(1, lx, ox, utf8.RuneCountInString(line), ld)
		case k == keyboard.KeySpace:
			line = insert(line, lx, ' ')
			lx, ox = move(1, lx, ox, utf8.RuneCountInString(line), ld)
		//Movement
		case k == keyboard.KeyArrowLeft:
			if lx > 0 {
				lx, ox = move(-1, lx, ox, ll, ld)
			}
		case k == keyboard.KeyArrowRight:
			if lx < utf8.RuneCountInString(line) {
				lx, ox = move(1, lx, ox, ll, ld)
			}
		case k == keyboard.KeyHome:
			lx, ox = move(-lx, lx, ox, ll, ld)
		case k == keyboard.KeyEnd:
			lx, ox = move(ll, lx, ox, ll, ld)
			//TODO: #74 Jump words
		//Delete
		case k == keyboard.KeyBackspace || k == keyboard.KeyBackspace2:
			if lx > 0 {
				lx, ox = move(-1, lx, ox, utf8.RuneCountInString(line), ld)
				line = delete(line, lx)
			}
		case k == keyboard.KeyDelete:
			if lx < utf8.RuneCountInString(line) {
				line = delete(line, lx)
			}
		//Discard
		case k == keyboard.KeyEsc || k == keyboard.KeyCtrlC:
			line, lx, ox = "", 0, 0
			fmt.Println()
		case k == keyboard.KeyCtrlL:
			line, lx, ox = "", 0, 0
			fmt.Print(term.AnsiCls)
		//history
		case k == keyboard.KeyArrowUp:
			if hist.id == 0 {
				hist.pending = line
			}
			line = hist.scroll(-1)
			lx, ox = move(len(line), lx, ox, utf8.RuneCountInString(line), ld)
		case k == keyboard.KeyArrowDown:
			if hist.id > 0 {
				line = hist.scroll(1)
				lx, ox = move(len(line), lx, ox, utf8.RuneCountInString(line), ld)
			}
		case k == keyboard.KeyF2:
			line = openView(view.History, line)
			lx, ox = move(len(line), lx, ox, utf8.RuneCountInString(line), ld)
			//Help
		case k == keyboard.KeyF1:
			openView(view.Help, "")
			//Complete
		case k == keyboard.KeyTab:
			line = tabComplete(line)
			lx, ox = move(len(line), lx, ox, utf8.RuneCountInString(line), ld)
		//Enter
		case k == keyboard.KeyEnter:
			fmt.Println()
			go func(){
				hist.id = db.HistoryWrite(line)
				hist.pending = ""
				db.CleanUp(cfg.HistoryMaxEntries, cfg.HistoryDelExitCmd)
			}()
			return
		}
		out, dx := output(prompt, line, lx, ox, ld)
		p := term.AnsiClrLine + term.AnsiPrompt + prompt + term.AnsiNormal
		fmt.Printf("%v%v\033[0G\033[%vC", p, out, dx)
	}
}

func openView(f func(string) string, line string) string {
	keyboard.Close()
	defer keyboard.Open()
	return f(line)
}
