package main

import (
  "log"
  "fmt"
  "bufio"
  "encoding/json"
  "github.com/tarm/serial"
  // "strings"
)

type TripleInfo struct {
  Event   string `json:"event"`
  Name    string `json:"name"`
  Version string `json:"version"`
  Memory  string `json:"memory"`
}

type TripleClient struct {
  PortSpec  string
  port      *serial.Port
}

func (c *TripleClient) ensureConnection() {
  if c.port == nil {

    var err error

    cfg := &serial.Config{Name: c.PortSpec, Baud: 115200}
    c.port, err = serial.OpenPort(cfg)

    if err != nil {
      log.Fatal(err)
    }

  }
}

func (c* TripleClient) OpenChannel() chan []byte {
  var channel = make(chan []byte)

  return channel
}

func (c *TripleClient) SetBus(busId byte, enabled byte) {
  c.ensureConnection()
  _, err := c.port.Write([]byte{0x03,busId,enabled,0x0000,0x0000})
  if err != nil {
    log.Fatal(err)
  }
}

func (c *TripleClient) GetInfo() string {
  c.ensureConnection()

  _, err := c.port.Write([]byte{0x01,0x01})

  if err != nil {
          log.Fatal(err)
  }

  reader := bufio.NewReader(c.port)
  reply, err := reader.ReadBytes('\x0a')
  if err != nil {
    panic(err)
  }

  jsobj := reply[:len(reply)]
  var info TripleInfo
  err = json.Unmarshal(jsobj, &info)

	if err != nil {
		panic(err)
	}

  ret := fmt.Sprintf("\n\n  Event:    %s\n  Name:     %s\n  Version:  %s\n  Memory:   %s\n",
    info.Event, info.Name, info.Version, info.Memory)

  return ret
}
