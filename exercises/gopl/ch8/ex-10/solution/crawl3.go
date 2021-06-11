/*
Упражнение 8.10. Запросы HTTP могут быть отменены с помощью закрытия не­
обязательного канала Cancel в структуре http.Request. Измените веб-сканер из
раздела 8.6 так, чтобы он поддерживал отмену.
Указание. Функция h ttp .G et не позволяет настроить R equest. Вместо этого соз­
дайте запрос с использованием h ttp .N ew R eq u est, установите его поле C ancel и
выполните запрос с помощью вызова h ttp .D e f a u ltC lie n t .D o (req ).
*/
package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"

	"golang.org/x/net/html"
)

var maxDepth int

type Link struct {
	url   string
	depth int
}

func init() {
	flag.IntVar(&maxDepth, "depth", 3, "set depth of crawling")
}

func createLinks(urls []string, depth int) []Link {
	list := make([]Link, 0, len(urls))
	for _, url := range urls {
		list = append(list, Link{url, depth})
	}
	return list
}

var done = make(chan struct{})

func cancelled() bool {
	select {
	case <-done:
		return true
	default:
		return false
	}
}

func Extract(ctx context.Context, url string) ([]string, error) {
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, err
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != http.StatusOK {
		resp.Body.Close()
		return nil, fmt.Errorf("getting %s: %s", url, resp.Status)
	}

	doc, err := html.Parse(resp.Body)
	resp.Body.Close()
	if err != nil {
		return nil, fmt.Errorf("parsing %s as HTML: %v", url, err)
	}

	var links []string
	visitNode := func(n *html.Node) {
		if n.Type == html.ElementNode && n.Data == "a" {
			for _, a := range n.Attr {
				if a.Key != "href" {
					continue
				}
				link, err := resp.Request.URL.Parse(a.Val)
				if err != nil {
					continue // ignore bad URLs
				}
				links = append(links, link.String())
			}
		}
	}
	forEachNode(doc, visitNode, nil)
	return links, nil
}

func forEachNode(n *html.Node, pre, post func(n *html.Node)) {
	if pre != nil {
		pre(n)
	}
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		forEachNode(c, pre, post)
	}
	if post != nil {
		post(n)
	}
}

func crawl(ctx context.Context, link Link) []Link {
	if cancelled() {
		return nil
	}
	fmt.Println(link.depth, link.url)
	urls, err := Extract(ctx, link.url)
	if err != nil {
		log.Print(err)
	}

	return createLinks(urls, link.depth+1)
}

//!+
func main() {
	flag.Parse()
	worklist := make(chan []Link)  // lists of URLs, may have duplicates
	unseenLinks := make(chan Link) // de-duplicated URLs

	// Add command-line arguments to worklist.
	go func() {
		worklist <- createLinks(flag.Args(), 1)
	}()

	go func() {
		os.Stdin.Read(make([]byte, 1)) // Чтение одного байта
		close(done)
	}()

	ctx, cancel := context.WithCancel(context.Background())

	// Create 20 crawler goroutines to fetch each unseen link.
	for i := 0; i < 20; i++ {
		go func() {
			for link := range unseenLinks {
				foundLinks := crawl(ctx, link)
				go func() { worklist <- foundLinks }()
			}
		}()
	}

	// The main goroutine de-duplicates worklist items
	// and sends the unseen ones to the crawlers.
	seen := make(map[string]bool)
	for {
		select {
		case list := <-worklist:
			for _, link := range list {
				if !seen[link.url] && link.depth <= maxDepth {
					seen[link.url] = true
					unseenLinks <- link
				}
			}
		case <-done:

			fmt.Println("case done")
			cancel()
			close(worklist)
			close(unseenLinks)
			for range worklist {

			}
			for range unseenLinks {

			}
			return
		}
	}
}

//!-
