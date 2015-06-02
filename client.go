// Copyright 2014 The gocui Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"os"
	"log"
  "fmt"
	"github.com/jroimartin/gocui"
)

type CanbusClient struct {
  filterView  	*gocui.View
  optionView  	*gocui.View
  mainView    	*gocui.View
  options     	CanbusClientOptions
	TripleClient	*TripleClient
	g 						*gocui.Gui
}

type CanbusClientOptions struct {
  bus1Enabled bool
  bus2Enabled bool
  bus3Enabled bool
}

func (c *CanbusClient) layout(g *gocui.Gui) error {
	maxX, maxY := g.Size()

  var err error

	if c.filterView, err = g.SetView("side-filters", -1, -1, 20, 10); err != nil {
		if err != gocui.ErrorUnkView {
		  return err
    }

    fmt.Fprint(c.filterView, "Active Filters\n-------------------\n")

  }

  if c.optionView, err = g.SetView("side-options", -1, 10, 20, maxY-2); err != nil {
		if err != gocui.ErrorUnkView {
		  return err
    }

    c.optionView.Highlight = true
    c.writeOptionsPane()
  }

	if c.mainView, err = g.SetView("main", 20, -1, maxX, maxY-2); err != nil {
		if err != gocui.ErrorUnkView {
			return err
		}

		c.mainView.Autoscroll = true

  }

	if v, err := g.SetView("cmdline", -1, maxY-2, maxX, maxY); err != nil {
    if err != gocui.ErrorUnkView {
		  return err
    }

		fmt.Fprint(v, "O: Set Options\t\t\t")
    fmt.Fprint(v, "P: Pause\t\t\t")
    fmt.Fprint(v, "V: Compact\t\t\t")
    fmt.Fprint(v, "M: Send CAN Message\t\t\t")
    fmt.Fprint(v, "F: Add Filter\t\t\t")
    fmt.Fprint(v, "I: Get Sys Info\t\t\t")
    fmt.Fprint(v, "C: Quit\t\t\t")
  }

  return nil
}

func quit(g *gocui.Gui, v *gocui.View) error {
	return gocui.Quit
}

func (c *CanbusClient) writeLoggingOptions() {
	if c.options.bus1Enabled {
		c.TripleClient.SetBus(0x01,0x01)
	} else {
		c.TripleClient.SetBus(0x01,0x00)
	}
	if c.options.bus2Enabled {
		c.TripleClient.SetBus(0x02,0x01)
	} else {
		c.TripleClient.SetBus(0x02,0x00)
	}
	if c.options.bus3Enabled {
		c.TripleClient.SetBus(0x03,0x01)
	} else {
		c.TripleClient.SetBus(0x03,0x00)
	}
}

func (c *CanbusClient) writeOptionsPane() error {
  v := c.optionView
  v.Clear()

  fmt.Fprint(v, "Options\n-------------------\n")

  fmt.Fprintf(v, "\n%s%+12s", "Bus 1", "")
  if c.options.bus1Enabled { fmt.Fprint(v, "\u2714") } else { fmt.Fprint(v, " ") }
  fmt.Fprintf(v, "\n%s%+12s", "Bus 2", "")
  if c.options.bus2Enabled { fmt.Fprint(v, "\u2714") } else { fmt.Fprint(v, " ") }
  fmt.Fprintf(v, "\n%s%+12s", "Bus 3", "")
  if c.options.bus3Enabled { fmt.Fprint(v, "\u2714") } else { fmt.Fprint(v, " ") }

  _, cy := v.Cursor()
  if (cy < 3) { v.SetCursor(0,2) }

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

func (c *CanbusClient) setOptions(g *gocui.Gui, v *gocui.View) error {
  if v != nil {
    _, cy := v.Cursor()

    if cy == 3 { c.options.bus1Enabled = !c.options.bus1Enabled }
    if cy == 4 { c.options.bus2Enabled = !c.options.bus2Enabled }
    if cy == 5 { c.options.bus3Enabled = !c.options.bus3Enabled }

    c.writeOptionsPane()
		c.writeLoggingOptions()
  }

  return nil
}

func (c *CanbusClient) switchToOptions(g *gocui.Gui, v *gocui.View) error {
  g.SetCurrentView("side-options")
  c.optionView.SetCursor(0,3)

  return nil
}

func (c *CanbusClient) requestTripleInfo(g *gocui.Gui, v *gocui.View) error {
	c.TripleClient.RequestInfo()
	return nil
}

func delMsg(g *gocui.Gui, v *gocui.View) error {

	if err := g.DeleteView("msg"); err != nil {
		return err
	}

	if err := g.SetCurrentView("main"); err != nil {
		return err
	}

	return nil
}

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

	return nil
}

func (c *CanbusClient) initCanChannel(ch chan CANPacket) {
	for {
    canPacket := <- ch
		log.Printf("%+v", canPacket)
		c.mainView.SetCursor(20,20)
		s := fmt.Sprintf("%+v\n", canPacket)
		c.mainView.Write([]byte(s))
		c.g.Flush()
  }
}


func (c *CanbusClient) initInfoChannel(ch chan TripleInfo) {

	for {
		info := <- ch

		txt := fmt.Sprintf("\n\n  Event:    %s\n  Name:     %s\n  Version:  %s\n  Memory:   %s\n",
			info.Event, info.Name, info.Version, info.Memory)

		maxX, maxY := c.g.Size()

		if v, err := c.g.SetView("msg", maxX/2-15, maxY/2-5, maxX/2+15, maxY/2+5); err != nil {

			if err != gocui.ErrorUnkView {
				panic(err)
			}

			v.SetOrigin(-10,-10)
			v.SetCursor(20,20)

			fmt.Fprint(v, txt)

			if err := c.g.SetCurrentView("msg"); err != nil {
				panic(err)
			}
		}
	}
}

func main() {

	var err error

	f, err := os.OpenFile("./canbustriple.log", os.O_RDWR | os.O_CREATE | os.O_APPEND, 0666)
	if err != nil {
	    panic(fmt.Sprintf("error opening file: %v", err))
	}
	defer f.Close()

	log.SetOutput(f)

	var port = os.Args[1]

	if port == "" {
		panic("Port not specified!")
	}

  tc := &TripleClient {PortSpec: port}
	c := &CanbusClient {TripleClient: tc}

  c.options.bus1Enabled = false
  c.options.bus2Enabled = false
  c.options.bus3Enabled = false

	canCh, infoCh := c.TripleClient.OpenChannels()


	g := gocui.NewGui()
	c.g = g

	if err := g.Init(); err != nil {
		panic(err)
	}
	defer g.Close()

	g.SetLayout(c.layout)
  if err := c.keybindings(g); err != nil {
		panic(err)
	}

	g.SelBgColor = gocui.ColorGreen
	g.SelFgColor = gocui.ColorBlack

  g.SetCurrentView("main")

	go c.initCanChannel(canCh)
	go c.initInfoChannel(infoCh)

	err = g.MainLoop()

  if err != nil && err != gocui.Quit {
		panic(err)
	}
}
