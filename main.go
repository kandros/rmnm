package main

import (
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/tj/go-spin"
)

const nodeModules = "node_modules"

func run() {
	c := make(chan os.Signal)
	go func() {
		<-c
		os.Exit(1)
	}()

	folderPath, _ := os.Getwd()

	var foldersToCheck []string
	var foldersToDelete []string

	walker := func(path string, info os.FileInfo, err error) error {
		// skil first folder because it's folderPath
		if err != nil {
			return err
		}

		if info.IsDir() {
			if info.Name() == nodeModules {
				foldersToDelete = append(foldersToDelete, path)

				fmt.Printf("%s marked for deletion\n", path)

				// if current folder is node_modules dont walk into it
				return filepath.SkipDir
			} else {
				foldersToCheck = append(foldersToCheck, path)
			}
		}

		return nil
	}

	s := spin.New()
	shouldSpin := true

	go func() {
		for shouldSpin {
			fmt.Printf("\r %s", s.Next())
			time.Sleep(100 * time.Millisecond)
		}
	}()

	var wg sync.WaitGroup

	wg.Add(1)
	go func() {
		filepath.Walk(folderPath, walker)
		wg.Done()
	}()

	wg.Wait()

	for _, f := range foldersToCheck {
		x := f
		wg.Add(1)
		go func() {
			filepath.Walk(x, walker)
			wg.Done()
		}()
	}
	wg.Wait()

	foldersDeleted := 0

	for i, f := range foldersToDelete {
		a, b := i, f
		wg.Add(1)
		go func() {
			foldersDeleted++
			err := os.RemoveAll(b)
			fmt.Printf("💀  - deleted %s\n", f)
			if err != nil {
				panic(err)
			}
			wg.Done()
			if a == len(foldersToDelete) {
				shouldSpin = false
			}
		}()
	}

	wg.Wait()

	fmt.Printf("\ndeleted %d node_modules folders\n", foldersDeleted)
}

func main() {
	run()
}
