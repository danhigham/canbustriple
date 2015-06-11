package main

import	"github.com/danhigham/gocui"

func (c *CanbusClient) keybindings(g *gocui.Gui) error {

  if err := g.SetKeybinding("", gocui.KeyCtrlO, gocui.ModNone, c.switchToOptions); err != nil {
		return err
	}

	if err := g.SetKeybinding("", gocui.KeyCtrlC, gocui.ModNone, quit); err != nil {
		return err
	}

	if err := g.SetKeybinding("", gocui.KeyCtrlI, gocui.ModNone, c.requestTripleInfo); err != nil {
		return err
	}

	if err := g.SetKeybinding("", gocui.KeyCtrlV, gocui.ModNone, c.toggleCompactView); err != nil {
		return err
	}

  if err := g.SetKeybinding("side-options", gocui.KeyArrowDown, gocui.ModNone, cursorDown); err != nil {
		return err
	}

	if err := g.SetKeybinding("side-options", gocui.KeyArrowUp, gocui.ModNone, cursorUp); err != nil {
		return err
	}

  if err := g.SetKeybinding("side-options", gocui.KeySpace, gocui.ModNone, c.setOptions); err != nil {
		return err
	}

	if err := g.SetKeybinding("msg", gocui.KeyEnter, gocui.ModNone, delMsg); err != nil {
		return err
	}

	if err := g.SetKeybinding("", gocui.KeySpace, gocui.ModNone, c.togglePause); err != nil {
		return err
	}

	if err := g.SetKeybinding("main", gocui.KeyArrowUp, gocui.ModNone, func (g *gocui.Gui, v *gocui.View) error {
		if c.SelectedLine > 0 { c.SelectedLine -- }
		return nil
	}); err != nil { return err }

	if err := g.SetKeybinding("main", gocui.KeyArrowDown, gocui.ModNone, func(g *gocui.Gui, v *gocui.View) error {
		if c.SelectedLine < len(c.Packets) - 1 { c.SelectedLine ++ }
		return nil
	}); err != nil { return err }

	if err := g.SetKeybinding("main", gocui.KeyEnter, gocui.ModNone, func(g *gocui.Gui, v *gocui.View) error {

		// if c.SelectedLine < len(c.Packets) - 1 { c.SelectedLine ++ }
		keys := c.packetKeys()

		p := c.Packets[keys[c.SelectedLine]]
		c.showActionDialog(p)

		return nil
	}); err != nil { return err }

	return nil
}

func cursorDown(g *gocui.Gui, v *gocui.View) error {
	if v != nil {
		cx, cy := v.Cursor()
    if (cy > 4) { return nil }

		if err := v.SetCursor(cx, cy+1); err != nil {
			ox, oy := v.Origin()
			if err := v.SetOrigin(ox, oy+1); err != nil {
				return err
			}
		}
	}
	return nil
}

func cursorUp(g *gocui.Gui, v *gocui.View) error {
	if v != nil {
		ox, oy := v.Origin()
		cx, cy := v.Cursor()

    if (cy < 4) { return nil }

		if err := v.SetCursor(cx, cy-1); err != nil && oy > 0 {
			if err := v.SetOrigin(ox, oy-1); err != nil {
				return err
			}
		}
	}
	return nil
}
