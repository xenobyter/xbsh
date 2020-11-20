package exec

import (
	"path/filepath"
)

func expandArg(args []string) (expArgs []string) {
	for _, arg := range args {
		match, _ := filepath.Glob(arg)
		if match == nil {
			match = append(match, arg)
		}
		expArgs = append(expArgs, match...)
	}
	return
}
