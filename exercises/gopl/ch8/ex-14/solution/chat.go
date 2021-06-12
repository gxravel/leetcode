/*
Упражнение 8.14. Измените сетевой протокол чат-сервера так, чтобы каждый
клиент предоставлял при подключении свое имя. Используйте это имя вместо сетево­
го адреса в префиксе сообщения.
*/
package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"strings"
	"time"
	"unicode/utf8"
)

//!+broadcaster
type client struct {
	message chan<- string // an outgoing message channel
	name    string
}

var (
	entering    = make(chan client)
	leaving     = make(chan client)
	messages    = make(chan string) // all incoming client messages
	checkUnique = make(chan client)
)

func isNameUnique(clients map[client]bool, name string) bool {
	for cli := range clients {
		if cli.name == name {
			return false
		}
	}
	return true
}

func broadcaster() {
	clients := make(map[client]bool) // all connected clients
	for {
		select {
		case msg := <-messages:
			// Broadcast incoming message to all
			// clients' outgoing message channels.
			for cli := range clients {
				cli.message <- msg
			}
		case cli := <-checkUnique:
			if len(clients) > 0 && !isNameUnique(clients, cli.name) {
				cli.message <- "the login is already in use"
			} else {
				cli.message <- "OK"
			}
		case cli := <-entering:
			if len(clients) == 0 {
				cli.message <- "there are no other users currently\n"
			} else {
				cli.message <- "the other users:"
				for other := range clients {
					cli.message <- other.name
				}
				cli.message <- ""
			}
			clients[cli] = true

		case cli := <-leaving:
			delete(clients, cli)
			close(cli.message)
		}
	}
}

//!-broadcaster

//!+handleConn
func handleConn(conn net.Conn) {
	var cli client
	ch := make(chan string) // outgoing client messages
	cli.message = ch

	var name = make([]byte, 20)
	for {
		fmt.Fprint(conn, "login: ") // NOTE: ignoring network errors
		n, err := conn.Read(name)
		if err != nil {
			conn.Close()
			return
		}

		cli.name = strings.TrimRight(string(name[:n]), "\r\n")
		if utf8.RuneCountInString(cli.name) < 3 {
			fmt.Fprintln(conn, "login must be at least 3 symbols length")
			continue
		}
		checkUnique <- cli
		msg := <-ch
		fmt.Fprintln(conn, msg) // NOTE: ignoring network errors
		if msg == "OK" {
			break
		}
	}

	written := make(chan struct{})
	go clientWriter(conn, ch)

	messages <- cli.name + " has arrived"
	entering <- cli

	d := 5 * time.Minute
	t := time.NewTimer(d)
	go func() {
		for {
			select {
			case <-t.C:
				leaving <- cli
				messages <- cli.name + " has left"
				conn.Close()
				return
			case <-written:
				t.Reset(d)
			}
		}
	}()

	input := bufio.NewScanner(conn)
	for input.Scan() {
		messages <- cli.name + ": " + input.Text()
		written <- struct{}{}
	}
	if t.Stop() {
		t.Reset(0)
	}
}

func clientWriter(conn net.Conn, ch <-chan string) {
	for msg := range ch {
		fmt.Fprintln(conn, msg) // NOTE: ignoring network errors
	}
}

//!-handleConn

//!+main
func main() {
	listener, err := net.Listen("tcp", "localhost:8000")
	if err != nil {
		log.Fatal(err)
	}
	go broadcaster()
	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Print(err)
			continue
		}
		go handleConn(conn)
	}
}

//!-main
