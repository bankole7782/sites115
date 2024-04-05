package sites115

import (
	"context"
	"fmt"
	"io/fs"
	"os"
	"slices"
	"strings"

	arrops "github.com/adam-hanna/arrayOperations"
	"github.com/kljensen/snowball"
	"github.com/mholt/archiver/v4"
	"github.com/russross/blackfriday"
)

type S1Object struct {
	mDTarPath  string
	iDXTarPath string
}

func Init(mdTarPath, idxTarPath string) (S1Object, error) {
	ret := S1Object{}
	if !doesPathExists(mdTarPath) {
		return ret, fmt.Errorf("file %s does not exists", mdTarPath)
	}
	if !doesPathExists(idxTarPath) {
		return ret, fmt.Errorf("file %s does not exists", idxTarPath)
	}

	return S1Object{mdTarPath, idxTarPath}, nil
}

func doesPathExists(p string) bool {
	if _, err := os.Stat(p); os.IsNotExist(err) {
		return false
	}
	return true
}

func (s1o *S1Object) ReadMDAsHTML(path string) (string, error) {
	mdFS, err := archiver.FileSystem(context.Background(), s1o.mDTarPath)
	if err != nil {
		return "", err
	}

	trueMDFS := mdFS.(fs.ReadFileFS)
	rawMD, err := trueMDFS.ReadFile(path)
	if err != nil {
		return "", err
	}

	return string(blackfriday.MarkdownCommon(rawMD)), nil
}

func getStopWords() []string {
	stopWordsStr := strings.ReplaceAll(string(StopWordsBytes), "\r", "")
	return strings.Split(stopWordsStr, "\n")
}

func (s1o *S1Object) ReadAllMD() ([]string, error) {
	mdFS, err := archiver.FileSystem(context.Background(), s1o.mDTarPath)
	if err != nil {
		return nil, err
	}

	trueMDFS := mdFS.(fs.ReadDirFS)
	dirFIs, err := trueMDFS.ReadDir("")
	if err != nil {
		return nil, err
	}

	allPaths := make([]string, 0)
	for _, dirFI := range dirFIs {
		allPaths = append(allPaths, dirFI.Name())
	}

	return allPaths, nil
}

func (s1o *S1Object) Search(searchStr string) ([]string, error) {
	idxFS, err := archiver.FileSystem(context.Background(), s1o.iDXTarPath)
	if err != nil {
		return nil, err
	}

	trueIdxFS := idxFS.(fs.ReadFileFS)

	allPaths := make([][]string, 0)

	stopWords := getStopWords()
	words := strings.Fields(searchStr)
	for _, word := range words {
		// stopwords check
		word = strings.ToLower(word)
		if slices.Contains(stopWords, word) {
			continue
		}

		stemmedWord, err := snowball.Stem(word, "english", true)
		if err != nil {
			continue
		}

		rawPaths, err := trueIdxFS.ReadFile(stemmedWord + ".txt")
		if err != nil {
			continue
		}

		pathsStr := strings.ReplaceAll(string(rawPaths), "\r", "")
		paths := strings.Split(pathsStr, "\n")
		allPaths = append(allPaths, paths)
	}

	return arrops.Union(allPaths...), nil
}
