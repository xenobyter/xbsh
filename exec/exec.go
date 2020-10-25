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
	"strings"
)

const (
	ansiPrompt = "\033[0;32m"
	ansiNormal = "\033[0m"
	ansiErr    = "\033[0;31m"
)

// Cmd takes a shell command and its arguments as a string. It returns
// the commands error
func Cmd(line string) (err error) {

	command, args, err := splitArgs(line)
	if err != nil {
		return
	}

	//store normal stderr
	oldStdErr := os.Stderr

	//setup the Pipe to capture stderr
	rStdErr, wStdErr, _ := os.Pipe()
	os.Stderr = wStdErr

	//handle the command
	switch command {
	case "cd":
		err = changeDir(args)
	default:
		cmd := exec.Command(command, args...)
		//run it
		cmd.Stdin = os.Stdin
		cmd.Stderr = wStdErr
		cmd.Stdout = os.Stdin
		cmd.Start()
		err = cmd.Wait()
	}

	//finish
	wStdErr.Close()
	stderr, _ := ioutil.ReadAll(rStdErr)
	os.Stderr = oldStdErr

	if len(stderr) > 0 {
		fmt.Fprint(os.Stderr, ansiErr, string(stderr), ansiNormal)
	}
	return
}

func splitArgs(line string) (command string, args []string, err error) {
	if line == "" {
		err = errors.New("errNoCommand")
	} else {
		fields := strings.Fields(line)
		command = fields[0]
		args = strings.Fields(line)[1:]
	}
	return
}

func changeDir(args []string) (err error) {
	//cd to homedir when called without argument
	if len(args) == 0 {
		usr, _ := user.Current()
		args = append(args, usr.HomeDir)
	}
	//acutally change dir
	if err = os.Chdir(args[0]); err != nil {
		fmt.Fprintln(os.Stderr, err)
		return
	}
	dir, _ := os.Getwd()
	fmt.Println(dir)

	return
}
