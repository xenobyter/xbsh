package cmd

import (
	"strings"
	"io/ioutil"
	"os"
	"os/exec"
	
)

// ExecCmd runs line and returns its output
func ExecCmd(line string) []byte {
	//store normal stdout
	oldStdOut := os.Stdout

	//setup the Pipe to capture output
	read, write, _ := os.Pipe()
	os.Stdout = write
	
	command, args := parseCmd(line)
	//handle the command 
	cmd := exec.Command(command, args...)

	//run it
	cmd.Stderr = os.Stderr //TODO: #10 redirect Stderr
	cmd.Stdout = write
	cmd.Run()

	//finish
	write.Close()
	out, _ := ioutil.ReadAll(read)
	os.Stdout = oldStdOut

	return out
}

func parseCmd(line string) (command string, args []string) {
	fields := strings.Fields(line)
	command = fields[0]
	args = strings.Fields(line)[1:]
	return
}