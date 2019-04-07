package main

import (
	"fmt"
	"os"
	"path/filepath"
)

const nodeModules = "node_modules"

func main() {
	folderPath := "folders"
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

	for _, f := range foldersToDelete {
		fmt.Printf("deleting %s\n", f)
		err := os.RemoveAll(f)
		if err != nil {
			panic(err)
		}
	}

	if err != nil {
		panic(err)
	}
}
