package view

import (
	"fmt"
	"strings"
	"time"
	"unicode/utf8"

	"github.com/jroimartin/gocui"
)

var (
	vStatus *gocui.View
)

// StatusBarView represents our single status bar
type tStatusBarView struct {
	name   string
}

// StatusBar returns a pointer to status bar
func vStatusBar(name string) *tStatusBarView {
	return &tStatusBarView{name: name}
}

// Layout implements the layout for Statusbar
func (i *tStatusBarView) Layout(gui *gocui.Gui) error {
	var err error
	maxX, maxY := gui.Size()

	vStatus, err = gui.SetView(i.name, 0, maxY-2, maxX, maxY)

	if err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		vStatus.Frame = false
	}
	return nil
}

func updateStatus(gui *gocui.Gui) {
	ticker := time.NewTicker(1000 * time.Millisecond)
	done := make(chan bool)
	for {
		select {
		case <-done:
			return
		case t := <-ticker.C:

			gui.Update(func(g *gocui.Gui) error {
				v, err := g.View("StatusBar")
				if err != nil {
					return err
				}

				v.Clear()
				vWidth, _ := v.Size()
				fmt.Fprintf(v, getStatusBarContent(t, vWidth-1))
				return nil
			})
		}
	}

}

func getStatusBarContent(t time.Time, contentLen int) (sContent string) {
	var (
		date   string
		fillUp int
	)
	sContent = "❱❱❱ Help: ^h"
	switch f := contentLen - utf8.RuneCountInString(sContent); {
	case f > 23:
		date = t.Format("2006-01-02 15:04:05")
		fillUp = f - 23
	case f > 14:
		date = t.Format("2006-01-02")
		fillUp = f - 14
	case f > 4:
		date = t.Format("2006-01-02")
		fillUp = f - 14
	}
	sContent += strings.Repeat(" ", fillUp) + date + " ❰❰❰"
	return
}
