package tui

import (
	"fmt"
	"io/ioutil"
	"log"

	"github.com/jroimartin/gocui"
	"github.com/mrgarelli/tui/reactions"
)

var (
	viewArr   = []string{"v1", "v2", "v3", "v4"}
	active    = 0
	boxHeight = 2
	fileBytes = []byte{}
)

func setCurrentViewOnTop(g *gocui.Gui, name string) (*gocui.View, error) {
	if _, err := g.SetCurrentView(name); err != nil {
		return nil, err
	}
	return g.SetViewOnTop(name)
}

func nextView(g *gocui.Gui, v *gocui.View) error {
	nextIndex := (active + 1) % len(viewArr)
	name := viewArr[nextIndex]

	out, err := g.View("v2")
	if err != nil {
		return err
	}
	fmt.Fprintln(out, "Going from view "+v.Name()+" to "+name)

	if _, err := setCurrentViewOnTop(g, name); err != nil {
		return err
	}

	if nextIndex == 0 || nextIndex == 3 {
		g.Cursor = true
	} else {
		g.Cursor = false
	}

	active = nextIndex
	return nil
}

func newLine(g *gocui.Gui, v *gocui.View) error {

	go g.Update(func(g *gocui.Gui) error {
		fileString := v.ViewBuffer()
		fileBytes = []byte(fileString)
		maxX, _ := g.Size()
		boxHeight = countRune(fileString, []rune("\n")[0]) + 2
		if _, err := g.SetView(v.Name(), 0, 0, maxX/2-1, boxHeight); err != nil {
			return err
		}

		v.EditNewLine()
		return nil
	})
	return nil
}

func keybindingsRTB(g *gocui.Gui) error {
	if err := g.SetKeybinding("", gocui.KeyCtrlC, gocui.ModNone, reactions.Quit); err != nil {
		log.Panicln(err)
	}
	if err := g.SetKeybinding("", gocui.KeyEnter, gocui.ModNone, nextView); err != nil {
		log.Panicln(err)
	}

	if err := g.SetKeybinding("v1", gocui.KeyEnter, gocui.ModAlt, newLine); err != nil {
		log.Panicln(err)
	}

	return nil
}

func ResizableTextBox() error {
	g, err := gocui.NewGui(gocui.OutputNormal)
	if err != nil {
		log.Panicln(err)
	}
	defer g.Close()

	g.Highlight = true
	g.Cursor = true
	g.SelFgColor = gocui.ColorGreen

	fileBytes, err = ioutil.ReadFile("./_example/rsrc/editable_file.txt")
	if err != nil {
		return err
	}

	g.SetManagerFunc(layoutRTB)

	if err := keybindingsRTB(g); err != nil {
		log.Panicln(err)
	}

	if err := g.MainLoop(); err != nil && err != gocui.ErrQuit {
		log.Panicln(err)
	}
	return nil
}

func countRune(s string, r rune) int {
	count := 0
	for _, c := range s {
		if c == r {
			count++
		}
	}
	return count
}

func layoutRTB(g *gocui.Gui) error {
	maxX, maxY := g.Size()
	boxHeight = countRune(string(fileBytes), []rune("\n")[0]) + 2
	if v, err := g.SetView("v1", 0, 0, maxX/2-1, boxHeight); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		v.Title = "v1 (editable)"
		v.Editable = true
		v.Write(fileBytes)

		if _, err = setCurrentViewOnTop(g, "v1"); err != nil {
			return err
		}
	}

	if v, err := g.SetView("v2", maxX/2-1, 0, maxX-1, maxY/2-1); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		v.Title = "v2"
		v.Wrap = true
		v.Autoscroll = true
	}
	if v, err := g.SetView("v3", 0, maxY/2-1, maxX/2-1, maxY-1); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		v.Title = "v3"
		v.Wrap = true
		v.Autoscroll = true
		fmt.Fprint(v, "Press TAB to change current view")
	}
	if v, err := g.SetView("v4", maxX/2, maxY/2, maxX-1, maxY-1); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		v.Title = "v4 (editable)"
		v.Editable = true
	}
	return nil
}
