package main

import (
	"fmt"
	"log"
	"math"
	"strings"

	"github.com/atotto/clipboard"
	"github.com/jroimartin/gocui"
)

const VERSION = "1.2.0"

var tree treeNode

func getPath(g *gocui.Gui, v *gocui.View) string {
	p := findTreePosition(g)
	for i, s := range p {
		transformed := s
		if !strings.HasPrefix(s, "[") && !strings.HasSuffix(s, "]") {
			transformed = fmt.Sprintf("[%q]", s)
		}
		p[i] = transformed
	}
	return strings.Join(p, "")
}

func drawPath(g *gocui.Gui, v *gocui.View) error {
	pv, err := g.View(pathView)
	if err != nil {
		log.Fatal("failed to get pathView", err)
	}
	p := getPath(g, v)
	pv.Clear()
	fmt.Fprintf(pv, p)
	return nil
}

func drawJSON(g *gocui.Gui, v *gocui.View) error {
	dv, err := g.View(textView)
	if err != nil {
		log.Fatal("failed to get textView", err)
	}
	p := findTreePosition(g)
	treeTodraw := tree.find(p)
	if treeTodraw != nil {
		dv.Clear()
		fmt.Fprintf(dv, treeTodraw.String(2, 0))
	}
	return nil
}

func lineBelow(v *gocui.View, d int) bool {
	_, y := v.Cursor()
	line, err := v.Line(y + d)
	return err == nil && line != ""
}

func countIndex(s string) int {
	count := 0
	for _, c := range s {
		if c == ' ' {
			count++
		}
	}
	return count
}

func getLine(s string, y int) string {
	lines := strings.Split(s, "\n")
	return lines[y]
}

var cleanPatterns = []string{
	treeSignUpEnding,
	treeSignDash,
	treeSignUpMiddle,
	treeSignVertical,
	" (+)",
}

func findTreePosition(g *gocui.Gui) treePosition {
	v, err := g.View(treeView)
	if err != nil {
		log.Fatal("failed to get treeview", err)
	}
	path := treePosition{}
	ci := -1
	_, yOffset := v.Origin()
	_, yCurrent := v.Cursor()
	y := yOffset + yCurrent
	s := v.Buffer()
	for cy := y; cy >= 0; cy-- {
		line := getLine(s, cy)
		for _, pattern := range cleanPatterns {
			line = strings.Replace(line, pattern, "", -1)
		}

		if count := countIndex(line); count < ci || ci == -1 {
			path = append(path, strings.TrimSpace(line))
			ci = count
		}
	}
	for i := len(path)/2 - 1; i >= 0; i-- {
		opp := len(path) - 1 - i
		path[i], path[opp] = path[opp], path[i]
	}

	return path[1:]
}

// This is a workaround for not having a Buffer
// function in gocui
func bufferLen(v *gocui.View) int {
	s := v.Buffer()
	return len(strings.Split(s, "\n")) - 1
}

func drawTree(g *gocui.Gui, v *gocui.View, tree treeNode) error {
	tv, err := g.View(treeView)
	if err != nil {
		log.Fatal("failed to get treeView", err)
	}
	tv.Clear()
	tree.draw(tv, 2, 0)
	maxY := bufferLen(tv)
	cx, cy := tv.Cursor()
	lastLine := maxY - 2
	if cy > lastLine {
		tv.SetCursor(cx, lastLine)
		tv.SetOrigin(0, 0)
	}

	return nil
}

func expandAll(g *gocui.Gui, v *gocui.View) error {
	tree.expandAll()
	return drawTree(g, v, tree)
}

func collapseAll(g *gocui.Gui, v *gocui.View) error {
	tree.collapseAll()
	return drawTree(g, v, tree)
}

func toggleExpand(g *gocui.Gui, v *gocui.View) error {
	p := findTreePosition(g)
	subTree := tree.find(p)
	subTree.toggleExpanded()
	return drawTree(g, v, tree)
}

func cursorMovement(d int) func(g *gocui.Gui, v *gocui.View) error {
	return func(g *gocui.Gui, v *gocui.View) error {
		dir := 1
		if d < 0 {
			dir = -1
		}
		distance := int(math.Abs(float64(d)))
		for ; distance > 0; distance-- {
			if lineBelow(v, dir) {
				v.MoveCursor(0, dir, false)
			}
		}

		drawJSON(g, v)
		drawPath(g, v)
		return nil
	}
}

func cursorJump(direction int) func(g *gocui.Gui, v *gocui.View) error {
	return func(g *gocui.Gui, v *gocui.View) error {
		xMax := bufferLen(v)

		if direction < 0 {
			return cursorMovement(-xMax)(g, v)
		} else {
			return cursorMovement(xMax)(g, v)
		}

		return nil
	}
}

func copyPathToClipboard(g *gocui.Gui, v *gocui.View) error {
	path := getPath(g, v)
	err := clipboard.WriteAll(path)
	if nil != err {
		return err
	}

	showMessage("Path coppied to clipboard")
	return nil
}

func copyValueToClipboard(g *gocui.Gui, v *gocui.View) error {
	p := findTreePosition(g)
	treeTodraw := tree.find(p)

	if nil == treeTodraw {
		return nil
	}

	err := clipboard.WriteAll(treeTodraw.String(2, 0))
	if nil != err {
		return err
	}

	showMessage("Current node coppied to clipboard")
	return nil
}

func toggleHelp(g *gocui.Gui, v *gocui.View) error {
	helpWindowToggle = !helpWindowToggle
	return nil
}

func clearFloatingViews(g *gocui.Gui, v *gocui.View) error {
	showMessage("")
	helpWindowToggle = false
	return nil
}

func quit(g *gocui.Gui, v *gocui.View) error {
	return gocui.ErrQuit
}
