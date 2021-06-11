/*
Упражнение 8.9. Напишите версию du, которая вычисляет и периодически выво­
дит отдельные итоговые величины для каждого из каталогов ro o t.
*/
package main

// The du3 variant traverses all directories in parallel.
// It uses a concurrency-limiting counting semaphore
// to avoid opening too many files at once.

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"sync"
	"time"
)

var vFlag = flag.Bool("v", false, "show verbose progress messages")

type Info struct {
	root   string
	nbytes int64
}

//!+
func main() {
	// ...determine roots...

	//!-
	flag.Parse()

	// Determine the initial directories.
	roots := flag.Args()
	if len(roots) == 0 {
		roots = []string{"."}
	}

	//!+
	// Traverse each root of the file tree in parallel.
	filesInfo := make(chan Info)
	var n sync.WaitGroup
	for _, root := range roots {
		n.Add(1)
		file, err := os.Open(root)
		if err != nil {
			log.Fatal(err)
		}
		fi, err := file.Stat()
		if err != nil {
			log.Fatal(err)
		}
		go walkDir(fi.Name(), root, &n, filesInfo)
		file.Close()
	}
	go func() {
		n.Wait()
		close(filesInfo)
	}()
	//!-

	// Print the results periodically.
	var tick <-chan time.Time
	if *vFlag {
		tick = time.Tick(500 * time.Millisecond)
	}
	var nbytes = make(map[string]int64, len(roots))
	var nfiles = make(map[string]int64, len(roots))
loop:
	for {
		select {
		case info, ok := <-filesInfo:
			if !ok {
				break loop // filesInfo was closed
			}
			nfiles[info.root]++
			nbytes[info.root] += info.nbytes
		case <-tick:
			printDiskUsage(nfiles, nbytes)
		}
	}

	printDiskUsage(nfiles, nbytes) // final totals
	//!+
	// ...select loop...
}

//!-

func printDiskUsage(nfiles, nbytes map[string]int64) {
	for k := range nbytes {
		fmt.Printf("%s: %d files  %.1f GB\n", k, nfiles[k], float64(nbytes[k])/1e9)
	}
}

// walkDir recursively walks the file tree rooted at dir
// and sends the size of each found file on filesInfo.
//!+walkDir
func walkDir(root string, dir string, n *sync.WaitGroup, filesInfo chan<- Info) {
	defer n.Done()
	for _, entry := range dirents(dir) {
		if entry.IsDir() {
			n.Add(1)
			subdir := filepath.Join(dir, entry.Name())
			go walkDir(root, subdir, n, filesInfo)
		} else {
			filesInfo <- Info{nbytes: entry.Size(), root: root}
		}
	}
}

//!-walkDir

//!+sema
// sema is a counting semaphore for limiting concurrency in dirents.
var sema = make(chan struct{}, 20)

// dirents returns the entries of directory dir.
func dirents(dir string) []os.FileInfo {
	sema <- struct{}{}        // acquire token
	defer func() { <-sema }() // release token
	// ...
	//!-sema

	entries, err := ioutil.ReadDir(dir)
	if err != nil {
		fmt.Fprintf(os.Stderr, "du: %v\n", err)
		return nil
	}
	return entries
}
