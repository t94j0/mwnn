package client

import (
	"fmt"
	"net"

	"github.com/t94j0/gocui"
)

// Layout for the main function's GUI to set up.
func gocuiLayout(g *gocui.Gui) error {
	maxX, maxY := g.Size()

	// How many lines the input box is.
	//heightOfBox := 4

	if v, err := g.SetView("input_box", 1, (maxY - (maxY / 10)), maxX-1, maxY-1); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}

		// Allow us to type into the editbox
		if err := g.SetCurrentView("input_box"); err != nil {
			return err
		}

		v.Editable = true
		v.Autoscroll = false
		v.Wrap = true
	}

	if v, err := g.SetView("messages_box", (maxX / 10), 0, maxX-1, (maxY - (maxY / 10))); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}

		v.Autoscroll = true
		v.Frame = false

		// Show an initial message
		fmt.Fprintln(v, "Welcome to messenger.")
	}

	if v, err := g.SetView("channel_box", 1, 1, (maxX / 10), (maxY - (maxY / 10) - 1)); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}

		v.Autoscroll = false
		v.Frame = true
		v.Title = "Channels:"

	}

	return nil
}

func createViewKeysWindow(g *gocui.Gui) (*gocui.View, error) {
	maxX, maxY := g.Size()
	view, err := g.SetView("view_keys", 0, 0, maxX-1, maxY-1)
	if err != nil {
		if err != gocui.ErrUnknownView {
			return nil, err
		}
		view.Autoscroll = true
		view.Title = "Known Public Keys"

		// TODO: Write the error handling locs
		if err := g.SetCurrentView("input_box"); err != nil {
			return nil, err
		}

	}

	return view, nil
}

func keybindings(g *gocui.Gui, c net.Conn) error {

	if err := g.SetKeybinding("", gocui.KeyEsc, gocui.ModNone,
		func(g *gocui.Gui, v *gocui.View) error {
			return gocui.ErrQuit
		}); err != nil {
		return err
	}
	if err := g.SetKeybinding("", gocui.KeyCtrlC, gocui.ModNone,
		func(g *gocui.Gui, v *gocui.View) error {
			return gocui.ErrQuit
		}); err != nil {
		return err
	}

	if err := g.SetKeybinding("", gocui.KeyCtrlQ, gocui.ModNone,
		func(g *gocui.Gui, v *gocui.View) error {
			return g.DeleteView("view_keys")
		}); err != nil {
		return err
	}

	return nil
}
