package main

import (
  "bufio"
  "bytes"
  "fmt"
  "io"
  "log"
  "net"
  "os"
)

func main() {
  messages := make(chan string)
  inputs := make(chan string)

  psock, err := net.Listen("tcp", ":5000")
  if err != nil {
    return
  }

  conn, err := psock.Accept()
  if err != nil {
    return
  }

  go handleMessages(messages)
  go handleInput(inputs, conn)
  go awaitInput(inputs)

  for {
    line, err := bufio.NewReader(conn).ReadBytes('\n')
    messages <- string(line)
    if err != nil {
      return
    }
  }
}

func awaitInput(inputs chan string) {
  scanner := bufio.NewScanner(os.Stdin)
  for scanner.Scan() {
    stdinText := scanner.Text()
    inputs <- stdinText
  }

  if err := scanner.Err(); err != nil {
    log.Fatal(err)
  }
}

func handleInput(inputs chan string, conn net.Conn) {
  for i := range inputs {
    io.Copy(conn, bytes.NewBufferString(i))
  }
}

func handleMessages(messages chan string) {
  for m := range messages {
    go handleMessage(m)
  }
}

func handleMessage(message string) {
  fmt.Println(message)
}
