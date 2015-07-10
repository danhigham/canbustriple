package main

import (
	"bufio"
	"encoding/json"
	"log"

	"github.com/tarm/serial"
)

type TripleInfo struct {
	Event   string `json:"event"`
	Name    string `json:"name"`
	Version string `json:"version"`
	Memory  string `json:"memory"`
}

type CANPacket struct {
	Bus       byte
	Mid0      byte
	Mid1      byte
	Data      []byte
	Length    byte
	BusStatus byte
	MessageID int
}

type TripleClient struct {
	PortSpec string
	port     *serial.Port
}

func (c *TripleClient) ensureConnection() error {
	if c.port == nil {

		var err error

		cfg := &serial.Config{Name: c.PortSpec, Baud: 115200}
		c.port, err = serial.OpenPort(cfg)

		if err != nil {
			return err
		}
	}

	return nil
}

func (c *TripleClient) peekConnection() (byte, error) {

	reader := bufio.NewReader(c.port)
	b, err := reader.Peek(1)

	if err != nil {
		return 0x00, err
	}

	return b[0], nil
}

func (c *TripleClient) OpenChannels() (error, chan CANPacket, chan TripleInfo) {

	var canChannel = make(chan CANPacket)
	var infoChannel = make(chan TripleInfo)

	err := c.ensureConnection()

	if err != nil {
		c.port.Close()
		return err, nil, nil
	}

	// send a request for info and then peek
	c.RequestInfo()
	_, err = c.peekConnection()

	if err != nil {
		c.port.Close()
		return err, nil, nil
	}

	go func() {

		for {

			reader := bufio.NewReader(c.port)
			m, err := reader.Peek(1)

			if err != nil {
				c.port.Close()
				panic(err)
			}

			log.Printf("%v", m[0])

			if m[0] == 0x03 { //CAN Packet

				buf := make([]byte, 14)
				_, err := reader.Read(buf)
				// buf, err := reader.ReadBytes('\x0d')

				if err != nil {
					c.port.Close()
					panic(err)
				}

				p := CANPacket{Bus: buf[1], Mid0: buf[2], Mid1: buf[3], Data: buf[4:12], Length: buf[12], BusStatus: buf[13]}
				p.MessageID = (int(p.Mid0) << 8) + int(p.Mid1)
				if p.Length > 0 {
					canChannel <- p
				}
			} else if m[0] == 0x7B { //JSON

				line, err := reader.ReadBytes('\x0d')

				if err != nil {
					c.port.Close()
					panic(err)
				}

				jsobj := line[:len(line)]

				var info TripleInfo
				err = json.Unmarshal(jsobj, &info)

				if err != nil {
					c.port.Close()
					panic(err)
				}

				infoChannel <- info

			} else {
				line, err := reader.ReadBytes('\x0d')
				log.Printf("%v", line[:len(line)])

				if err != nil {
					c.port.Close()
					panic(err)
				}
			}
		}
	}()

	return nil, canChannel, infoChannel
}

func (c *TripleClient) RequestInfo() {
	c.ensureConnection()
	_, err := c.port.Write([]byte{0x01, 0x01})

	if err != nil {
		panic(err)
	}
}

func (c *TripleClient) SetBus(busId byte, enabled byte) {
	c.ensureConnection()
	_, err := c.port.Write([]byte{0x03, busId, enabled, 0x0000, 0x0000})
	if err != nil {
		panic(err)
	}
}
