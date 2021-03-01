package view

const helptext = `Enter commands:
Write to command line and hit <CR>. Use <LeftArrow> and <RightArrow> to move your cursor by one char. Jump one word with <CTRL-L> and <CTRL-R>.

Help: 
Scroll up and down with arrow keys.

History:
Start scrolling backwards through command history with <UpArrow>. <DownArrow> scrolls forward. Hit <F2> to open history search. Refining your search pattern refreshes the results. Use <UpArrow> and <DownArrows> as usual. Select a line with <Enter>.

Completion:
Press <TAB> to complete. If more than one possible completions are found, a new view is opened. Navigate with Arrows aus usual, <ENTER> or <TAB> selects. Hit <EsC> to close the view without actually completing. Fi you type only "cd" and hit <TAB>, a special view for quick directory changes opens.

Change directories:
Hit <F3> to quickly change deep into the directory tree. <ArrowLeft> moves up one directory, <ArrowRight> dives into the directory under your cursor. As always <Enter> selects. Besides moving with <ArrowUp> and <ArrowDown> you can jump by typing the first letter of any directory name.

Background jobs:
Background jobs can be startet by prefixing normal commands with "bg". Any output to stdout or stderr is stored. Jobs are viewed and managed in a dedicated view. This view opens with <F4>. Use <Tab> to toggle the focus between the list of job in the left view and the output of one job on the right side. Enter "bg ls" to view a list of jobs on the console. By typing "bg" followed by an integer, the output of this job is printed.

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

Renaming files in bulk:
With pressing <F5> a view is opened for renaming files or diretories in bulk. The left part of this view shows a simple editor to write down rules for renaming. The right part shows files and folders of the current diretory and previews the effect of currently active rules. once the preview shows the wanted result, pressing <F6> actually performs the renaming. Rules are as follows:
--------------------------------------------------------
Insert
ins string pre 				insert "string" as prefix
ins string suf 				insert "string" as suffix
ins string pos 2			insert "string" at position 2
ins string aft substrings	insert "string" aft first occurence of substring
In case string parses as integer, ins increments string for every file

Delete
del string 					delete first occurence of string
del string pre 				delete string if its the prefix
del string suf 				delete string it its the suffix
del string any 				delete any matching strings
del 1						delete from position 1 to end
del 1 3						delete from position 1 to 3
del -3						delete last 3 characters

Replace
rep old new 				replace the first occurence of "old" with "new"
rep old new pre 			replace "old" with "new" when old is a prefix
rep old new suf 			replace "old" with "new" when old is a suffix
rep old new any 			replace any occurence of "old" with "new"
In case "new" parses as integer, ins increments string for every file

Change capitalization
case upp 					first Char to uppercase
case upp wrd 				first char of every word uppercase
case upp any [string]		any matching substring to uppercase, on empty string, any char is uppercase
Use low instead of upp to lowercase characters

Date
dat 						use filedate as name
dat pre 					use filedate as prefix
dat suf 					use filedate as suffix

Rules for files only or directoies only
mod all 					following rules apply to files and diretories
mod fil 					rules apply to files only
mod dir 					rules apply to directories only

Include and exclude
inc regex 					following rules apply only to files or folders that match regex 
exc regex 					following rules don't apply to files or folders that match regex
--------------------------------------------------------
`