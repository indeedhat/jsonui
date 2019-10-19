# JSONUI

`jsonui` is an interactive JSON explorer in your command line. You can pipe any JSON into `jsonui` and explore it, copy the path for each element.

*Note:* this is a fork of [gulyasm/jsonui](https://github.com/gulyasm/jsonui) and all the hard work along with the credit
lies squarely on the shoulders of gulyasm.\
I have mearly added some quality of life updates

![](img/jsonui.gif)

## Install
`go get -u github.com/indeedhat/jsonui`

## Usage

### Standard output
```bash
cat test_big.json | jsonui
```
### Clipboard
```bash
jsonui -c
```
Clipboard support is handled by [atotto/clipboard](https://github.com/atotto/clipboard)
and is supported on:
- OSX
- Windows 7 (probably work on other Windows)
- Linux, Unix (requires 'xclip' or 'xsel' command to be installed)

### From File
```bash
jsonui -f /path/to/file.json
```

## QOL Changes
 - Load json from clipboard
 - Load json from filepath
 - Copy selected path to clipboard
 - Copy selected value to clipboard
 - Extra key bindings

## Keys

#### `j`, `DownArrow`
Move down a line

#### `k`, `DownUp`
Move up a line

#### `J/PageDown`
Move down 15 lines

#### `K/PageUp`
Move up 15 lines

#### `h/?`
Toggle Help view

#### `e/o`
Toggle node (expend or collapse)

#### `E/O`
Expand all nodes

#### `C`
Collapse all nodes

#### `y`
Copy the path (currently displayed in path panel) to the clipboard\
(only available if clipboard is supported on your system)

#### `Y`
Copy the value (currently displayed in text panel) to the clipboard\
(only available if clipboard is supported on your system)

#### `q/Ctrl+C`
Quit jsonui


## Acknowledgments
Special thanks for [asciimoo](https://github.com/asciimoo) and the [wuzz](https://github.com/asciimoo/wuzz) project for all the help and suggestions.  

