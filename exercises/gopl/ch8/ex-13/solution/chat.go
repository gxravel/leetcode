/*
Упражнение 8.13. Заставьте сервер отключать простаивающих клиентов, которые
не прислали ни одного сообщения за последние 5 минут.
Указание: вызов co n n . C lose () в другой go-подпрограмме деблокирует активный
вызов Read, такой, как выполняемый вызовом in p u t .S c a n ().
*/
package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"time"
)

//!+broadcaster
type client struct {
	message chan<- string // an outgoing message channel
	name    string
}

var (
	entering = make(chan client)
	leaving  = make(chan client)
	messages = make(chan string) // all incoming client messages
)

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

		case cli := <-entering:
			if len(clients) == 0 {
				cli.message <- "there are no other users currently\n"
			} else {
				cli.message <- "the other users:"
				for other := range clients {
					cli.message <- other.name
				}
				cli.message <- "\n"
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
	ch := make(chan string) // outgoing client messages
	written := make(chan struct{})
	go clientWriter(conn, ch)

	who := conn.RemoteAddr().String()
	ch <- "You are " + who
	messages <- who + " has arrived"
	entering <- client{message: ch, name: who}

	d := 5 * time.Minute
	t := time.NewTimer(d)
	go func() {
		for {
			select {
			case <-t.C:
				leaving <- client{message: ch, name: who}
				messages <- who + " has left"
				conn.Close()
				return
			case <-written:
				t.Reset(d)
			}
		}
	}()

	input := bufio.NewScanner(conn)
	for input.Scan() {
		messages <- who + ": " + input.Text()
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
