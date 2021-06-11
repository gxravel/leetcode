// Copyright © 2016 Alan A. A. Donovan & Brian W. Kernighan.
// License: https://creativecommons.org/licenses/by-nc-sa/4.0/

// See page 243.

// Crawl3 crawls web links starting with the command-line arguments.
//
// This version uses bounded parallelism.
// For simplicity, it does not address the termination problem.
//
package main

import (
	"flag"
	"fmt"
	"log"

	"github.com/gxravel/leetcode/exercises/gopl/ch8/ex-6/links"
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

func crawl(link Link) []Link {
	fmt.Println(link.depth, link.url)
	urls, err := links.Extract(link.url)
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
