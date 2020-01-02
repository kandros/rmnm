package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sync"
	"sync/atomic"

	"github.com/dustin/go-humanize"
	"github.com/fatih/color"
)

func main() {
	if len(os.Args) != 2 {
		fmt.Print("Expected first argument to be a path")
		os.Exit(1)
	}
	path := os.Args[1]

	// shouldSpin := true
	// s := spin.New()

	// go func() {
	// 	for shouldSpin {
	// 		fmt.Printf("\r %s", s.Next())
	// 		time.Sleep(100 * time.Millisecond)
	// 	}
	// }()

	foundNodeModuleChan := make(chan string)

	var totalSize uint64 = 0

	wg := sync.WaitGroup{}
	go func() {
		for p := range foundNodeModuleChan {
			wg.Add(1)
			go func(p string) {
				defer wg.Done()
				s := folderSize(p)
				atomic.AddUint64(&totalSize, s)
				err := os.RemoveAll(p)
				if err != nil {
					panic(err)
				}

				fmt.Printf("\n❌  [%s]  \"%s\"", humanize.Bytes(s), color.HiYellowString(p))
			}(p)
		}
	}()

	scanForNodeModules(path, foundNodeModuleChan)
	wg.Wait()
	fmt.Printf("\ntotal size %s", humanize.Bytes(totalSize))

}

func isNodemoduleDir(fi os.FileInfo) bool {
	return fi.IsDir() && fi.Name() == "node_modules"
}

func scanForNodeModules(folderPath string, foundNodeModuleChan chan<- string) {
	filepath.Walk(folderPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			log.Fatal(err)
		}

		if isNodemoduleDir(info) {
			foundNodeModuleChan <- path
			return filepath.SkipDir
		}

		return nil
	})
}

func folderSize(path string) uint64 {
	var size int64 = 0
	err := filepath.Walk(path, func(path string, info os.FileInfo, err error) error {
		size += info.Size()
		if err != nil {
			return err
		}
		return nil
	})

	if err != nil {
		panic(err)
	}

	return uint64(size)
}
