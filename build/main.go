package main

import (
	"os"
	"path/filepath"
)

func main() {
	if len(os.Args) < 2 {
		panic("expecting a markdown files folder")
	}
	markdownPath := os.Args[1]
	// project := filepath.Base(markdownPath)

	dirFIs, err := os.ReadDir(markdownPath)
	if err != nil {
		panic(err)
	}

	indexesPath := filepath.Join(os.TempDir(), "s115", UntestedRandomString(5))
	os.MkdirAll(indexesPath, 0777)

	for _, dirFI := range dirFIs {
		if dirFI.IsDir() {
			// repeat what happens to files
			innerDirFIs, _ := os.ReadDir(filepath.Join(markdownPath, dirFI.Name()))
			for _, innerDirFI := range innerDirFIs {
				innerPath := filepath.Join(markdownPath, dirFI.Name(), innerDirFI.Name())
				makeIndex(markdownPath, indexesPath, innerPath)
			}
		} else {
			makeIndex(markdownPath, indexesPath, filepath.Join(markdownPath, dirFI.Name()))
		}
	}
}
