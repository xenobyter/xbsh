package line

import (
	"io/ioutil"
	"os"
	"strings"
	"unicode/utf8"

	"github.com/eiannone/keyboard"
	"github.com/xenobyter/xbsh/db"
	"github.com/xenobyter/xbsh/view"
)

func tabComplete(line string) (res string) {
	item, path, bin := lastItem(line)

	//handle cd
	if line == "cd" || line == "cd " {
		keyboard.Close()
		defer keyboard.Open()
		p := view.CD(path)
		if p != "" {
			return "cd " + p
		}
	}

	completions := completionSearch(item, path)
	if bin {
		completions = append(completions, db.PathComplete(item)...)
	}
	switch len(completions) {
	case 0:
		return line
	case 1:
		res = line + completions[0].Name()[utf8.RuneCountInString(item):]
		if completions[0].IsDir() {
			res += "/"
		}
	default:
		keyboard.Close()
		defer keyboard.Open()
		res = line
		if c := view.Completion(completions); len(c) > 0 {
			res += c[utf8.RuneCountInString(item):]
		}
	}
	return
}

//completionSearch takes search term and directory. It returns a slice of files and directories
//the names of wich start with search term
func completionSearch(srch, dir string) (completions []os.FileInfo) {
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		return
	}

	for _, file := range files {
		if strings.HasPrefix(file.Name(), srch) {
			completions = append(completions, file)
		}
	}
	return
}

var itemSep = " "          //TODO: #63 Add itemSep to config management
var groupSep = []rune{'|'} //TODO: #69 Add groupSep to config management
var pathSep = "/"          //TODO: #70 Add pathSep to config management

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

// lastItem takes the whole buffer from InputView and returns the
// item to be possibly completed and it's path. If the item is the
// first item in a cmd group and without path, bin is true.
// In this case files from PATH can be completed.
func lastItem(line string) (cmd, path string, bin bool) {
	bin = true

	//Trim trainling \n
	line = strings.TrimSuffix(line, "\n")
	path, _ = os.Getwd()
	if strings.TrimSpace(line) == "" {
		return
	}

	groups := splitBuffer(line, groupSep)
	items := strings.Split(groups[len(groups)-1], itemSep)
	item := items[len(items)-1]

	//cmd is absolute path, discard workdir
	if strings.HasPrefix(item, pathSep) {
		bin = false
		path = ""
	}

	//handle paths inside
	parts := strings.Split(item, pathSep)
	for i, p := range parts {
		if p != "" && p != "." && i < len(parts)-1 {
			bin = false
			path += pathSep + p
		}
	}
	//cmd is absolute path, and path still empty
	if strings.HasPrefix(item, pathSep) && path == "" {
		path = pathSep
	}
	cmd = parts[len(parts)-1]

	//buffer ends with blank, thus cmd is empty
	if strings.HasSuffix(line, " ") {
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
