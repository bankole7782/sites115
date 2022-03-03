package main

import (
  "encoding/json"
  "path/filepath"
  "strings"
  "os"
  "fmt"
  "strconv"
  "github.com/microcosm-cc/bluemonday"
  "github.com/bankole7782/sites115/sites115s"
)

func generateIndexes(siteName string) {

  rootPath, _ := GetRootPath()

  indexesPath := filepath.Join(rootPath, siteName, "out", "_indexes")
  os.MkdirAll(indexesPath, 0777)

  dir := filepath.Join(rootPath, siteName, "out")

  // get all pages
  allPages := make([]string, 0)
  err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
    if err != nil {
      return err
    }


    if ! info.IsDir() {
      if ! strings.HasSuffix(dir, "/") {
        dir += "/"
      }

      toRemove := filepath.Join(rootPath, siteName, "out")
      if strings.HasPrefix(path, filepath.Join(toRemove, "static")) {
        return nil
      }
      if strings.HasPrefix(path, filepath.Join(toRemove, "_indexes")) {
        return nil
      }

      pathToWrite := strings.Replace(path, dir, "", 1)
      allPages = append(allPages, pathToWrite)

    }
    return nil
  })

  if err != nil {
    panic(err)
  }

  allPagesMap := make(map[int]string)
  for i, page := range allPages {
    allPagesMap[i+1] = page
  }

  jsonBytes, err := json.Marshal(allPagesMap)
  if err != nil {
    panic(err)
  }
  os.WriteFile(filepath.Join(rootPath, siteName, "out", "_indexes", "allpages.json"), jsonBytes, 0777)

  // generate indexes
  err = filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
    if err != nil {
      return err
    }


    if ! info.IsDir() {
      if ! strings.HasSuffix(dir, "/") {
        dir += "/"
      }

      toRemove := filepath.Join(rootPath, siteName, "out")
      if strings.HasPrefix(path, filepath.Join(toRemove, "static")) {
        return nil
      }
      if strings.HasPrefix(path, filepath.Join(toRemove, "_indexes")) {
        return nil
      }

      pathToWrite := strings.Replace(path, dir, "", 1)
      index := 0
      for k, v := range allPagesMap {
        if v == pathToWrite {
          index = k
          break
        }
      }

      doIndex(path, strconv.Itoa(index), siteName)
    }
    return nil
  })
}


func doIndex(textPath, index, siteName string) {
	raw, err := os.ReadFile(textPath)
	if err != nil {
		return
	}

	textStrippedOfHtml := bluemonday.StrictPolicy().Sanitize(string(raw))
	words := strings.Fields(textStrippedOfHtml)

	wordCountMap := make(map[string]int64)
	for _, word := range words {
		// clean the word.
		word = sites115s.CleanWord(word)
		if word == "" {
			continue
		}
		if sites115s.FindIn(sites115s.STOP_WORDS, word) != -1 {
			continue
		}

		oldCount, ok := wordCountMap[word]
		if ! ok {
			wordCountMap[word] = 1
		} else {
			wordCountMap[word] = oldCount + 1
		}
	}

	rootPath, _ := GetRootPath()

	if ! strings.HasSuffix(rootPath, "/") {
		rootPath += "/"
	}

	for word, wordCount := range wordCountMap {
		dirToMake := filepath.Join(rootPath, siteName, "out", "_indexes", word)
		os.MkdirAll(dirToMake, 0777)
		err = os.WriteFile(filepath.Join(dirToMake, index), []byte(fmt.Sprintf("%d", wordCount)), 0777)
		if err != nil {
			fmt.Printf("word is : '%s'\n", word)
      return
		}
	}

}
