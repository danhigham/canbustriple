package main

import (
	"github.com/danhigham/gocui"
  "strings"
  "fmt"
)

func (c *CanbusClient) showActionDialog(packet CANPacket) {
  maxX, maxY := c.g.Size()

  width := 80
  height := 15

  if v, err := c.g.SetView("action_dialog", maxX/2-(width/2), maxY/2-(height/2),
    maxX/2+(width/2), maxY/2+(height/2)); err != nil {

    if err != gocui.ErrorUnkView {
      panic(err)
    }

    v.SetOrigin(-10,-10)
    v.SetCursor(20,20)

    fmt.Fprintf(v, padRight(" Bus", " ", 10))
		fmt.Fprintf(v, padRight("| Message ID", " ", 10))
		fmt.Fprintf(v, padRight("| Data", " ", 50))
		fmt.Fprintf(v, padRight("| Length", " ", 10))
    fmt.Fprintf(v, "\n%s\n", strings.Repeat("-", 80))
    fmt.Fprintf(v, string(packet.lineEntry(true)))

    if err := c.g.SetCurrentView("action_dialog"); err != nil {
      panic(err)
    }
  }

}
