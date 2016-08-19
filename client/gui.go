package main

import (
	"fmt"
	"net"

	"github.com/t94j0/gocui"
)

// Layout for the main function's GUI to set up.
func gocuiLayout(g *gocui.Gui) error {
	maxX, maxY := g.Size()

	// How many lines the input box is.
	heightOfBox := 4

	if v, err := g.SetView("input_box", 1, maxY-heightOfBox+1, maxX-1, maxY-1); err != nil {
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

	if v, err := g.SetView("messages_box", 0, 0, maxX-1, maxY-heightOfBox-1); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}

		v.Autoscroll = true
		v.Frame = false

		// Show an initial message
		fmt.Fprintln(v, "Welcome to messenger.")
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
	// If we are on any view and the enter button is pressed, submit whats in the editbox buffer
	// to the server.
	// messageHandler is in main.go
	if err := g.SetKeybinding("", gocui.KeyEnter, gocui.ModNone, messageHandler); err != nil {
		return err
	}

	return nil
}
