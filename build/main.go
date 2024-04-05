package main

import (
	"context"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"github.com/mholt/archiver/v4"
)

func main() {
	if len(os.Args) < 2 {
		panic("expecting a markdown files folder")
	}
	markdownPath := os.Args[1]

	if !strings.HasSuffix(markdownPath, "/") {
		markdownPath += "/"
	}
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

	// make archive of markdown files
	mdFilesMap := make(map[string]string)
	filepath.Walk(markdownPath, func(path string, info fs.FileInfo, err error) error {
		if err != nil {
			fmt.Printf("prevent panic by handling failure accessing a path %q: %v\n", path, err)
			return err
		}
		if info.IsDir() {
			return nil
		} else {
			// fmt.Println(path)
			shortPath := strings.ReplaceAll(path, markdownPath, "")
			mdFilesMap[path] = filepath.Join("root", shortPath)
		}

		return nil
	})
	mdFiles, _ := archiver.FilesFromDisk(nil, mdFilesMap)

	projectName := filepath.Base(markdownPath)
	projectDir := filepath.Dir(markdownPath)

	mdZipPath := filepath.Join(projectDir, projectName+"_md.tar.gz")
	mdZipPathHandle, err := os.Create(mdZipPath)
	if err != nil {
		panic(err)
	}
	defer mdZipPathHandle.Close()

	format := archiver.CompressedArchive{
		Compression: archiver.Gz{},
		Archival:    archiver.Tar{},
	}

	err = format.Archive(context.Background(), mdZipPathHandle, mdFiles)
	if err != nil {
		panic(err)
	}

	// make archive of indexes
	idxFilesMap := make(map[string]string)
	filepath.Walk(indexesPath, func(path string, info fs.FileInfo, err error) error {
		if err != nil {
			fmt.Printf("prevent panic by handling failure accessing a path %q: %v\n", path, err)
			return err
		}
		if info.IsDir() {
			return nil
		} else {
			// fmt.Println(path)
			shortPath := strings.ReplaceAll(path, indexesPath, "")
			idxFilesMap[path] = shortPath
		}

		return nil
	})
	idxFiles, _ := archiver.FilesFromDisk(nil, idxFilesMap)

	idxZipPath := filepath.Join(projectDir, projectName+"_idx.tar.gz")
	idxZipPathHandle, err := os.Create(idxZipPath)
	if err != nil {
		panic(err)
	}
	defer idxZipPathHandle.Close()

	err = format.Archive(context.Background(), idxZipPathHandle, idxFiles)
	if err != nil {
		panic(err)
	}

}
