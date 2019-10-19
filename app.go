package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/atotto/clipboard"
	"github.com/indeedhat/gli"
)

const (
	ERROR_NONE = iota
	ERROR_JSON
	ERROR_UI
	ERROR_CLIPBOARD
	ERROR_FILEPATH
)

type JsonUi struct {
	Help          bool   `gli:"help,h" desscription:"Show this help page"`
	FilePath      string `gli:"file,f" description:"Load json directly from a file"`
	FromClipboard bool   `gli:"clipboard,clip,c" description:"Load json from the clipboard"`
}

func (j *JsonUi) Run() int {
	var err error

	if "" != j.FilePath {
		path, err := filepath.Abs(j.FilePath)
		if nil != err {
			fmt.Println(err)
			return ERROR_FILEPATH
		}

		json, err := ioutil.ReadFile(path)
		if nil != err {
			fmt.Println(err)
			return ERROR_FILEPATH
		}

		tree, err = fromBytes(json)
	} else if j.FromClipboard {
		if clipboard.Unsupported {
			fmt.Println("Clipboard is unsupported on your system")
			return ERROR_CLIPBOARD
		}

		json, err := clipboard.ReadAll()
		if nil != err {
			fmt.Println(err)
		}

		tree, err = fromBytes([]byte(json))
	} else {
		tree, err = fromReader(os.Stdin)
	}

	if nil != err {
		fmt.Println(err)
		return ERROR_JSON
	}

	if err = setupUi(); nil != err {
		fmt.Println(err)
		return ERROR_UI
	}

	return ERROR_NONE
}

func (j *JsonUi) NeedHelp() bool {
	return j.Help
}

func main() {
	app := gli.NewApplication(&JsonUi{}, "Terminal user interface for exploring a json structure")
	app.Run()
}
