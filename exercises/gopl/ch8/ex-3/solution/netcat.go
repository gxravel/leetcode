/*
Упражнение 8.3. В программе n e tc a t3 значение интерфейса conn имеет конкрет­
ный тип * n et .TCPConn, который представляет TCP-соединение. ТСР-соединение
состоит из двух половин, которые могут быть закрыты независимо с использованием
методов C loseR ead и C lo seW rite. Измените главную go-подпрограмму n e tc a t3
так, чтобы она закрывала только записывающую половину соединения, так, чтобы
программа продолжала выводить последние эхо от сервера r e v e r b l даже после того,
как стандартный ввод будет закрыт. (Сделать это для сервера re v e rb 2 труднее; см.
упражнение 8.4.)
*/
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

	time.AfterFunc(5*time.Second, func() {
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
