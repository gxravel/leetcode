/*
Упражнение 8.7. Напишите параллельную программу, которая создает локальное
зеркало веб-сайта, загружая все доступные страницы и записывая их в каталог на
локальном диске. Выбираться должны только страницы в пределах исходного домена
(например, g o la n g .o rg ). URL в страницах зеркала должны при необходимости быть
изменены таким образом, чтобы они ссылались на зеркальную страницу, а не на ори­
гинал.
*/
package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"golang.org/x/net/html"
)

var site string
var mirror string

var maxDepth int

type Link struct {
	url   string
	depth int
}

func init() {
	flag.IntVar(&maxDepth, "depth", 3, "set depth of crawling")
	flag.StringVar(&site, "site", "golang.org", "a site to be mirrored")
	flag.StringVar(&mirror, "mirror", "mirror.com", "a mirror site")
}

func Extract(url string) ([]string, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != http.StatusOK {
		resp.Body.Close()
		return nil, fmt.Errorf("getting %s: %s", url, resp.Status)
	}
	if !strings.EqualFold(resp.Header.Get("content-type"), "text/html; charset=UTF-8") {
		resp.Body.Close()
		return nil, nil
	}

	doc, err := html.Parse(resp.Body)
	resp.Body.Close()
	if err != nil {
		return nil, fmt.Errorf("parsing %s as HTML: %v", url, err)
	}

	var links []string
	visitNode := func(n *html.Node) {
		if n.Type == html.ElementNode && n.Data == "a" {
			for i, a := range n.Attr {
				if a.Key != "href" {
					continue
				}
				link, err := resp.Request.URL.Parse(a.Val)
				if err != nil {
					continue // ignore bad URLs
				}
				if strings.Contains(link.String(), site) {
					links = append(links, link.String())
					n.Attr[i].Val = strings.Replace(link.String(), site, mirror, 1)
				}
			}
		}
	}
	forEachNode(doc, visitNode, nil)
	i := strings.Index(url, "//") + 2
	j := i + strings.Index(url[i:], "/")
	var dir, fname string
	if len(url)-1 == j {
		fname = mirror + ".html"
		k := strings.Index(url, site)
		if i != k {
			fname = url[i:k] + fname
		}
	} else {
		j++
		dir = filepath.Dir(url[j:])
		fname = filepath.Base(url[j:])
		if filepath.Ext(fname) == "" {
			fname += ".html"
		}
	}
	if dir != "" {
		err = os.MkdirAll(dir, 0755)
		if err != nil {
			return nil, err
		}
	}
	file, err := os.Create(filepath.Join(dir, fname))
	if err != nil {
		return nil, err
	}
	err = html.Render(file, doc)
	if err != nil {
		return nil, err
	}
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

func createLinks(urls []string, depth int) []Link {
	list := make([]Link, 0, len(urls))
	for _, url := range urls {
		list = append(list, Link{url, depth})
	}
	return list
}

func crawl(link Link) []Link {
	urls, err := Extract(link.url)
	if err != nil {
		log.Print(err)
	}
	if urls != nil {
		fmt.Println(link.depth, link.url)
	}
	return createLinks(urls, link.depth+1)
}

//!+
func main() {
	flag.Parse()
	worklist := make(chan []Link)  // lists of URLs, may have duplicates
	unseenLinks := make(chan Link) // de-duplicated URLs

	// Add command-line arguments to worklist.
	go func() { worklist <- createLinks([]string{"http://" + site + "/"}, 1) }()

	// Create 20 crawler goroutines to fetch each unseen link.
	for i := 0; i < 20; i++ {
		go func() {
			for link := range unseenLinks {
				foundLinks := crawl(link)
				go func() { worklist <- foundLinks }()
			}
		}()
	}

	// The main goroutine de-duplicates worklist items
	// and sends the unseen ones to the crawlers.
	seen := make(map[string]bool)
	for list := range worklist {
		for _, link := range list {
			if !seen[link.url] && link.depth <= maxDepth {
				seen[link.url] = true
				unseenLinks <- link
			}
		}
	}
}

//!-
