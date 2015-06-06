// Copyright 2014 The gocui Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"os"
	"log"
  "fmt"
	"strings"
	"sort"
	"github.com/jroimartin/gocui"
)

type CanbusClient struct {
  filterView  	*gocui.View
  optionView  	*gocui.View
  mainView    	*gocui.View
  options     	CanbusClientOptions
	TripleClient	*TripleClient
	g 						*gocui.Gui
	PauseOutput		bool
	Packets				map[int]CANPacket
	ShowCompact 	bool
}

type CanbusClientOptions struct {
  bus1Enabled bool
  bus2Enabled bool
  bus3Enabled bool
}

func PadRight(str, pad string, length int) string {
    for {
        str += pad
        if len(str) > length {
            return str[0:length]
        }
    }
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

	if v, err := g.SetView("headers", 20, -1, maxX, 1); err != nil {
    if err != gocui.ErrorUnkView {
		  return err
    }

		fmt.Fprintf(v, PadRight(" Bus", " ", 13))
		fmt.Fprintf(v, PadRight("| Message ID", " ", 14))
		fmt.Fprintf(v, PadRight("| Data", " ", 76))
		fmt.Fprintf(v, PadRight("| Length", " ", 12))
  }

	if v, err := g.SetView("cmdline", -1, maxY-2, maxX, maxY); err != nil {
    if err != gocui.ErrorUnkView {
		  return err
    }

		fmt.Fprintf(v, "O: Set Options\t\t\t")
    fmt.Fprintf(v, "P: Pause\t\t\t")
    fmt.Fprintf(v, "V: Compact\t\t\t")
    fmt.Fprintf(v, "M: Send CAN Message\t\t\t")
    fmt.Fprintf(v, "F: Add Filter\t\t\t")
    fmt.Fprintf(v, "I: Get Sys Info\t\t\t")
    fmt.Fprintf(v, "C: Quit\t\t\t")
  }

  return nil
}

func quit(g *gocui.Gui, v *gocui.View) error {
	return gocui.Quit
}

func (c *CanbusClient) togglePause(g *gocui.Gui, v *gocui.View) error {
	c.PauseOutput = !c.PauseOutput
	return nil
}

func (c *CanbusClient) toggleCompactView(g *gocui.Gui, v *gocui.View) error {
	c.ShowCompact = !c.ShowCompact
	return nil
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

	if err := g.SetKeybinding("", gocui.KeyCtrlP, gocui.ModNone, c.togglePause); err != nil {
		return err
	}

	return nil
}

func (canPacket *CANPacket) lineEntry(compact bool) []byte {

	bus := fmt.Sprintf("%v", canPacket.Bus)
	messageId := fmt.Sprintf("%v", canPacket.MessageID)
	length := fmt.Sprintf("%v", canPacket.Length)
	hexdata := make([]string, 8)

	data := make([]rune, 8)

	for i := 0; i< 8; i ++ {
		hexdata[i] = fmt.Sprintf("%02X", canPacket.Data[i])
		if canPacket.Data[i] > 31 && canPacket.Data[i] < 127 {
			data[i] = rune(canPacket.Data[i])
		} else {
			data[i] = '.'
		}
	}

	// format packet for display
	s := fmt.Sprintf(
		" %s| %s| %s| %s| %s",
		PadRight(bus, " ", 12),
		PadRight(messageId, " ", 12),
		PadRight(strings.Join(hexdata, "  "), " ", 48),
		PadRight(string(data), " ", 24),
		PadRight(length, " ", 12))

	return []byte(s)
}

func (c *CanbusClient) initCanChannel(ch chan CANPacket) {

	c.Packets = make(map[int]CANPacket)

	for {
    canPacket := <- ch
		c.Packets[canPacket.MessageID] = canPacket

		if c.ShowCompact {

			c.mainView.Clear()
			c.drawCompactView()

		} else {

			c.mainView.Write(canPacket.lineEntry(false))
			c.mainView.Write([]byte("\n"))
		}

		if !c.PauseOutput { c.g.Flush() }
  }
}

func (c *CanbusClient) drawCompactView() {

	// sort Keys

	var keys []int
	for k := range c.Packets {
		keys = append(keys, k)
	}
	sort.Ints(keys)

	//order packets here
	for _, k := range keys {

		p := c.Packets[k]
		c.mainView.Write(p.lineEntry(true))
		c.mainView.Write([]byte("\n"))

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
	c.PauseOutput = false
	c.ShowCompact = false

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
