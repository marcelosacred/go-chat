package main

import (
    "net"
    "bufio"
    "fmt"
)

func handleConnection(conn net.Conn) {
    defer conn.Close()

    reader := bufio.NewReader(conn)
    for {
        msg, err := reader.ReadString('\n')
        if err != nil {
            fmt.Println("client disconnected")
            return
        }
        fmt.Printf("got message: %s", msg)
        conn.Write([]byte("echo: " + msg))
    }
}