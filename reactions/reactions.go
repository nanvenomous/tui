package reactions

import (
	"log"

	"github.com/jroimartin/gocui"
)

func Quit(g *gocui.Gui, v *gocui.View) error {
	return gocui.ErrQuit
}

const (
	UP = iota
	DOWN
)

func Move(g *gocui.Gui, v *gocui.View, direction int) error {
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
