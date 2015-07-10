// Copyright 2014 The gocui Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"fmt"
	"log"
	"os"
	"sort"
	"strings"
	"time"

	"github.com/danhigham/gocui"
)

type CanbusClient struct {
	filterView   *gocui.View
	optionView   *gocui.View
	mainView     *gocui.View
	options      CanbusClientOptions
	TripleClient *TripleClient
	g            *gocui.Gui
	ShowAbout    bool
	PauseOutput  bool
	Packets      map[int]CANPacket
	ShowCompact  bool
	SelectedLine int
	Actions      [0]ActionButton
}

type CanbusClientOptions struct {
	bus1Enabled bool
	bus2Enabled bool
	bus3Enabled bool
}

func (c *CanbusClient) layout(g *gocui.Gui) error {

	if c.ShowAbout {
		c.showAboutDialog()
		return nil
	}

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

	if c.mainView, err = g.SetView("main", 20, 1, maxX, maxY-2); err != nil {

		if err != gocui.ErrorUnkView {
			return err
		}

		c.mainView.Autoscroll = true
		c.mainView.Overwrite = true

		// Initialise Canbus Triple connection
		err, canCh, infoCh := c.TripleClient.OpenChannels()

		go c.initCanChannel(canCh)
		go c.initInfoChannel(infoCh)

		if err != nil {
			return gocui.Quit
		}
	}

	if v, err := g.SetView("headers", 20, -1, maxX, 1); err != nil {
		if err != gocui.ErrorUnkView {
			return err
		}

		fmt.Fprintf(v, padRight(" Bus", " ", 13))
		fmt.Fprintf(v, padRight("| Message ID", " ", 14))
		fmt.Fprintf(v, padRight("| Data", " ", 76))
		fmt.Fprintf(v, padRight("| Length", " ", 12))
	}

	c.createMainMenu(g)

	return nil
}

func (c *CanbusClient) togglePause(g *gocui.Gui, v *gocui.View) error {
	c.PauseOutput = !c.PauseOutput
	return nil
}

func (c *CanbusClient) writeLoggingOptions() {
	if c.options.bus1Enabled {
		c.TripleClient.SetBus(0x01, 0x01)
	} else {
		c.TripleClient.SetBus(0x01, 0x00)
	}
	if c.options.bus2Enabled {
		c.TripleClient.SetBus(0x02, 0x01)
	} else {
		c.TripleClient.SetBus(0x02, 0x00)
	}
	if c.options.bus3Enabled {
		c.TripleClient.SetBus(0x03, 0x01)
	} else {
		c.TripleClient.SetBus(0x03, 0x00)
	}
}

func (c *CanbusClient) writeOptionsPane() error {
	v := c.optionView
	v.Clear()

	fmt.Fprint(v, "Options\n-------------------\n")

	fmt.Fprintf(v, "\n%s%+12s", "Bus 1", "")
	if c.options.bus1Enabled {
		fmt.Fprint(v, "\u2714")
	} else {
		fmt.Fprint(v, " ")
	}
	fmt.Fprintf(v, "\n%s%+12s", "Bus 2", "")
	if c.options.bus2Enabled {
		fmt.Fprint(v, "\u2714")
	} else {
		fmt.Fprint(v, " ")
	}
	fmt.Fprintf(v, "\n%s%+12s", "Bus 3", "")
	if c.options.bus3Enabled {
		fmt.Fprint(v, "\u2714")
	} else {
		fmt.Fprint(v, " ")
	}

	_, cy := v.Cursor()
	if cy < 3 {
		v.SetCursor(0, 2)
	}

	return nil
}

func (c *CanbusClient) setOptions(g *gocui.Gui, v *gocui.View) error {
	if v != nil {
		_, cy := v.Cursor()

		if cy == 3 {
			c.options.bus1Enabled = !c.options.bus1Enabled
		}
		if cy == 4 {
			c.options.bus2Enabled = !c.options.bus2Enabled
		}
		if cy == 5 {
			c.options.bus3Enabled = !c.options.bus3Enabled
		}

		c.writeOptionsPane()
		c.writeLoggingOptions()
	}

	return nil
}

func (c *CanbusClient) switchToOptions(g *gocui.Gui, v *gocui.View) error {
	g.SetCurrentView("side-options")
	c.optionView.SetCursor(0, 3)

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

func (canPacket *CANPacket) lineEntry(hideDataString bool) []byte {

	bus := fmt.Sprintf("%v", canPacket.Bus)
	messageId := fmt.Sprintf("%v", canPacket.MessageID)
	length := fmt.Sprintf("%v", canPacket.Length)
	hexdata := make([]string, 8)

	data := make([]rune, 8)

	for i := 0; i < 8; i++ {
		hexdata[i] = fmt.Sprintf("%02X", canPacket.Data[i])
		if canPacket.Data[i] > 31 && canPacket.Data[i] < 127 {
			data[i] = rune(canPacket.Data[i])
		} else {
			data[i] = '.'
		}
	}

	var s string

	// format packet for display
	if hideDataString {

		s = fmt.Sprintf(
			" %s| %s| %s| %s",
			padRight(bus, " ", 9),
			padRight(messageId, " ", 8),
			padRight(strings.Join(hexdata, "  "), " ", 48),
			padRight(length, " ", 10))

	} else {

		s = fmt.Sprintf(
			" %s| %s| %s| %s| %s",
			padRight(bus, " ", 12),
			padRight(messageId, " ", 12),
			padRight(strings.Join(hexdata, "  "), " ", 48),
			padRight(string(data), " ", 24),
			padRight(length, " ", 12))

	}

	return []byte(s)
}

func (c *CanbusClient) initCanChannel(ch chan CANPacket) {

	c.Packets = make(map[int]CANPacket)

	for {

		canPacket := <-ch
		c.Packets[canPacket.MessageID] = canPacket

		if c.ShowCompact {
			c.mainView.DirtyClear()
			c.drawCompactView()
			c.mainView.SetCursor(0, c.SelectedLine)

		} else {

			c.mainView.Write(canPacket.lineEntry(false))
			c.mainView.Write([]byte("\n"))

		}

		if !c.PauseOutput {
			c.g.Flush()
		}
	}
}

func (c *CanbusClient) packetKeys() []int {
	var keys []int

	for k := range c.Packets {
		keys = append(keys, k)
	}

	sort.Ints(keys)
	return keys
}

func (c *CanbusClient) drawCompactView() {

	// sort Keys
	keys := c.packetKeys()

	//order packets here
	for _, k := range keys {

		p := c.Packets[k]
		c.mainView.Write(p.lineEntry(false))
		c.mainView.Write([]byte("\n"))

	}

}

func (c *CanbusClient) initInfoChannel(ch chan TripleInfo) {

	for {
		info := <-ch

		txt := fmt.Sprintf("\n\n  Event:    %s\n  Name:     %s\n  Version:  %s\n  Memory:   %s\n",
			info.Event, info.Name, info.Version, info.Memory)

		maxX, maxY := c.g.Size()

		if v, err := c.g.SetView("triple-info", maxX/2-15, maxY/2-5, maxX/2+15, maxY/2+5); err != nil {

			if err != gocui.ErrorUnkView {
				panic(err)
			}

			v.SetOrigin(-10, -10)
			v.SetCursor(20, 20)

			fmt.Fprint(v, txt)

			if err := c.g.SetCurrentView("triple-info"); err != nil {
				panic(err)
			}

			c.g.Flush()
		}
	}
}

func main() {

	var err error

	f, err := os.OpenFile("./canbustriple.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		panic(fmt.Sprintf("error opening file: %v", err))
	}
	defer f.Close()

	log.SetOutput(f)

	var port = os.Args[1]

	if port == "" {
		panic("Port not specified!")
	}

	tc := &TripleClient{PortSpec: port}
	c := &CanbusClient{TripleClient: tc}

	c.options.bus1Enabled = false
	c.options.bus2Enabled = false
	c.options.bus3Enabled = false
	c.PauseOutput = false
	c.ShowCompact = false
	c.SelectedLine = 0
	c.ShowAbout = true

	g := gocui.NewGui()
	c.g = g

	go func() {
		time.Sleep(2 * time.Second)
		c.ShowAbout = false
		g.Flush()
	}()

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

	err = g.MainLoop()

	if err != nil && err != gocui.Quit {
		panic(err)
	}

	fmt.Println("Hello")
}
