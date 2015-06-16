package main

import (
  "github.com/danhigham/gocui"
  "strings"
  "bytes"
  "log"
  "fmt"
)

type ActionButton struct {
  handler func(*gocui.Gui, *gocui.View) error
  caption string
  shortcut gocui.Key
}

var viewName = "action_dialog"

var actions = [...]ActionButton {
  ActionButton{
    caption: "Label",
    shortcut: gocui.KeyCtrlL,
    handler: func(g *gocui.Gui, v *gocui.View) error {
      log.Printf("Label Packet")
      return nil
  }},
  ActionButton{
    caption: "Ignore",
    shortcut: gocui.KeyCtrlI,
    handler: func(g *gocui.Gui, v *gocui.View) error {
      log.Printf("Ignore Packet")
      return nil
  }},
  ActionButton{
    caption: "Send",
    shortcut: gocui.KeyCtrlS,
    handler: func(g *gocui.Gui, v *gocui.View) error {
      log.Printf("Send Packet")
      return nil
  }},
  ActionButton{
    caption: "Cancel",
    shortcut: gocui.KeyCtrlO,
    handler: func(g *gocui.Gui, v *gocui.View) error {
      g.DeleteView(viewName)
      g.SetCurrentView("main")
      return nil
  }}}

func (c *CanbusClient) showActionDialog(packet CANPacket) {

  maxX, maxY := c.g.Size()

  width := 80
  height := 10

  if v, err := c.g.SetView(viewName, maxX/2-(width/2), maxY/2-(height/2),
    maxX/2+(width/2), maxY/2+(height/2)); err != nil {

    v.Overwrite = true
    v.Underlined = []gocui.Pos{{1, 0}, {2, 3}, {4, 5}}

    if err != gocui.ErrorUnkView {
      panic(err)
    }

    fmt.Fprint(v, padRight(" Bus", " ", 10))
		fmt.Fprint(v, padRight("| Message ID", " ", 10))
		fmt.Fprint(v, padRight("| Data", " ", 50))
		fmt.Fprint(v, padRight("| Length", " ", 10))
    fmt.Fprintf(v, "\n%s\n", strings.Repeat("-", 80))
    fmt.Fprint(v, string(packet.lineEntry(true)))
    fmt.Fprint(v, "\n\n\n\n\n")

    var btnBuffer bytes.Buffer
    buttonWidth := 0

    for i, a := range actions {
      actions[i].caption = fmt.Sprintf("[%v]", a.caption)
      buttonWidth += len(actions[i].caption)
    }

    spc := (width - buttonWidth) / len(actions)

    for _, a := range actions {
      pleft := padLeft(a.caption, " ", (spc / 2) + len(a.caption))
      padded := padRight(pleft, " ", (spc / 2) + len(pleft))

      btnBuffer.WriteString(padded)
    }

    v.Write(btnBuffer.Bytes())

    if err := c.g.SetCurrentView(viewName); err != nil {
      panic(err)
    }
  }

  for _, a := range actions {
    c.g.SetKeybinding(viewName, a.shortcut, gocui.ModNone, a.handler)
  }

}
