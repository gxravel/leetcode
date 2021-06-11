/*
Упражнение 8.8. Используя инструкцию select, добавьте к эхо-серверу из раз­
дела 8.3 тайм-аут, чтобы он отключал любого клиента, который ничего не передает в
течение 10 секунд.
*/
package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"strings"
	"time"
)

func echo(c net.Conn, shout string, delay time.Duration) {
	fmt.Fprintln(c, "\t", strings.ToUpper(shout))
	time.Sleep(delay)
	fmt.Fprintln(c, "\t", shout)
	time.Sleep(delay)
	fmt.Fprintln(c, "\t", strings.ToLower(shout))
}

//!+
func handleConn(c net.Conn) {
	defer c.Close()
	input := bufio.NewScanner(c)
	echoed := make(chan struct{})
	go func() {
		for input.Scan() {
			go echo(c, input.Text(), 1*time.Second)
			echoed <- struct{}{}
		}
	}()
	t := time.NewTimer(10 * time.Second)
	for {
		select {
		case <-t.C:
			return
		case <-echoed:
			t.Reset(10 * time.Second)
		}
	}
	// NOTE: ignoring potential errors from input.Err()
}

//!-

func main() {
	l, err := net.Listen("tcp", "localhost:8000")
	if err != nil {
		log.Fatal(err)
	}
	for {
		conn, err := l.Accept()
		if err != nil {
			log.Print(err) // e.g., connection aborted
			continue
		}
		go handleConn(conn)
	}
}
