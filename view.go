package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"math"
	"strings"

	"github.com/atotto/clipboard"
	"github.com/jroimartin/gocui"
)

const (
	treeView    = "TREE"
	textView    = "TEXT"
	pathView    = "PATH"
	helpView    = "HELP"
	messageView = "MESSGE"
)

type position struct {
	prc    float32
	margin int
}

func (p position) getCoordinate(max int) int {
	// value = prc * MAX + abs
	return int(p.prc*float32(max)) - p.margin
}

type viewPosition struct {
	x0, y0, x1, y1 position
}

func logFile(s string) error {
	d1 := []byte(s + "\n")
	return ioutil.WriteFile("log.txt", d1, 0644)
}

func (vp viewPosition) getCoordinates(maxX, maxY int) (int, int, int, int) {
	var x0 = vp.x0.getCoordinate(maxX)
	var y0 = vp.y0.getCoordinate(maxY)
	var x1 = vp.x1.getCoordinate(maxX)
	var y1 = vp.y1.getCoordinate(maxY)
	return x0, y0, x1, y1
}

var helpWindowToggle = false
var messageViewText = ""

var viewPositions = map[string]viewPosition{
	treeView: {
		position{0.0, 0},
		position{0.0, 0},
		position{0.3, 1},
		position{0.9, 1},
	},
	textView: {
		position{0.3, 0},
		position{0.0, 0},
		position{1.0, 1},
		position{0.9, 1},
	},
	pathView: {
		position{0.0, 0},
		position{0.89, 0},
		position{1.0, 1},
		position{1.0, 1},
	},
}

func setupUi() error {
	var err error

	g, err := gocui.NewGui(gocui.OutputNormal)
	if err != nil {
		return err
	}
	defer g.Close()

	g.SetManagerFunc(layout)

	if err := bindGlobals(g); nil != err {
		return err
	}

	if err := bindTreeView(g); nil != err {
		return err
	}

	g.SelFgColor = gocui.ColorBlack
	g.SelBgColor = gocui.ColorGreen

	if err := g.MainLoop(); err != nil && err != gocui.ErrQuit {
		return err
	}

	return nil
}

func bindGlobals(g *gocui.Gui) error {
	if err := g.SetKeybinding("", gocui.KeyCtrlC, gocui.ModNone, quit); err != nil {
		return err
	}

	if err := g.SetKeybinding("", 'q', gocui.ModNone, quit); err != nil {
		return err
	}

	if err := g.SetKeybinding("", 'h', gocui.ModNone, toggleHelp); err != nil {
		return err
	}

	if err := g.SetKeybinding("", '?', gocui.ModNone, toggleHelp); err != nil {
		return err
	}

	return nil
}

func bindTreeView(g *gocui.Gui) error {
	if err := g.SetKeybinding(treeView, 'k', gocui.ModNone, cursorMovement(-1)); err != nil {
		return err
	}

	if err := g.SetKeybinding(treeView, 'j', gocui.ModNone, cursorMovement(1)); err != nil {
		return err
	}

	if err := g.SetKeybinding(treeView, gocui.KeyArrowUp, gocui.ModNone, cursorMovement(-1)); err != nil {
		return err
	}

	if err := g.SetKeybinding(treeView, gocui.KeyArrowDown, gocui.ModNone, cursorMovement(1)); err != nil {
		return err
	}

	if err := g.SetKeybinding(treeView, 'K', gocui.ModNone, cursorMovement(-15)); err != nil {
		return err
	}

	if err := g.SetKeybinding(treeView, 'J', gocui.ModNone, cursorMovement(15)); err != nil {
		return err
	}

	if err := g.SetKeybinding(treeView, gocui.KeyPgup, gocui.ModNone, cursorMovement(-15)); err != nil {
		return err
	}

	if err := g.SetKeybinding(treeView, gocui.KeyPgdn, gocui.ModNone, cursorMovement(15)); err != nil {
		return err
	}

	if err := g.SetKeybinding(treeView, 'g', gocui.ModNone, cursorJump(-1)); err != nil {
		return err
	}

	if err := g.SetKeybinding(treeView, 'G', gocui.ModNone, cursorJump(1)); err != nil {
		return err
	}

	if err := g.SetKeybinding(treeView, 'e', gocui.ModNone, toggleExpand); err != nil {
		return err
	}

	if err := g.SetKeybinding(treeView, 'o', gocui.ModNone, toggleExpand); err != nil {
		return err
	}

	if err := g.SetKeybinding(treeView, 'E', gocui.ModNone, expandAll); err != nil {
		return err
	}

	if err := g.SetKeybinding(treeView, 'O', gocui.ModNone, expandAll); err != nil {
		return err
	}

	if err := g.SetKeybinding(treeView, 'C', gocui.ModNone, collapseAll); err != nil {
		return err
	}

	if err := g.SetKeybinding(treeView, gocui.KeyEnter, gocui.ModNone, clearFloatingViews); err != nil {
		return err
	}

	if !clipboard.Unsupported {
		if err := g.SetKeybinding(treeView, 'y', gocui.ModNone, copyPathToClipboard); nil != err {
			return err
		}

		if err := g.SetKeybinding(treeView, 'Y', gocui.ModNone, copyValueToClipboard); nil != err {
			return err
		}
	}

	return nil
}

func showMessage(message string) {
	if "" == message {
		messageViewText = ""
	} else {
		messageViewText = " " + message
	}
}

func helpMessage() string {

	helpMessage := `
 JSONUI - Help
----------------------------------------------------
 j/ArrowDown     ═   Move a line down
 k/ArrowUp       ═   Move a line up
 J/PageDown      ═   Move 15 line down
 K/PageUp        ═   Move 15 line up
 g               =   Move to top of the tree
 G               =   Move to the bottom of the tree
 e/o             ═   Toggle expend/collapse node
 E/O             ═   Expand all nodes
 C               ═   Collapse all nodes`

	if !clipboard.Unsupported {
		helpMessage += `
 y               =   Copy path to clipboard
 Y               =   Copy value to clipboard`
	}

	helpMessage += `
 Enter           =   Close Message/Help window
 q/ctrl+c        ═   Exit
 h/?             ═   Toggle help message`
	return helpMessage
}

func layout(g *gocui.Gui) error {
	var views = []string{treeView, textView, pathView}
	maxX, maxY := g.Size()
	for _, view := range views {
		x0, y0, x1, y1 := viewPositions[view].getCoordinates(maxX, maxY)
		if v, err := g.SetView(view, x0, y0, x1, y1); err != nil {
			v.SelFgColor = gocui.ColorBlack
			v.SelBgColor = gocui.ColorGreen

			v.Title = " " + view + " "
			if err != gocui.ErrUnknownView {
				return err

			}
			if v.Name() == treeView {
				v.Highlight = true
				drawTree(g, v, tree)
				// v.Autoscroll = true
			}
			if v.Name() == textView {
				drawJSON(g, v)
			}

		}
	}

	if err := renderHelpView(g); nil != err {
		return err
	}

	if err := renderMessageView(g); nil != err {
		return err
	}

	_, err := g.SetCurrentView(treeView)
	if err != nil {
		log.Fatal("failed to set current view: ", err)
	}
	return nil

}

func renderHelpView(g *gocui.Gui) error {
	maxX, maxY := g.Size()

	if helpWindowToggle {
		height := strings.Count(helpMessage(), "\n") + 1
		width := -1

		for _, line := range strings.Split(helpMessage(), "\n") {
			width = int(math.Max(float64(width), float64(len(line)+2)))
		}

		if v, err := g.SetView(helpView, maxX/2-width/2, maxY/2-height/2, maxX/2+width/2, maxY/2+height/2); err != nil {
			if err != gocui.ErrUnknownView {
				return err
			}

			fmt.Fprintln(v, helpMessage())
		}
	} else {
		g.DeleteView(helpView)
	}

	return nil
}

func renderMessageView(g *gocui.Gui) error {
	maxX, maxY := g.Size()

	if "" != messageViewText {
		height := strings.Count(messageViewText, "\n") + 2
		width := -1

		for _, line := range strings.Split(messageViewText, "\n") {
			width = int(math.Max(
				float64(width),
				float64(len(line)+2),
			))
		}

		if v, err := g.SetView(messageView, maxX/2-width/2, maxY/2-height/2, maxX/2+width/2, maxY/2+height/2); err != nil {
			if err != gocui.ErrUnknownView {
				return err
			}

			fmt.Fprintln(v, messageViewText)
		}
	} else {
		g.DeleteView(messageView)
	}

	return nil
}
