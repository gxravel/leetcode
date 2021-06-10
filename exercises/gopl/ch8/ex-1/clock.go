package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"net"
	"time"
)

var port int
var tz string
var loc *time.Location

func init() {
	flag.IntVar(&port, "port", 8000, "port")
	flag.StringVar(&tz, "TZ", "Europe/Moscow", "timezone")
}

func handleConn(c net.Conn) {
	defer c.Close()
	for {
		_, err := io.WriteString(c, time.Now().In(loc).Format("15:04:05\n"))
		if err != nil {
			return // e.g., client disconnected
		}
		time.Sleep(1 * time.Second)
	}
}

func main() {
	flag.Parse()
	var err error
	loc, err = time.LoadLocation(tz)
	if err != nil {
		log.Fatal(err)
	}
	listener, err := net.Listen("tcp", fmt.Sprintf("localhost:%d", port))
	if err != nil {
		log.Fatal(err)
	}
	//!+
	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Print(err) // e.g., connection aborted
			continue
		}
		go handleConn(conn) // handle connections concurrently
	}
	//!-
}
