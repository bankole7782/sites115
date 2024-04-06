package sites115

import (
	"context"
	"fmt"
	"io"
	"io/fs"
	"os"
	"slices"
	"strings"

	"github.com/PuerkitoBio/goquery"
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

	trueMDFS := mdFS.(fs.ReadDirFS)
	mdHandler, err := trueMDFS.Open("root/" + path)
	if err != nil {
		return "", err
	}
	defer mdHandler.Close()

	rawMD, err := io.ReadAll(mdHandler)
	if err != nil {
		return "", err
	}

	return string(blackfriday.MarkdownCommon(rawMD)), nil
}

func (s1o *S1Object) ReadMDTitle(path string) (string, error) {
	html, err := s1o.ReadMDAsHTML(path)
	if err != nil {
		return "", err
	}

	doc, err := goquery.NewDocumentFromReader(strings.NewReader(html))
	if err != nil {
		return "", err
	}

	goQSelection := doc.Find("h1").First()
	return goQSelection.Text(), nil
}

func (s1o *S1Object) ReadMDAbstract(path string) (string, error) {
	html, err := s1o.ReadMDAsHTML(path)
	if err != nil {
		return "", err
	}

	doc, err := goquery.NewDocumentFromReader(strings.NewReader(html))
	if err != nil {
		return "", err
	}

	goQSelection := doc.Find("p").First()
	return goQSelection.Text(), nil
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
	dirFIs, err := trueMDFS.ReadDir("root")
	if err != nil {
		return nil, err
	}

	allPaths := make([]string, 0)
	for _, dirFI := range dirFIs {
		if dirFI.IsDir() {
			innerDirFIs, err := trueMDFS.ReadDir("root/" + dirFI.Name())
			if err != nil {
				return nil, err
			}
			for _, innerFI := range innerDirFIs {
				allPaths = append(allPaths, dirFI.Name()+"/"+innerFI.Name())

			}
		} else {
			allPaths = append(allPaths, dirFI.Name())
		}
	}

	return allPaths, nil
}

func (s1o *S1Object) Search(searchStr string) ([]string, error) {
	idxFS, err := archiver.FileSystem(context.Background(), s1o.iDXTarPath)
	if err != nil {
		return nil, err
	}

	// trueIdxFS := idxFS.(fs.ReadDirFS)

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

		idxHandler, err := idxFS.Open("root/" + stemmedWord + ".txt")
		if err != nil {
			continue
		}
		defer idxHandler.Close()

		rawPaths, err := io.ReadAll(idxHandler)
		if err != nil {
			continue
		}

		pathsStr := strings.ReplaceAll(string(rawPaths), "\r", "")
		paths := strings.Split(pathsStr, "\n")
		allPaths = append(allPaths, paths)
	}

	return arrops.Intersect(allPaths...), nil
}
