/*
Упражнение 8.4. Модифицируйте сервер re v e rb 2 так, чтобы он использовал по
одному объекту s y n c . W aitG roup для каждого соединения для подсчета количества
активных go-подпрограмм echo. Когда он обнуляется, закрывайте пишущую полови­
ну TCP-соединения, как описано в упражнении 8.3. Убедитесь, что вы изменили кли­
ентскую программу n e tc a t3 из этого упражнения так, чтобы она ожидала последние
ответы от параллельных go-подпрограмм сервера даже после закрытия стандартного
ввода.
*/
package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"strings"
	"sync"
	"time"
)

func echo(c net.Conn, shout string, delay time.Duration, wg *sync.WaitGroup) {
	defer wg.Done()
	fmt.Fprintln(c, "\t", strings.ToUpper(shout))
	time.Sleep(delay)
	fmt.Fprintln(c, "\t", shout)
	time.Sleep(delay)
	fmt.Fprintln(c, "\t", strings.ToLower(shout))
}

//!+
func handleConn(c net.Conn) {
	input := bufio.NewScanner(c)
	var wg sync.WaitGroup
	var firstIteration = true

	for input.Scan() {
		wg.Add(1)
		if firstIteration {
			go func() {
				wg.Wait()
				c.(*net.TCPConn).CloseWrite()
			}()
			firstIteration = false
		}
		go echo(c, input.Text(), 1*time.Second, &wg)
	}

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
