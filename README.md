# GOSH

The following line operators are supported:

| Operator | Action                                                                           |
| -------- | -------------------------------------------------------------------------------- |
| &&       | AND Operator (executes the following command if the previous one was successful) |
| \|\|     | OR Operator (executes the following command if the previous one failed)          |
| ;        | Semicolon Operator (always executes the following command)                       |
| &        | Execute command in background                                                    |
| $?       | Get last process' exit code                                                      |

--

The following line editing commands are supported:

| Keystroke             | Action                                              |
| --------------------- | --------------------------------------------------- |
| Ctrl-A, Home          | Move cursor to beginning of line                    |
| Ctrl-E, End           | Move cursor to end of line                          |
| Ctrl-B, Left          | Move cursor one character left                      |
| Ctrl-F, Right         | Move cursor one character right                     |
| Ctrl-Left, Alt-B      | Move cursor to previous word                        |
| Ctrl-Right, Alt-F     | Move cursor to following word                       |
| Ctrl-T                | Transpose previous character with current character |
| Delete                | Delete current character                            |
| Alt-D                 | Delete next word                                    |
| Ctrl-W, Alt-BackSpace | Delete current word                                 |
| Ctrl-U                | Delete from start of line to cursor                 |
| Ctrl-K                | Delete from cursor to end of line                   |
| Ctrl-H, BackSpace     | Delete character before cursor                      |
| Ctrl-L                | Clear screen                                        |
| Ctrl-D                | End of File - if line is empty quits application    |
| Ctrl-C                | Reset input (create new empty prompt)               |
| Ctrl-P, Up            | Previous command from history                       |
| Ctrl-N, Down          | Next command from history                           |
| Tab                   | Word completion                                     |
| Winch Signal          | Window change                                       |

--

Additional commands also supported:

| Command                | Action                           |
| ---------------------- | -------------------------------- |
| cd                     | Go to system root directory      |
| cd \<path\>            | Go to required directory         |
| cd -                   | Go to last registered directory  |
| echo $\<env variable\> | Print environment variable value |
