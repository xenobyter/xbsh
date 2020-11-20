package term

import (
	"fmt"
	"strings"
	"time"
)

var (
	ws     *WinSize
	ticker *time.Ticker
)

//StatusLine displays the bottom line. It is intended to run
//in it's own go routine.
func StatusLine() {
	ticker = time.NewTicker(time.Second)
	done := make(chan bool)
	for {
		select {
		case <-done:
			return
		case t := <-ticker.C:
			ws, _ = GetWinsize()
			bar := getStatusBarContent(t, ws.Col)
			if bar == "" {
				return
			}
			fmt.Printf(AnsiSave+"\033[%v;0H%v"+AnsiRestore, ws.Col,bar)
		}
	}
}

//PauseStatus is called to pause the ticker
func PauseStatus() {
	ticker.Stop()
}

//RestartStatus is called to start over the ticker
func RestartStatus() {
	ticker.Reset(time.Second)
}

func getStatusBarContent(t time.Time, cLen uint16) string {
	bar := "| F1: Help | F2: Hist | F3: cd | F4: Jobs"
	date := t.Format("15:04:05") + " |"
	fill := int(cLen) - len(bar) - len(date)
	if fill < 0 {
		return ""
	}
	return AnsiBGGreen + AnsiBlack + bar + strings.Repeat(" ", fill) + date + AnsiNormal
}
