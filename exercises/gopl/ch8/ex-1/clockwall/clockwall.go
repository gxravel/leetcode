package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"strings"
	"time"
)

func main() {
	flag.Parse()
	fmt.Println("args: ", flag.Args())
	for _, address := range flag.Args() {
		go func(address string) {
			iEquality := strings.Index(address, "=")
			conn, err := net.Dial("tcp", address[iEquality+1:])
			if err != nil {
				log.Fatal(err)
			}
			defer conn.Close()
			mustCopy(os.Stdout, conn)
		}(address)
	}
	time.Sleep(1 * time.Minute)
}

func mustCopy(dst io.Writer, src io.Reader) {
	if _, err := io.Copy(dst, src); err != nil {
		log.Fatal(err)
	}
}
