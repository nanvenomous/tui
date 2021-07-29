package tui

import (
	"fmt"
	"log"
	"math"

	"github.com/jroimartin/gocui"
)

var (
	choices     []string
	finalChoice string
)

const (
	UP = iota
	DOWN
)

func LaunchSelection(inputChoices []string) string {
	finalChoice = ""
	if len(inputChoices) < 1 {
		log.Panicln("must provide at least one choice")
	}
	choices = inputChoices

	g, err := gocui.NewGui(gocui.OutputNormal)
	if err != nil {
		log.Panicln(err)
	}
	defer g.Close()

	g.Cursor = true
	g.Mouse = true

	g.SetManagerFunc(layout)

	if err := keybindings(g); err != nil {
		log.Panicln(err)
	}

	if err := g.MainLoop(); err != nil && err != gocui.ErrQuit {
		log.Panicln(err)
	}
	return finalChoice
}

func getWidth() int {
	max := 0.0
	for _, ch := range choices {
		max = math.Max(max, float64(len(ch)))
	}
	return int(max)
}

func layout(g *gocui.Gui) error {

	yMessagePos := 0
	message := "Select an interface from the list below"
	if v, err := g.SetView("instructions", 0, yMessagePos, len(message)+1, 2); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		fmt.Fprintln(v, message)
	}

	yStart := yMessagePos + 3
	yEnd := yStart + 1 + len(choices)
	xStart := 3
	xEnd := xStart + 1 + getWidth()
	if v, err := g.SetView("but1", xStart, yStart, xEnd, yEnd); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		v.Highlight = true
		v.SelBgColor = gocui.ColorGreen
		v.SelFgColor = gocui.ColorBlack
		for _, choice := range choices {
			fmt.Fprintln(v, choice)
		}
		g.SetCurrentView(v.Name())
	}
	return nil
}

func getKey(symbol string) rune {
	return []rune(symbol)[0]
}

func keybindings(g *gocui.Gui) error {
	if err := g.SetKeybinding("", gocui.KeyCtrlC, gocui.ModNone, quit); err != nil {
		return err
	}

	if err := g.SetKeybinding("", getKey("q"), gocui.ModNone, quit); err != nil {
		return err
	}

	if err := g.SetKeybinding("but1", getKey("j"), gocui.ModNone, func(g *gocui.Gui, v *gocui.View) error {
		err := move(g, v, DOWN)
		if err != nil {
			return err
		}
		return nil
	}); err != nil {
		return err
	}

	if err := g.SetKeybinding("but1", getKey("k"), gocui.ModNone, func(g *gocui.Gui, v *gocui.View) error {
		err := move(g, v, UP)
		if err != nil {
			return err
		}
		return nil
	}); err != nil {
		return err
	}

	if err := g.SetKeybinding("but1", gocui.MouseLeft, gocui.ModNone, showMsg); err != nil {
		return err
	}

	if err := g.SetKeybinding("but1", gocui.KeyEnter, gocui.ModNone, showMsg); err != nil {
		return err
	}

	if err := g.SetKeybinding("msg", gocui.MouseLeft, gocui.ModNone, delMsg); err != nil {
		return err
	}
	return nil
}

func quit(g *gocui.Gui, v *gocui.View) error {
	return gocui.ErrQuit
}

func move(g *gocui.Gui, v *gocui.View, direction int) error {
	var inc int

	if _, err := g.SetCurrentView(v.Name()); err != nil {
		return err
	}

	switch direction {
	case UP:
		inc = -1
	case DOWN:
		inc = 1
	default:
		log.Panic("incorrect direction given")
	}

	_, cy := v.Cursor()
	if err := v.SetCursor(0, cy+inc); err == nil {
		return err
	}
	return nil
}

func showMsg(g *gocui.Gui, v *gocui.View) error {
	var l string
	var err error

	if _, err := g.SetCurrentView(v.Name()); err != nil {
		return err
	}

	_, cy := v.Cursor()
	if l, err = v.Line(cy); err != nil {
		l = ""
	}
	finalChoice = l

	return quit(g, v)
}

func delMsg(g *gocui.Gui, v *gocui.View) error {
	if err := g.DeleteView("msg"); err != nil {
		return err
	}
	return nil
}
