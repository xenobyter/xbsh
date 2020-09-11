package cmd

import (
	"errors"
	"io/ioutil"
	"os"
	"os/exec"
	"strings"
)

// ExecCmd runs line and returns its output
func ExecCmd(line string) (stdout, stderr []byte) {
	command, args, err := parseCmd(line) 
	if err != nil {
		return
	}

	//store normal stdout/stderr
	oldStdOut := os.Stdout
	oldStdErr := os.Stderr

	//setup the Pipe to capture output
	rStdOut, wStdOut, _ := os.Pipe()
	os.Stdout = wStdOut
	rStdErr, wStdErr, _ := os.Pipe()
	os.Stderr = wStdErr

	//handle the command
	switch command {
	case "cd": //ToDo: #26 Fix panic on cd without arguments
		if err := os.Chdir(args[0]); err!=nil {
			stderr = []byte(err.Error() + "\n")
			return
		}
	default:
		cmd := exec.Command(command, args...)
		//run it
		cmd.Stderr = wStdErr
		cmd.Stdout = wStdOut
		cmd.Run()
	}

	//finish
	wStdOut.Close()
	wStdErr.Close()
	stdout, _ = ioutil.ReadAll(rStdOut)
	os.Stdout = oldStdOut
	stderr, _ = ioutil.ReadAll(rStdErr)
	os.Stderr = oldStdErr

	return
}

func parseCmd(line string) (command string, args []string, err error) {
	if line == "" {
		err = errors.New("errNoCommand")
	} else {
		fields := strings.Fields(line)
		command = fields[0]
		args = strings.Fields(line)[1:]
	}
	return
}
