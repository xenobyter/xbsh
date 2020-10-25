package term

//ANSI Sequences
const (
	//Colors
	AnsiPrompt  = "\033[0;32m"
	AnsiNormal  = "\033[0m"
	AnsiErr     = "\033[0;31m"
	AnsiBGGreen = "\033[42m"
	AnsiBlack   = "\033[30m"

	// Cursor Movement
	AnsiSave    = "\0337"
	AnsiRestore = "\0338"
	AnsiLeft    = "\033[1D"
	AnsiRight   = "\033[1C"
	AnsiUp      = "\033[1A"
	AnsiDown    = "\033[1B"

	//Delete
	AnsiClrLine = "\033[0G\033[2K"
	AnsiCls     = "\033[2J\033[0;0r"

	AnsiUnderLine = "\u001b[4m"
)
