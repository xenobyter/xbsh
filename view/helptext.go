package view

const helptext = `Enter commands:
Write to command line and hit <CR>.

Help: 
Scroll up and down with arrow keys.

History:
Start scrolling backwards through command history with <UpArrow>. <DownArrow> scrolls forward. Hit <F2> to open history search. Refining your search pattern refreshes the results. Use <UpArrow> and <DownArrows> as usual. Select a line with <Enter>.

Completion:
Press <TAB> to complete. If more than one possible completions are found, a new view is opened. Navigate with Arrows aus usual, <ENTER> or <TAB> selects. Hit <EsC> to close the view without actually completing. Fi you type only "cd" and hit <TAB>, a special view for quick directory changes opens.

Change directories:
Hit <F3> to quickly change deep into the directory tree. <ArrowLeft> moves up one directory, <ArrowRight> dives into the directory under your cursor. As always <Enter> selects. Besides moving with <ArrowUp> and <ArrowDown> you can jump by typing the first letter of any directory name.

Background jobs:
Background jobs can be startet by prefixing normal commands with "bg". Any output to stdout or stderr is stored. Jobs are viewed and managed in a dedicated view. This view opens with <F4>. Use <Tab> to toggle the focus between the list of job in the left view and the output of one job on the right side.

Redirections:
One can redirect stdout to any file by adding ">" followed by a filename. If file exists, its conetents will be overridden. To append stdout at the end of a file without overriding use ">>". To redirect stderr use ">!" followed by a filename, to append stderr to a file use ">>!".

Pipes:
Use "|" to pipe one commands stdout to next commands stdin.

Configuration:
To customize some aspects xbsh reads a simple ini-file from ~/.xbsh/xbshrc. Example:
--------------------------------------------------------
#xbsh config
[separators]
pathsep=/
pipesep=|

[history]
maxentries=1000
delExitCmd=exit
--------------------------------------------------------
`