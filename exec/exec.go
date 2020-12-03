/*
Package exec implements functions to execute shell commands.
*/
package exec

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"os/user"
	"path/filepath"
	"strings"

	"github.com/xenobyter/xbsh/term"
)

// Cmd takes a shell command and its arguments as a string. It returns
// the commands error
func Cmd(line string) (err error) {

	command, args, err := splitArgs(line)
	if err != nil {
		return
	}

	stdin, stdout, stderr, args, err := redirect(args)
	if err != nil {
		return err
	}

	err = run(stdin, stdout, stderr, command, args...)
	return err
}

func redirect(args []string) (stdin, stdout, stderr *os.File, resArgs []string, err error) {
	stdin, stdout, stderr = os.Stdin, os.Stdout, os.Stderr
	for i, arg := range args {
		switch {
		case strings.HasPrefix(arg, ">>!"):
			if len(arg) == 3 {
				if len(args) <= i+1 {
					err = errors.New("redirect error")
					return
				}
				args[i+1] = ">>!" + args[i+1]
				break
			}
			file, e := filepath.Abs(strings.TrimLeft(arg, ">>!"))
			if e != nil {
				err = e
				break
			}
			if stderr, err = os.OpenFile(file, os.O_CREATE|os.O_APPEND|os.O_RDWR, os.ModePerm); err != nil {
				stderr = os.Stderr
			}
		case strings.HasPrefix(arg, ">!"):
			if len(arg) == 2 {
				if len(args) <= i+1 {
					err = errors.New("redirect error")
					return
				}
				args[i+1] = ">!" + args[i+1]
				break
			}
			file, e := filepath.Abs(strings.TrimLeft(arg, ">!"))
			if e != nil {
				err = e
				break
			}
			if stderr, err = os.Create(file); err != nil {
				stderr = os.Stderr
			}
		case strings.HasPrefix(arg, ">>"):
			if len(arg) == 2 {
				if len(args) <= i+1 {
					err = errors.New("redirect error")
					return
				}
				args[i+1] = ">>" + args[i+1]
				break
			}
			file, e := filepath.Abs(strings.TrimLeft(arg, ">>"))
			if e != nil {
				err = e
				break
			}
			if stdout, err = os.OpenFile(file, os.O_CREATE|os.O_APPEND|os.O_RDWR, os.ModePerm); err != nil {
				stdout = os.Stdout
			}
		case strings.HasPrefix(arg, ">"):
			if len(arg) == 1 {
				if len(args) <= i+1 {
					err = errors.New("redirect error")
					return
				}
				args[i+1] = ">" + args[i+1]
				break
			}
			file, e := filepath.Abs(strings.TrimLeft(arg, ">"))
			if e != nil {
				err = e
				break
			}
			if stdout, err = os.Create(file); err != nil {
				stdout = os.Stdout
			}
		case strings.HasPrefix(arg, "<"):
			if len(arg) == 1 {
				if len(args) <= i+1 {
					err = errors.New("redirect error")
					return
				}
				args[i+1] = "<" + args[i+1]
				break
			}
			file, e := filepath.Abs(strings.TrimLeft(arg, "<"))
			if e != nil {
				err = e
				break
			}
			if stdin, err = os.Open(file); err != nil {
				stdin = os.Stdin
			}
		default:
			resArgs = append(resArgs, arg)
		}
	}
	return stdin, stdout, stderr, resArgs, err
}

func run(stdin, stdout, stderr *os.File, command string, args ...string) (err error) {
	//handle the command
	switch command {
	case "bg":
		j := BgJob(args)
		fmt.Printf("background job %v started\n", j)
	case "cd":
		err = changeDir(args)
	default:
		cmd := exec.Command(command, args...)
		//run it
		cmd.Stdin = stdin
		cmd.Stdout = stdout
		pipe, _ := cmd.StderrPipe()
		err = cmd.Start()
		out, _ := ioutil.ReadAll(pipe)
		if len(out) > 0 {
			fmt.Fprint(stderr, term.AnsiErr, string(out), term.AnsiNormal)
		}
		err = cmd.Wait()
	}
	return err
}

func splitArgs(line string) (command string, expArgs []string, err error) {
	var args []string
	if line == "" {
		err = errors.New("errNoCommand")
	} else {
		fields := strings.Fields(line)
		command = fields[0]
		args = strings.Fields(line)[1:]
	}
	expArgs = expandArg(args)
	return
}

func changeDir(args []string) (err error) {
	//cd to homedir when called without argument
	if len(args) == 0 {
		usr, _ := user.Current()
		args = append(args, usr.HomeDir)
	}
	//actually change dir
	if err = os.Chdir(args[0]); err != nil {
		fmt.Fprintln(os.Stderr, err)
		return
	}
	dir, _ := os.Getwd()
	fmt.Println(dir)

	return
}
