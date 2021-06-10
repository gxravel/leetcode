package main

import (
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"time"
)

//!+
func main() {
	conn, err := net.Dial("tcp", "localhost:8000")
	if err != nil {
		log.Fatal(err)
	}
	done := make(chan struct{})
	go func() {
		io.Copy(os.Stdout, conn) // NOTE: ignoring errors
		log.Println("done")
		done <- struct{}{} // signal the main goroutine
	}()

	time.AfterFunc(3*time.Second, func() {
		fmt.Println("writing is closed")
		conn.(*net.TCPConn).CloseWrite()
	})
	mustCopy(conn, os.Stdin)
	<-done // wait for background goroutine to finish
}

//!-

func mustCopy(dst io.Writer, src io.Reader) {
	if _, err := io.Copy(dst, src); err != nil {
		log.Fatal(err)
	}
}
