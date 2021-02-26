package main

import (
	"fmt"
	"os"

	"github.com/xenobyter/xbsh/db"
	"github.com/xenobyter/xbsh/exec"
	"github.com/xenobyter/xbsh/line"
	"github.com/xenobyter/xbsh/term"
)

func main() {
	go term.StatusLine()
	go db.PathCache()

	//setup scroll region
	ws, _ := term.GetWinsize()
	fmt.Printf("\033[2J\033[0;%vr", ws.Row-1)
	term.ClearStatus()
	for {
		p := line.Prompt()
		line := line.Read(p)

		switch line {
		case "exit":
			os.Exit(0)
		case "cls":
			fmt.Print(term.AnsiCls)
		case "top":
			term.PauseStatus()
			exec.Cmd(line)
			fmt.Print(term.AnsiCls)
			term.RestartStatus()
		default:
			term.PauseStatus()
			exec.Cmd(line)
			term.RestartStatus()
		}
	}
}
