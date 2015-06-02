package main

import (
  "os"
  "fmt"
  "bufio"
  "log"
  "strings"
  // "encoding/binary"
  // "github.com/tarm/serial"
)

func main() {

  // cfg := &serial.Config{Name: "/dev/ptyp6", Baud: 115200}
  // port, err := serial.OpenPort(cfg)

  file, err := os.Open("./canbuslog")

  if err != nil {
    log.Fatal(err)
  }

  defer file.Close()

  scanner := bufio.NewScanner(file)

  for scanner.Scan() {
    line := scanner.Text()
    bytes := strings.Split(line, " ")
    
    fmt.Printf("%+v", bytes)
    // buf := make([]byte, len(bytes))

    for i := range bytes {
      // nbyte := []uint8(bytes[i])

      // binary.PutUvarint(buf, nbyte)

      // fmt.Printf("%+v", nbyte)
      // port.Write(buf)
    }


    // fmt.Println(scanner.Text())
  }

}
