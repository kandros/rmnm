package main

import (
	"fmt"
	"os"
	"path/filepath"
)

const nodeModules = "node_modules"

func main() {
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

				// if current folder is node_modules dont walk into it
				return filepath.SkipDir
			} else {
				foldersToCheck = append(foldersToCheck, path)
			}
		}

		return nil
	}

	err := filepath.Walk(folderPath, walker)

	for _, f := range foldersToCheck {
		filepath.Walk(f, walker)
	}

	foldersDeleted := 0
	for _, f := range foldersToDelete {
		foldersDeleted++
		err := os.RemoveAll(f)
		fmt.Printf("💀  - deleted %s\n", f)
		if err != nil {
			panic(err)
		}
	}

	fmt.Printf("\ndeleted %d node_modules folders\n", foldersDeleted)

	if err != nil {
		panic(err)
	}
}
