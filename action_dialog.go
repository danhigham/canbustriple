package main

import (
  "github.com/danhigham/gocui"
  "strings"
  "log"
  "fmt"
)

type ActionButton struct {
  handler func(CANPacket)
  caption string
  shortcut gocui.Key
}

var viewName = "action_dialog"

var actions = [...]ActionButton {
  ActionButton{
    caption: fmt.Sprintf("Log"),
    shortcut: gocui.KeyCtrlL,
    handler: func(packet CANPacket){
      log.Printf("%v", packet)
  }},
  ActionButton{
    caption: "Ignore",
    shortcut: gocui.KeyCtrlI,
    handler: func(packet CANPacket){
      log.Printf("%v", packet)
  }}}

func (c *CanbusClient) showActionDialog(packet CANPacket) {
  maxX, maxY := c.g.Size()

  width := 80
  height := 15

  if v, err := c.g.SetView(viewName, maxX/2-(width/2), maxY/2-(height/2),
    maxX/2+(width/2), maxY/2+(height/2)); err != nil {

    if err != gocui.ErrorUnkView {
      panic(err)
    }

    fmt.Fprint(v, padRight(" Bus", " ", 10))
		fmt.Fprint(v, padRight("| Message ID", " ", 10))
		fmt.Fprint(v, padRight("| Data", " ", 50))
		fmt.Fprint(v, padRight("| Length", " ", 10))
    fmt.Fprintf(v, "\n%s\n", strings.Repeat("-", 80))
    fmt.Fprint(v, string(packet.lineEntry(true)))
    fmt.Fprint(v, "\n\n")

    for _, a := range actions {
      x, y := v.Cursor()
      log.Printf("*** %v *** %v",x ,y)
      v.Write([]byte(a.caption))

      v.WriteUnderlinedRune(5, 5, 'x')
      // termbox.SetCell(5, 5, 'x', termbox.AttrUnderline, termbox.AttrUnderline)
    }

    if err := c.g.SetCurrentView(viewName); err != nil {
      panic(err)
    }
  }

  c.g.SetKeybinding(viewName, gocui.KeyEsc, gocui.ModNone, func(g *gocui.Gui, v *gocui.View) error {
    c.g.DeleteView(viewName)
    c.g.SetCurrentView(c.mainView.Name())
    return nil
  })

}
