package cmd

import (
	"io/ioutil"
	"os"
	"os/exec"
	"strings"
)

// ExecCmd runs line and returns its output
func ExecCmd(line string) (stdout, stderr []byte) {
	//store normal stdout/stderr
	oldStdOut := os.Stdout
	oldStdErr := os.Stderr

	//setup the Pipe to capture output
	rStdOut, wStdOut, _ := os.Pipe()
	os.Stdout = wStdOut
	rStdErr, wStdErr, _ := os.Pipe()
	os.Stderr = wStdErr

	command, args := parseCmd(line) //TODO: #13 handle empty command string
	//handle the command
	switch command {
	case "cd":
		setWorkDir(args[0])
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

func parseCmd(line string) (command string, args []string) {
	fields := strings.Fields(line)
	command = fields[0]
	args = strings.Fields(line)[1:]
	return
}
