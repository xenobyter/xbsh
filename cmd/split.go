package cmd

import (
	"strings"
)

var itemSep = " "  //TODO: #63 Add itemSep to config management
var groupSep = []rune{'|'} //TODO: #69 Add groupSep to config management
var pathSep = "/"  //TODO: #70 Add pathSep to config management

// LastItem takes the whole buffer from InputView and returns the
// item to be possibly completed and it's path. If the item is the
// first item in a cmd group and without path, bin is true. 
// In this case files from PATH can be completed.
func LastItem(buffer string) (cmd, path string, bin bool) {
	bin=true

	//Trim trainling \n
	buffer = strings.TrimSuffix(buffer, "\n")

	path, buffer = seperatePrompt(buffer)
	if strings.TrimSpace(buffer) == "" {
		return
	}

	groups := splitBuffer(buffer, groupSep)
	items := strings.Split(groups[len(groups)-1], itemSep)
	item := items[len(items)-1]

	//cmd is absolute path, discard workdir
	if strings.HasPrefix(item, pathSep) {
		bin=false
		path = ""
	}

	//handle paths inside
	parts := strings.Split(item, pathSep)
	for i, p := range parts {
		if p!="" && p != "." && i < len(parts)-1 {
			bin=false
			path += pathSep + p
		}
	}
	//cmd is absolute path, and path still empty
	if strings.HasPrefix(item, pathSep) && path == "" {
		path = pathSep
	}
	cmd = parts[len(parts)-1]

	//buffer ends with blank, thus cmd is empty
	if strings.HasSuffix(buffer, " ") {
		cmd = ""
		bin = false
		return
	}

	if strings.HasSuffix(cmd, pathSep) {
		path += pathSep + strings.TrimSuffix(cmd, pathSep)
		cmd = ""
	}

	//Only the first in a group can be from bin
	if len(items) > 1 {
		bin = false
	}

	return
}

func splitBuffer(buffer string, sep []rune) []string {
	f := func(c rune) bool {
		for _, r := range sep {
			if r == c {
				return true
			}
		}
		return false
	}
	return strings.FieldsFunc(buffer, f)
}

func seperatePrompt(buffer string) (path, cmd string) {
	pStart := strings.Index(buffer, ":")
	pEnd := strings.Index(buffer, "$ ")
	if pStart != -1 && pEnd > pStart {
		path = buffer[pStart+1 : pEnd]
		cmd = buffer[pEnd+2:]
	}
	return
}
