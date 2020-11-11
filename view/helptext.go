package view

const helptext = `
Enter commands:
Write to command line and hit <CR>.

Help: 
Scroll up and down with arrow keys.

History:
Start scrolling backwards through command history with <UpArrow>. <DownArrow> scrolls forward. Hit <F2> to open history search. Refining your search pattern refreshes the results. Use <UpArrow> and <DownArrows> as usual. Select a line with <Enter>.

Completion:
Press <TAB> to complete. If more than one possible completions are found, a new view is opened. Navigate with Arrows aus usual, <ENTER> or <TAB> selects. Hit <EsC> to close the view without actually completing. Fi you type only "cd" and hit <TAB>, a special view for quick directory changes opens.

Change directories:
Hit <F3> to quickly change deep into the directory tree. <ArrowLeft> moves up one directory, <ArrowRight> dives into the directory under your cursor. As always <Enter> selects. Besides moving with <ArrowUp> and <ArrowDown> you can jump by typing the first letter of any directory name.
`