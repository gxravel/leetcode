Упражнение 8.1. Измените программу clock2 таким образом, чтобы она прини­
мала номер порта, и напишите программу clockwall, которая действует в качестве
клиента нескольких серверов одновременно, считывая время из каждого и выводя
результаты в виде таблицы, сродни настенным часам, которые можно увидеть в не­
которых офисах. Если у вас есть доступ к географически разнесенным компьютерам,
запустите экземпляры серверов удаленно; в противном случае запустите локальные
экземпляры на разных портах с поддельными часовыми поясами.

$ TZ=US/Eastern ./clock2 -port 8010 &
$ TZ=Asia/Tokyo ./clock2 -port 8020 &
$ TZ=Europe/London ./clock2 -port 8030 &
$ clockwall NewYork=localhost:8010 London=localhost:8030 Tokyo=localhost:8020

clock2.go

package main

import (
	"io"
	"log"
	"net"
	"time"
)

func handleConn(c net.Conn) {
	defer c.Close()
	for {
		_, err := io.WriteString(c, time.Now().Format("15:04:05\n"))
		if err != nil {
			return // e.g., client disconnected
		}
		time.Sleep(1 * time.Second)
	}
}

func main() {
	listener, err := net.Listen("tcp", "localhost:8000")
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
