package main

import (
  "os"
  "fmt"
  "bufio"
  "log"
  "strings"
  "strconv"
  "encoding/binary"
  // "time"
  "github.com/tarm/serial"
)

func main() {

  cfg := &serial.Config{Name: "/dev/ptyp6", Baud: 115200}
  port, err := serial.OpenPort(cfg)

  file, err := os.Open("./canbuslog")

  if err != nil {
    log.Fatal(err)
  }

  defer file.Close()


  for {

    file.Seek(0, 0)
    scanner := bufio.NewScanner(file)

    for scanner.Scan() {

      line := scanner.Text()
      bytes := strings.Split(line, " ")

      buf := make([]byte, len(bytes))

      for i := range bytes {
        n, _ := strconv.ParseInt(bytes[i], 10, 8)

        b := make([]byte, 2)
        binary.LittleEndian.PutUint16(b, uint16(n))

        buf[i] = b[0]
      }

      // time.Sleep(1*time.Second)
      port.Write(buf)

      fmt.Printf("%+v\n", buf)

    }
  }
}
