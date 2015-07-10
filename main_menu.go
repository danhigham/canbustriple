package main

import (
	"fmt"

	"github.com/danhigham/gocui"
)

func (c *CanbusClient) createMainMenu(g *gocui.Gui) error {

	maxX, maxY := g.Size()

	if v, err := g.SetView("cmdline", -1, maxY-2, maxX, maxY); err != nil {

		if err != gocui.ErrorUnkView {
			return err
		}

		v.Underlined = []gocui.Pos{{0, 0}, {10, 0}, {18, 0}, {43, 0}, {53, 0}, {70, 0}, {77, 0}}

		fmt.Fprintf(v, "Options\t\t\t")
		fmt.Fprintf(v, "Pause\t\t\t")
		fmt.Fprintf(v, "View: Compact\t\t\t")
		fmt.Fprintf(v, "Send CAN Message\t\t\t")
		fmt.Fprintf(v, "Filter\t\t\t")
		fmt.Fprintf(v, "Get Sys Info\t\t\t")
		fmt.Fprintf(v, "Quit\t\t\t")

		if err := g.SetKeybinding("", gocui.KeyCtrlO, gocui.ModNone, c.switchToOptions); err != nil {
			return err
		}

		if err := g.SetKeybinding("", gocui.KeyCtrlI, gocui.ModNone, func(g *gocui.Gui, v *gocui.View) error {
			c.TripleClient.RequestInfo()
			return nil
		}); err != nil {
			return err
		}

		if err := g.SetKeybinding("", gocui.KeyCtrlV, gocui.ModNone, func(g *gocui.Gui, v *gocui.View) error {
			c.ShowCompact = !c.ShowCompact
			c.mainView.Highlight = c.ShowCompact
			c.SelectedLine = 0
			g.SetCurrentView("main")

			c.mainView.Clear()
			return nil
		}); err != nil {
			return err
		}

		if err := g.SetKeybinding("", gocui.KeyCtrlQ, gocui.ModNone, func(g *gocui.Gui, v *gocui.View) error {
			return gocui.Quit
		}); err != nil {
			return err
		}

	}

	return nil
}
