package cmd

import (
	"io/ioutil"
	"os"
	"os/exec"
	
)

// ParseCmd takes a string, handles any processing and runs it
func ParseCmd(cmdString string) []byte {
	//store normal stdout
	oldStdOut := os.Stdout

	//setup the Pipe to capture output
	read, write, _ := os.Pipe()
	os.Stdout = write
	
	//handle the command 
	cmd := exec.Command(cmdString)

	//run it
	cmd.Stderr = os.Stderr
	cmd.Stdout = write
	cmd.Run()

	//finish
	write.Close()
	out, _ := ioutil.ReadAll(read)
	os.Stdout = oldStdOut

	return out
}
