package main

import (
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"slices"

	"github.com/kljensen/snowball"
	stripMarkdown "github.com/writeas/go-strip-markdown"
)

func makeIndex(basePath, indexesPath, path string) error {
	if !strings.HasSuffix(path, ".md") {
		return nil
	}

	mdShortPath := strings.Replace(path, basePath, "", 1)

	rawMD, _ := os.ReadFile(path)
	strippedMD := stripMarkdown.Strip(string(rawMD))

	strippedMD = regexp.MustCompile(`[^a-zA-Z0-9 ]+`).ReplaceAllString(strippedMD, " ")

	words := strings.Fields(strippedMD)
	stopWords := getStopWords()

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

		indexedWordPath := filepath.Join(indexesPath, stemmedWord+".txt")
		if DoesPathExists(indexedWordPath) {
			raw, err := os.ReadFile(indexedWordPath)
			if err != nil {
				continue
			}
			pathsStr := strings.ReplaceAll(string(raw), "\r", "")
			paths := strings.Split(pathsStr, "\n")
			if !slices.Contains(paths, mdShortPath) {
				paths = append(paths, mdShortPath)
			}
			toWrite := []byte(strings.Join(paths, "\n"))
			os.WriteFile(indexedWordPath, toWrite, 0777)
		} else {
			os.WriteFile(indexedWordPath, []byte(mdShortPath), 0777)
		}
	}

	return nil
}

func getStopWords() []string {
	stopWordsStr := strings.ReplaceAll(string(StopWordsBytes), "\r", "")
	return strings.Split(stopWordsStr, "\n")
}
